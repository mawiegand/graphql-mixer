package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nautilus/gateway"
	"github.com/nautilus/graphql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the gateway server",
	Run:   Run,
}

func init() {
	serverCmd.Flags().StringP("host", "H", "localhost", "Bind socket to this host.")
	if err := viper.BindPFlag("host", serverCmd.Flags().Lookup("host")); err != nil {
		fmt.Println("Could not bind host flag:", err.Error())
		os.Exit(1)
	}

	serverCmd.Flags().StringP("port", "P", "4000", "Bind socket to this port.")
	if err := viper.BindPFlag("port", serverCmd.Flags().Lookup("port")); err != nil {
		fmt.Println("Could not bind port flag:", err.Error())
		os.Exit(1)
	}

	serverCmd.Flags().StringSliceP("services", "s", []string{}, "Services to be federated.")
	if err := viper.BindPFlag("services", serverCmd.Flags().Lookup("services")); err != nil {
		fmt.Println("Could not bind services flag:", err.Error())
		os.Exit(1)
	}

	rootCmd.AddCommand(serverCmd)
}

func Run(cmd *cobra.Command, args []string) {
	// get args from viper
	host := viper.GetString("host")
	port := viper.GetString("port")
	services := viper.GetStringSlice("services")

	if len(services) == 0 {
		fmt.Println("Please specify services for schema introspection!")
		os.Exit(1)
	}

	// introspect schemas
	schemas, err := graphql.IntrospectRemoteSchemas(services...)
	if err != nil {
		fmt.Println("Could not introspect schemas:", err.Error())
		os.Exit(1)
	}

	// create a gateway instance
	gw, err := gateway.New(schemas)
	if err != nil {
		fmt.Println("Could not create gateway:", err.Error())
		os.Exit(1)
	}

	// add graphql endpoints
	http.HandleFunc("/graphql", gw.PlaygroundHandler)

	// serve the HTTP server
	fmt.Printf("🚀 GraphQL-Mixer Gateway is ready at http://%s:%s/graphql\n", host, port)
	err = http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), nil)
	if err != nil {
		fmt.Println("Could not run server", err.Error())
		os.Exit(1)
	}
}
