package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Comcast/golang-discovery-client/service"
	"github.com/Comcast/webpa-common/concurrent"
	"github.com/Comcast/webpa-common/device"
	"github.com/Comcast/webpa-common/fact"
	"github.com/Comcast/webpa-common/handler"
	"github.com/Comcast/webpa-common/hash"
	"github.com/Comcast/webpa-common/health"
	"github.com/Comcast/webpa-common/logging"
	"github.com/Comcast/webpa-common/logging/golog"
	"github.com/Comcast/webpa-common/server"
	"github.com/billhathaway/consistentHash"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
)

const (
	release               = "Developer"
	defaultVnodeCount int = 211
)

// Configuration hold all the configurable options for petasos
type Configuration struct {
	server.Configuration
	AlternateAddress string                   `json:"alternateAddress"`
	LoggerFactory    golog.LoggerFactory      `json:"log"`
	DiscoveryBuilder service.DiscoveryBuilder `json:"discovery"`
	VnodeCount       int                      `json:"vnodeCount"`
}

func main() {
	viper := viper.New()
	if err := server.ReadInConfig("petasos", viper, nil, nil); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	configuration := new(Configuration)
	if err := server.ReadConfigurationFile(configurationFile, configuration); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read configuration file: %v\n", err)
		os.Exit(1)
	}

	logger, err := configuration.LoggerFactory.NewLogger("petasos")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Logger could not be created: %v\n", err)
		os.Exit(1)
	}

	if len(configuration.DiscoveryBuilder.Watches) != 1 {
		logger.Error("There must be exactly (1) watched service")
		os.Exit(1)
	}

	watchedServiceName := configuration.DiscoveryBuilder.Watches[0]
	vnodeCount := configuration.VnodeCount
	if vnodeCount < 1 {
		vnodeCount = defaultVnodeCount
	}

	logger.Info("Using configuration: %s", configuration)

	os.Exit(func() int {
		serviceHashHolder := &hash.ServiceHashHolder{}
		printlogger := logging.PrintLogger{logger}
		discovery, err := configuration.DiscoveryBuilder.New(printlogger)
		if err != nil {
			logger.Error("Unable to create discovery client: %s", err)
			return 1
		}

		discovery.AddListener(
			watchedServiceName,
			service.ListenerFunc(func(serviceName string, instances service.Instances) {
				logger.Info("Rehashing service nodes [%s]: %s", serviceName, instances)
				newHash := consistentHash.New()
				newHash.SetVnodeCount(vnodeCount)
				instances.ToKeys(service.HttpAddress, newHash)
				serviceHashHolder.Update(newHash)
			}),
		)

		petasosHealth := health.New(
			configuration.HealthCheckInterval(),
			logger,
			handler.TotalRequestsReceived,
			handler.TotalRequestSuccessfullyServiced,
			handler.TotalRequestDenied,
		)

		healthServer := (&server.Builder{
			Name:    "petasos-health",
			Address: configuration.HealthAddress(),
			Logger:  logger,
			Handler: petasosHealth,
		}).Build()

		pprofServer := (&server.Builder{
			Name:    "petasos-pprof",
			Address: configuration.PprofAddress(),
			Logger:  logger,
			Handler: http.DefaultServeMux,
		}).Build()

		ctx := fact.SetLogger(context.Background(), logger)
		petasosHandler := handler.Chain{
			handler.Listen(handler.NewHealthRequestListener(petasosHealth)),
			handler.DeviceId(),
			handler.Convey(),
		}.Decorate(ctx, handler.Hash(serviceHashHolder))

		petasosPrimaryServer := (&server.Builder{
			Name:            "petasos",
			Address:         configuration.PrimaryAddress(),
			CertificateFile: configuration.CertificateFile,
			KeyFile:         configuration.KeyFile,
			Logger:          logger,
			Handler:         petasosHandler,
		}).Build()

		runnables := concurrent.RunnableSet{
			discovery,
			petasosHealth,
			healthServer,
			pprofServer,
			petasosPrimaryServer,
		}

		if alternateAddress, ok := configuration.AlternateAddress(); ok {
			runnables = append(
				runnables,
				(&server.Builder{
					Name:    "petasos-alt",
					Address: alternateAddress,
					Logger:  logger,
					Handler: petasosHandler,
				}).Build(),
			)
		}

		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)
		err = concurrent.Await(runnables, signals)
		if err != nil {
			logger.Error("Petasos exiting: %v", err)
		} else {
			logger.Info("Petasos exiting")
		}

		return 0
	}())
}
