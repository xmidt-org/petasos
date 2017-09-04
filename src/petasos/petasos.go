package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/Comcast/webpa-common/concurrent"
	"github.com/Comcast/webpa-common/device"
	"github.com/Comcast/webpa-common/logging"
	"github.com/Comcast/webpa-common/server"
	"github.com/Comcast/webpa-common/service"
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

		logger, webPA, err = server.Initialize(applicationName, arguments, f, v)
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to initialize Viper environment: %s\n", err)
		return 1
	}

	//
	// Now, initialize the service discovery infrastructure
	//

	serviceOptions, err := service.FromViper(service.Sub(v))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read service discovery options: %s\n", err)
		return 2
	}

	logging.Info(logger).Log("configurationFile", v.ConfigFileUsed(), "serviceOptions", serviceOptions)
	serviceOptions.Logger = logger
	services, err := service.New(serviceOptions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to initialize service discovery: %s\n", err)
		return 2
	}

	instancer, err := services.NewInstancer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to obtain service discovery instancer: %s\n", err)
		return 2
	}

	var (
		accessor     = new(service.UpdatableAccessor)
		subscription = service.Subscribe(serviceOptions, instancer)

		redirectHandler = &service.RedirectHandler{
			Logger:       logger,
			KeyFunc:      device.IDHashParser,
			Accessor:     accessor,
			RedirectCode: http.StatusTemporaryRedirect,
		}

		_, runnable = webPA.Prepare(logger, nil, redirectHandler)
		signals     = make(chan os.Signal, 1)
	)

	accessor.Consume(subscription)

	//
	// Execute the runnable, which runs all the servers, and wait for a signal
	//

	if err := concurrent.Await(runnable, signals); err != nil {
		fmt.Fprintf(os.Stderr, "Error when starting %s: %s", applicationName, err)
		return 4
	}

	return 0
}

func main() {
	os.Exit(petasos(os.Args))
}
