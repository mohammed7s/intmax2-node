package balance_checker

import (
	"context"
	"intmax2-node/configs"
	"intmax2-node/internal/logger"

	"github.com/spf13/cobra"
)

type Balance struct {
	Context context.Context
	Config  *configs.Config
	Log     logger.Logger
	SB      ServiceBlockchain
}

func NewBalanceCmd(b *Balance) *cobra.Command {
	const (
		use   = "balance"
		short = "Manage balance"
	)

	balanceCmd := &cobra.Command{
		Use:   use,
		Short: short,
	}
	balanceCmd.AddCommand(getBalanceCmd(b))

	return balanceCmd
}

func getBalanceCmd(b *Balance) *cobra.Command {
	const (
		use                    = "get"
		short                  = "Get balance of specified INTMAX account"
		userPrivateKeyKey      = "private-key"
		emptyKey               = ""
		userAddressDescription = "specify user address. use as --private-key \"0x0000000000000000000000000000000000000000000000000000000000000000\""
	)

	cmd := cobra.Command{
		Use:   use,
		Short: short,
	}

	var userEthPrivateKey string
	cmd.PersistentFlags().StringVar(&userEthPrivateKey, userPrivateKeyKey, emptyKey, userAddressDescription)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		l := b.Log.WithFields(logger.Fields{"module": use})

		err := newCommands().GetBalance(b.Config, b.Log, b.SB).Do(b.Context, args, userEthPrivateKey)
		if err != nil {
			const msg = "failed to get balance: %v"
			l.Fatalf(msg, err.Error())
		}
	}

	return &cmd
}
