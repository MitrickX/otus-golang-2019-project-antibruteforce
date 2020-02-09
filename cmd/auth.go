package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/grpc"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/domain/entities"

	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth <login> <password> <ip> [flags]",
	Short: "Auth",
	Long: `Check that auth is allowed for login, password and ip
All arguments are required
IP should be of host

See examples below:
	auth test 1234 193.70.18.123 - check that auth is allowed for login=test, password=1234 and ip=192.70.18.123
`,
	DisableFlagsInUseLine: true,
	Args: func(cmd *cobra.Command, args []string) error {
		return validateAuthCmdArgs(args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		runAuthCommand(args[0], args[1], args[2])
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}

func validateAuthCmdArgs(args []string) error {
	//nolint:gomnd
	if len(args) < 1 {
		return errors.New("<login>, <password> and <ip> is required. Run with --help for more information")
	}

	//nolint:gomnd
	if len(args) < 2 {
		return errors.New("<password> and <ip> is required. Run with --help for more information")
	}

	//nolint:gomnd
	if len(args) < 3 {
		return errors.New("<ip> is required. Run with --help for more information")
	}

	_, err := entities.NewWithoutMaskPart(args[2])
	if err != nil {
		return fmt.Errorf("%s. Run with --help for more information", err)
	}

	return nil
}

func runAuthCommand(login, password, ip string) {
	cfg := getGRPCClientConfig()

	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()

	result, err := cfg.apiClient.Auth(ctx, &grpc.AuthRequest{
		Login:    login,
		Password: password,
		Ip:       ip,
	})

	if err != nil {
		fmt.Printf("FAIL: %s\n", err)
		return
	}

	if result.Ok {
		fmt.Printf("OK: auth for login=%s, password=%s, ip=%s is allowed\n", login, password, ip)
	} else {
		fmt.Printf("NOT OK: auth for login=%s, password=%s, ip=%s is NOT allowed\n", login, password, ip)
	}
}
