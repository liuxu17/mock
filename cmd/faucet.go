package cmd

import (
	"fmt"
	"github.com/kaifei-bianjie/mock/conf"
	"github.com/kaifei-bianjie/mock/util/constants"
	"github.com/kaifei-bianjie/mock/util/helper"
	"github.com/kaifei-bianjie/mock/util/helper/account"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func FaucetInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "faucet-init",
		Short: "init mock faucet account",
		Long: `init mock faucet account
Note the account must has many token, so that this account can transfer token to other account.
Example:
	mock faucet-init --seed="recycle light kid ..."
`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			seed := viper.GetString(FlagFaucetSeed)
			address, err := account.CreateAccount(constants.MockFaucetName, constants.MockFaucetPassword, seed)

			if err != nil {
				return err
			}

			fmt.Printf("faucet address is: %v\n", address)

			// TODO: can't read default value
			homeDir := viper.GetString(FlagHome)
			confFilePath := fmt.Sprintf("%v/%v", homeDir, constants.ConfigFileName)
			helper.WriteFile(confFilePath, []byte(fmt.Sprintf("%v=%v\n", conf.KeyFaucetSeed, seed)))
			helper.WriteFile(confFilePath, []byte(fmt.Sprintf("%v=%v\n", conf.KeyFaucetAddress, address)))

			return nil
		},
	}

	cmd.Flags().AddFlagSet(faucetFlagSet)

	cmd.MarkFlagRequired(FlagFaucetSeed)

	return cmd
}
