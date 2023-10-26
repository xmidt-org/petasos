// SPDX-FileCopyrightText: 2016 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0
package main

import (
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof" //nolint:gosec
	"os"
	"os/signal"
	"runtime"

	"github.com/justinas/alice"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/xmidt-org/candlelight"
	"github.com/xmidt-org/webpa-common/v2/adapter"
	"github.com/xmidt-org/webpa-common/v2/concurrent"          // nolint: staticcheck
	"github.com/xmidt-org/webpa-common/v2/device"              // nolint: staticcheck
	"github.com/xmidt-org/webpa-common/v2/server"              // nolint: staticcheck
	"github.com/xmidt-org/webpa-common/v2/service"             // nolint: staticcheck
	"github.com/xmidt-org/webpa-common/v2/service/monitor"     // nolint: staticcheck
	"github.com/xmidt-org/webpa-common/v2/service/servicecfg"  // nolint: staticcheck
	"github.com/xmidt-org/webpa-common/v2/service/servicehttp" // nolint: staticcheck
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

const (
	applicationName  = "petasos"
	tracingConfigKey = "tracing"
)

var (
	GitCommit = "undefined"
	Version   = "undefined"
	BuildTime = "undefined"
)

func loadTracing(v *viper.Viper, appName string) (candlelight.Tracing, error) {

	var config candlelight.Config
	err := v.UnmarshalKey(tracingConfigKey, &config)
	if err != nil {
		return candlelight.Tracing{}, err
	}
	config.ApplicationName = appName
	tracing, err := candlelight.New(config)
	if err != nil {
		return candlelight.Tracing{}, err
	}
	return tracing, nil
}

// petasos is the driver function for Petasos.  It performs everything main() would do,
// except for obtaining the command-line arguments (which are passed to it).
func petasos(arguments []string) int {
	//
	// Initialize the server environment: command-line flags, Viper, logging, and the WebPA instance
	//

	var (
		f = pflag.NewFlagSet(applicationName, pflag.ContinueOnError)
		v = viper.New()

		logger, metricsRegistry, webPA, err = server.Initialize(applicationName, arguments, f, v, service.Metrics)
	)

	if parseErr, done := printVersion(f, arguments); done {
		// if we're done, we're exiting no matter what
		if parseErr != nil {
			logger.Error("failed to parse arguments. detailed error:", zap.Error(parseErr))
			os.Exit(1)
		}
		os.Exit(0)
	}

	if err != nil {
		logger.Error("Unable to initialize Viper environment", zap.Error(err))
		return 1
	}

	//
	// Now, initialize the service discovery infrastructure
	//
	var log = &adapter.Logger{
		Logger: logger,
	}

	e, err := servicecfg.NewEnvironment(log, v.Sub("service"))
	if err != nil {
		logger.Error("Unable to initialize service discovery environment", zap.Error(err))
		return 2
	} else if e == nil {
		logger.Error("Petasos requires service discovery")
		return 2
	}

	logger.Info("configuration file successfully unmarshaled", zap.Any("configurationFile", v.ConfigFileUsed()))
	tracing, err := loadTracing(v, applicationName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to build tracing component: %v \n", err)
		return 1
	}
	logger.Info("tracing status", zap.Bool("enabled", !tracing.IsNoop()))
	accessor := new(service.UpdatableAccessor)

	redirectHandler := &servicehttp.RedirectHandler{
		KeyFunc:      device.IDHashParser,
		Accessor:     accessor,
		RedirectCode: http.StatusTemporaryRedirect,
	}

	options := []otelhttp.Option{
		otelhttp.WithPropagators(tracing.Propagator()),
		otelhttp.WithTracerProvider(tracing.TracerProvider()),
	}
	decoratedHandler := alice.New(setLogger(logger, header("X-Webpa-Device-Name", "device_id")), candlelight.EchoFirstTraceNodeInfo(tracing.Propagator())).Then(redirectHandler)

	handler := otelhttp.NewHandler(decoratedHandler, "mainSpan", options...)

	_, petasosServer, done := webPA.Prepare(logger, nil, metricsRegistry, handler)
	signals := make(chan os.Signal, 10)

	_, err = monitor.New(
		monitor.WithLogger(logger),
		monitor.WithEnvironment(e),
		monitor.WithListeners(
			monitor.NewMetricsListener(metricsRegistry),
			monitor.NewAccessorListener(e.AccessorFactory(), accessor.Update),
		),
		monitor.WithFilter(monitor.NewNormalizeFilter(e.DefaultScheme())),
	)

	if err != nil {
		logger.Error("Unable to start service discovery monitor", zap.Error(err))
		return 3
	}

	//
	// Execute the runnable, which runs all the servers, and wait for a signal
	//
	waitGroup, shutdown, err := concurrent.Execute(petasosServer)
	if err != nil {
		logger.Error("Ubale to start petasos", zap.Error(err))
		return 4
	}

	signal.Notify(signals, os.Interrupt)
	for exit := false; !exit; {
		select {
		case s := <-signals:
			logger.Info("exiting due to signal", zap.Any("signal", s))
			exit = true
		case <-done:
			logger.Error("one or more servers exited")
			exit = true
		}
	}

	close(shutdown)
	waitGroup.Wait()

	return 0
}

func printVersion(f *pflag.FlagSet, arguments []string) (error, bool) {
	printVer := f.BoolP("version", "v", false, "displays the version number")
	if err := f.Parse(arguments); err != nil {
		return err, true
	}

	if *printVer {
		printVersionInfo(os.Stdout)
		return nil, true
	}
	return nil, false
}

func printVersionInfo(writer io.Writer) {
	fmt.Fprintf(writer, "%s:\n", applicationName)
	fmt.Fprintf(writer, "  version: \t%s\n", Version)
	fmt.Fprintf(writer, "  go version: \t%s\n", runtime.Version())
	fmt.Fprintf(writer, "  built time: \t%s\n", BuildTime)
	fmt.Fprintf(writer, "  git commit: \t%s\n", GitCommit)
	fmt.Fprintf(writer, "  os/arch: \t%s/%s\n", runtime.GOOS, runtime.GOARCH)
}

func main() {
	os.Exit(petasos(os.Args))
}
