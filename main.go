package main

import (
	"fmt"

	"github.com/err0r500/go-realworld-clean/implem/dummy.articleValidator"
	"github.com/err0r500/go-realworld-clean/implem/gin.server"
	"github.com/err0r500/go-realworld-clean/implem/gosimple.slugger"
	"github.com/err0r500/go-realworld-clean/implem/jwt.authHandler"
	"github.com/err0r500/go-realworld-clean/implem/logrus.logger"
	"github.com/err0r500/go-realworld-clean/implem/memory.articleRW"
	"github.com/err0r500/go-realworld-clean/implem/memory.commentRW"
	"github.com/err0r500/go-realworld-clean/implem/memory.tagsRW"
	"github.com/err0r500/go-realworld-clean/implem/memory.userRW"
	"github.com/err0r500/go-realworld-clean/implem/user.validator"
	"github.com/err0r500/go-realworld-clean/infra"
	"github.com/err0r500/go-realworld-clean/uc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

func run() {
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
