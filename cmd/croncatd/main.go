package main

import (
	"fmt"
	"os"

	"github.com/CronCats/croncat-go/internal/app"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Logger *logrus.Entry

const Logo = "                                         _  \n" +
	"           ___ _ __ ___  _ __   ___ __ _| |_ \n" +
	"  /\\_/\\   / __| '__/ _ \\| '_ \\ / __/ _` | __/ \n" +
	" ( o.o ) | (__| | | (_) | | | | (_| (_| | |_ \n" +
	"  > ^ <   \\___|_|  \\___/|_| |_|\\___\\__,_|\\__\\ \n" +
	"===============================================\n\n"

func initConfig() {
	// Set level from env var
	level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true, DisableLevelTruncation: true})
}

func initContextLogger(cmd *cobra.Command) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"agent":    cmd.Flag("agent").Value.String(),
		"chain-id": cmd.Flag("chain-id").Value.String(),
	})
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.croncatd.yaml)")

	rootCmd.PersistentFlags().StringP("agent", "a", "", "agent id")
	rootCmd.MarkPersistentFlagRequired("agent")
	rootCmd.PersistentFlags().StringP("chain-id", "i", "", "chain id")
	rootCmd.MarkPersistentFlagRequired("chain-id")

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "debug mode")
	rootCmd.PersistentFlags().BoolP("no-frills", "", false, "wether to display cute things")

	rootCmd.AddCommand(goCommand)
	rootCmd.AddCommand(registerCommand)
	rootCmd.AddCommand(unregisterCommand)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func preRun(cmd *cobra.Command, args []string) {
	Logger = initContextLogger(cmd)

	if cmd.Flag("no-frills").Value.String() != "true" {
		fmt.Print(Logo)
	}
}

var rootCmd = &cobra.Command{
	PersistentPreRun: preRun,
	Use:              "croncatd",
	Short:            "croncatd is a daemon that runs croncat tasks",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var goCommand = &cobra.Command{
	Use:   "go",
	Short: "Run croncat tasks",
	Run: func(cmd *cobra.Command, args []string) {
		chainId := cmd.Flag("chain-id").Value.String()

		app := app.NewApp(chainId, Logger)
		app.Run()
	},
}

var registerCommand = &cobra.Command{
	Use:   "register",
	Short: "Register a croncat agent",
	Run: func(cmd *cobra.Command, args []string) {
		chainId := cmd.Flag("chain-id").Value.String()

		app := app.NewApp(chainId, Logger)
		app.Register()
	},
}

var unregisterCommand = &cobra.Command{
	Use:   "unregister",
	Short: "Unregister a croncat agent",
	Run: func(cmd *cobra.Command, args []string) {
		chainId := cmd.Flag("chain-id").Value.String()

		app := app.NewApp(chainId, Logger)
		app.Unregister()
	},
}
