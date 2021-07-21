package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "graphql-mixer",
	Short: "A gateway for GraphQL APIs.",
	Long: "GraphQL Mixer is a gateway to federate GraphQL APIs and serve the result via HTTP.\n" +
		"It's based on https://github.com/nautilus/gateway.",
}

var configFile string

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(
		&configFile,
		"config",
		"",
		"Path to config file (default search paths are /etc/graphql-mixer/, "+
			"$HOME/.graphql-mixer, working directory)")
}

func initConfig() {
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")               // name of config file (without extension)
		viper.SetConfigType("yaml")                 // REQUIRED if the config file does not have the extension in the name
		viper.AddConfigPath("/etc/graphql-mixer/")  // path to look for the config file in
		viper.AddConfigPath("$HOME/.graphql-mixer") // call multiple times to add many search paths
		viper.AddConfigPath(".")                    // optionally look for config in the working directory
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Printf("Running without config file!\n")
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("Fatal error config file: %w \n", err))
		}
	} else {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
