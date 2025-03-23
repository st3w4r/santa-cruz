package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/st3w4r/santa-cruz/manager"
	"github.com/st3w4r/santa-cruz/server"
)

func serve(cmd *cobra.Command, args []string) error {
	host := cmd.Flag("host").Value.String()
	port := cmd.Flag("port").Value.String()

	controlDb, err := manager.InitDBManger()
	if err != nil {
		return err
	}
	dataDbs := []manager.ManagedDb{}

	fmt.Println("Starting server on " + host + ":" + port)
	server := server.NewServer(controlDb, dataDbs)
	err = server.Serve(host, port)
	if err != nil {
		return err
	}

	return nil
}

func serveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve --host <host> --port <port>",
		Short: "Start the server",
		RunE:  serve,
	}
	cmd.Flags().StringP("host", "H", "localhost", "Host to bind to")
	cmd.Flags().StringP("port", "p", "8080", "Port to bind to")
	return cmd
}
