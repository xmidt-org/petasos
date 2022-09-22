/**
 * Copyright 2016 Comcast Cable Communications Management, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package main

import (
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof" //nolint:gosec
	"os"
	"os/signal"
	"runtime"

	"github.com/go-kit/log/level"
	"github.com/justinas/alice"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/xmidt-org/candlelight"
	"github.com/xmidt-org/webpa-common/v2/concurrent"          // nolint: staticcheck
	"github.com/xmidt-org/webpa-common/v2/device"              // nolint: staticcheck
	"github.com/xmidt-org/webpa-common/v2/logging"             // nolint: staticcheck
	"github.com/xmidt-org/webpa-common/v2/logging/logginghttp" // nolint: staticcheck
	"github.com/xmidt-org/webpa-common/v2/server"              // nolint: staticcheck
	"github.com/xmidt-org/webpa-common/v2/service"             // nolint: staticcheck
	"github.com/xmidt-org/webpa-common/v2/service/monitor"     // nolint: staticcheck
	"github.com/xmidt-org/webpa-common/v2/service/servicecfg"  // nolint: staticcheck
	"github.com/xmidt-org/webpa-common/v2/service/servicehttp" // nolint: staticcheck
	"github.com/xmidt-org/webpa-common/v2/xhttp/xcontext"      // nolint: staticcheck
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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
	var traceConfig candlelight.TraceConfig
	tracing, err := candlelight.New(config)
	if err != nil {
		return candlelight.Tracing{}, err
	}

	error := v.UnmarshalKey(tracingConfigKey, &config)
	if error != nil {
		return candlelight.Tracing{}, err
	}
	config.ApplicationName = appName
	tracerProvider, err := candlelight.ConfigureTracerProvider(config)
	if err != nil {
		return candlelight.Tracing{}, err
	}
	// if len(config.Provider) != 0 && config.Provider != candlelight.DefaultTracerProvider {
	// 	tracing.Enabled = true
	// }
	traceConfig.TraceProvider = tracerProvider
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
		infoLog                             = logging.Info(logger)
		errorLog                            = logging.Error(logger)
	)

	if parseErr, done := printVersion(f, arguments); done {
		// if we're done, we're exiting no matter what
		if parseErr != nil {
			friendlyError := fmt.Sprintf("failed to parse arguments. detailed error: %s", parseErr)
			logging.Error(logger).Log(
				logging.ErrorKey(),
				friendlyError)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if err != nil {
		errorLog.Log(logging.MessageKey(), "Unable to initialize Viper environment", logging.ErrorKey(), err)
		return 1
	}

	//
	// Now, initialize the service discovery infrastructure
	//

	e, err := servicecfg.NewEnvironment(logger, v.Sub("service"))
	if err != nil {
		errorLog.Log(logging.MessageKey(), "Unable to initialize service discovery environment", logging.ErrorKey(), err)
		return 2
	} else if e == nil {
		errorLog.Log(logging.MessageKey(), "Petasos requires service discovery")
		return 2
	}

	infoLog.Log("configurationFile", v.ConfigFileUsed())

	tracing, err := loadTracing(v, applicationName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to build tracing component: %v \n", err)
		return 1
	}
	//infoLog.Log(logging.MessageKey(), "tracing status", "enabled", tracing.Enabled)

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
	requestFunc := logginghttp.SetLogger(logger, logginghttp.Header("X-Webpa-Device-Name", "device_id"), logginghttp.Header("Authorization", "authorization"), logginghttp.LoggerFunc(candlelight.InjectTraceInfoInLogger()))
	decoratedHandler := alice.New(xcontext.Populate(requestFunc), candlelight.EchoFirstTraceNodeInfo(tracing.Propagator())).Then(redirectHandler)

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
		errorLog.Log(logging.MessageKey(), "Unable to start service discovery monitor", logging.ErrorKey(), err)
		return 3
	}

	//
	// Execute the runnable, which runs all the servers, and wait for a signal
	//
	waitGroup, shutdown, err := concurrent.Execute(petasosServer)
	if err != nil {
		errorLog.Log(logging.MessageKey(), "Unable to start petasos", logging.ErrorKey(), err)
		return 4
	}

	signal.Notify(signals, os.Interrupt)
	for exit := false; !exit; {
		select {
		case s := <-signals:
			logger.Log(level.Key(), level.ErrorValue(), logging.MessageKey(), "exiting due to signal", "signal", s)
			exit = true
		case <-done:
			logger.Log(level.Key(), level.ErrorValue(), logging.MessageKey(), "one or more servers exited")
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
