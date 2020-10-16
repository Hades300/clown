package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "clown",
	Short: "clown is a arp spoofing tool which work in len env",
	Long: `clown is a arp spoofing tool which work in len env
			Check https://github/hades300/clown for more info`,
}

func init() {
	rootCmd.AddCommand(pretendCmd)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
