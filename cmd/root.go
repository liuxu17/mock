package cmd

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "mock",
		Short: "mock test",
	}

	rootCmd.PersistentFlags().AddFlagSet(rootFlagSet)

	rootCmd.MarkFlagRequired(FlagChainId)
	rootCmd.MarkFlagRequired(FlagNodeUrl)

	return rootCmd
}
