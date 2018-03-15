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
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"

	"github.com/Comcast/webpa-common/concurrent"
	"github.com/Comcast/webpa-common/device"
	"github.com/Comcast/webpa-common/logging"
	"github.com/Comcast/webpa-common/server"
	"github.com/Comcast/webpa-common/service"
	"github.com/Comcast/webpa-common/service/monitor"
	"github.com/Comcast/webpa-common/service/servicecfg"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	applicationName       = "petasos"
	release               = "Developer"
	defaultVnodeCount int = 211
)

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
	}

	infoLog.Log("configurationFile", v.ConfigFileUsed())

	var (
		accessor = new(service.UpdatableAccessor)

		redirectHandler = &service.RedirectHandler{
			Logger:       logger,
			KeyFunc:      device.IDHashParser,
			Accessor:     accessor,
			RedirectCode: http.StatusTemporaryRedirect,
		}

		_, petasosServer = webPA.Prepare(logger, nil, metricsRegistry, redirectHandler)
		signals          = make(chan os.Signal, 1)
	)

	_, err = monitor.New(
		monitor.WithLogger(logger),
		monitor.WithEnvironment(e),
		monitor.WithListeners(
			monitor.NewMetricsListener(metricsRegistry),
			monitor.NewAccessorListener(e.AccessorFactory(), accessor.Update),
		),
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

	signal.Notify(signals)
	s := server.SignalWait(infoLog, signals, os.Kill, os.Interrupt)
	errorLog.Log(logging.MessageKey(), "exiting due to signal", "signal", s)
	close(shutdown)
	waitGroup.Wait()

	return 0
}

func main() {
	os.Exit(petasos(os.Args))
}
