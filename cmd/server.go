package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nautilus/gateway"
	"github.com/nautilus/graphql"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the gateway server",
	Run:   Run,
}

var Host string
var Port string
var Services []string

func init() {
	serverCmd.Flags().StringVarP(&Host, "host", "H", "localhost", "Bind socket to this host.")

	serverCmd.Flags().StringVarP(&Port, "port", "P", "4000", "Bind socket to this port.")

	serverCmd.Flags().StringSliceVarP(&Services, "services", "s", []string{}, "Services to be federated.")
	serverCmd.MarkFlagRequired("services")

	rootCmd.AddCommand(serverCmd)
}

func Run(cmd *cobra.Command, args []string) {
	// introspect schemas
	schemas, err := graphql.IntrospectRemoteSchemas(Services...)
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
	fmt.Printf("🚀 GraphQL-Mixer Gateway is ready at http://%s:%s/graphql\n", Host, Port)
	err = http.ListenAndServe(fmt.Sprintf("%s:%s", Host, Port), nil)
	if err != nil {
		fmt.Println("Could not run server", err.Error())
		os.Exit(1)
	}
}
