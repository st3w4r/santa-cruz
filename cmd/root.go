package cmd

import "github.com/spf13/cobra"

func NewCLI() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "cruz",
		Short: "Santa Cruz CLI",
	}

	rootCmd.AddCommand(
		listDbsCmd(),
		addDbCmd(),
	)

	return rootCmd
}
