package cmd

import "github.com/spf13/cobra"

func NewCLI() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:          "cruz",
		Short:        "Santa Cruz CLI",
		SilenceUsage: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	rootCmd.AddCommand(
		addDbCmd(),
		listDbsCmd(),
		showDbCmd(),
		removeDbCmd(),
		editDbCmd(),
	)

	rootCmd.AddCommand(
		listTablesCmd(),
	)

	return rootCmd
}
