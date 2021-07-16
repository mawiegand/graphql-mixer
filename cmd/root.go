package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "graphql-mixer",
	Short: "A gateway for GraphQL APIs.",
	Long: `GraphQL Mixer is a gateway to federate GraphQL API
				and serve the result via HTTP.
				It's based on https://github.com/nautilus/gateway.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
