package cmd

import (
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <kind> <ip> [flags]",
	Short: "Delete IP from black or white list",
	Long: `Delete IP from black or white list. 
IP could be as of host or as of subnet

<kind> = black | white
<ip> = host IP or subnet IP in CIDR notation

See examples below:
	delete black 191.80.28.0/24	- delete subnet IP from black list 
	delete black 128.73.19.123	- delete host IP from black list
	delete white 280.71.28.0/24	- delete subnet IP from white list
	delete white 196.60.28.123	- delete host IP from white list
`,
	DisableFlagsInUseLine: true,
	Args: func(cmd *cobra.Command, args []string) error {
		return validateListCmdArgs(args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		runDeleteCommand(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
