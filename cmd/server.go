package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/codefluence-x/altair/controller"
	"github.com/codefluence-x/altair/forwarder"
	"github.com/codefluence-x/altair/loader"
	"github.com/codefluence-x/altair/provider"
	"github.com/codefluence-x/journal"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run API gateway services.",
	Run: func(cmd *cobra.Command, args []string) {
		defer closeConnection()
		if err := fabricateConnection(); err != nil {
			journal.Error("Error running altair:", err).SetTags("altair", "main").Log()
			return
		}

		if err := runAPI(); err != nil {
			journal.Error("Error running altair API:", err).SetTags("altair", "main").Log()
		}
	},
}

func runAPI() error {
	gin.SetMode(gin.ReleaseMode)

	apiEngine = gin.New()
	apiEngine.GET("/health", controller.Health)

	pluginEngine := apiEngine.Group("/_plugins/", gin.BasicAuth(gin.Accounts{
		appConfig.BasicAuthUsername(): appConfig.BasicAuthPassword(),
	}))

	appBearer := loader.AppBearer(pluginEngine, appConfig)
	dbBearer := loader.DatabaseBearer(databases, dbConfigs)

	provider.Metric(appBearer)
	_ = provider.Plugin(appBearer, dbBearer, pluginBearer)

	// Route Engine
	routeCompiler := forwarder.Route().Compiler()
	routeObjects, err := routeCompiler.Compile("./routes")
	if err != nil {
		journal.Error("Error compiling routes", err).
			SetTags("altair", "main").
			Log()
		return err
	}

	metricProvider, _ := appBearer.MetricProvider()
	err = forwarder.Route().Generator().Generate(apiEngine, metricProvider, routeObjects, appBearer.DownStreamPlugins())
	if err != nil {
		journal.Error("Error generating routes", err).
			SetTags("altair", "main").
			Log()
		return err
	}

	gracefulSignal := make(chan os.Signal, 1)
	signal.Notify(gracefulSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%d", appConfig.Port()),
			Handler: apiEngine,
		}

		journal.Info(fmt.Sprintf("Running Altair in: %d", appConfig.Port())).Log()

		if err := srv.ListenAndServe(); err != nil {
			journal.Error("Error running api engine", err).
				SetTags("altair", "main").
				Log()
		}
	}()

	closeSignal := <-gracefulSignal
	journal.Info(fmt.Sprintf("Receiving %s signal.... Cleaning up processes.", closeSignal.String())).Log()
	return nil
}
