package cmd

import "github.com/spf13/cobra"

func NewCLI() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:          "cruz",
		Short:        "Santa Cruz CLI",
		SilenceUsage: false,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
	// Manage databases
	rootCmd.AddCommand(
		addDbCmd(),
		listDbsCmd(),
		showDbCmd(),
		removeDbCmd(),
		editDbCmd(),
	)
	// Operations on managed databases
	rootCmd.AddCommand(
		listTablesCmd(),
		queryCmd(),
	)

	return rootCmd
}
