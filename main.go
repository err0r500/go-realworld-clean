package main

import (
	"fmt"
	"io"
	"log"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-lib/metrics"

	articleValidator "github.com/err0r500/go-realworld-clean/implem/dummy.articleValidator"
	server "github.com/err0r500/go-realworld-clean/implem/gin.server"
	slugger "github.com/err0r500/go-realworld-clean/implem/gosimple.slugger"
	"github.com/err0r500/go-realworld-clean/implem/jwt.authHandler"
	logger "github.com/err0r500/go-realworld-clean/implem/logrus.logger"
	articleRW "github.com/err0r500/go-realworld-clean/implem/memory.articleRW"
	commentRW "github.com/err0r500/go-realworld-clean/implem/memory.commentRW"
	tagsRW "github.com/err0r500/go-realworld-clean/implem/memory.tagsRW"
	userRW "github.com/err0r500/go-realworld-clean/implem/memory.userRW"
	validator "github.com/err0r500/go-realworld-clean/implem/user.validator"
	"github.com/err0r500/go-realworld-clean/infra"
	"github.com/err0r500/go-realworld-clean/uc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

// Build number and versions injected at compile time, set yours
var (
	Version = "unknown"
	Build   = "unknown"
)

// the command to run the server
var rootCmd = &cobra.Command{
	Use:   "go-realworld-clean",
	Short: "Runs the server",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show build and version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Build: %s\nVersion: %s\n", Build, Version)
	},
}

func main() {
	rootCmd.AddCommand(versionCmd)
	cobra.OnInitialize(infra.CobraInitialization)

	infra.LoggerConfig(rootCmd)
	infra.ServerConfig(rootCmd)
	infra.DatabaseConfig(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal()
	}
}

func setTracer() (opentracing.Tracer, io.Closer) {
	cfg, err := jaegercfg.FromEnv()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)
	if err != nil {
		log.Fatal(err)
	}
	return tracer, closer
}

func run() {
	tracer, closer := setTracer()
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	ginServer := infra.NewServer(
		viper.GetInt("server.port"),
		infra.DebugMode,
	)

	authHandler := jwt.New(viper.GetString("jwt.Salt"))
	routerLogger := logger.NewLogger("TEST",
		viper.GetString("log.level"),
		viper.GetString("log.format"),
	)

	server.NewRouterWithLogger(
		uc.HandlerConstructor{
			Logger:           routerLogger,
			UserRW:           userRW.New(),
			ArticleRW:        articleRW.New(),
			UserValidator:    validator.New(),
			AuthHandler:      authHandler,
			Slugger:          slugger.New(),
			ArticleValidator: articleValidator.New(),
			TagsRW:           tagsRW.New(),
			CommentRW:        commentRW.New(),
		}.New(),
		authHandler,
		routerLogger,
	).SetRoutes(ginServer.Router)

	ginServer.Start()
}
