package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/tendermint/tmlibs/cli"

	"github.com/cosmos/cosmos-sdk/client/commands"
	"github.com/cosmos/cosmos-sdk/client/commands/auto"
	"github.com/cosmos/cosmos-sdk/client/commands/commits"
	"github.com/cosmos/cosmos-sdk/client/commands/keys"
	"github.com/cosmos/cosmos-sdk/client/commands/proxy"
	"github.com/cosmos/cosmos-sdk/client/commands/query"
	rpccmd "github.com/cosmos/cosmos-sdk/client/commands/rpc"
	"github.com/cosmos/cosmos-sdk/client/commands/search"
	txcmd "github.com/cosmos/cosmos-sdk/client/commands/txs"
	authcmd "github.com/cosmos/cosmos-sdk/modules/auth/commands"
	basecmd "github.com/cosmos/cosmos-sdk/modules/base/commands"
	coincmd "github.com/cosmos/cosmos-sdk/modules/coin/commands"
	feecmd "github.com/cosmos/cosmos-sdk/modules/fee/commands"
	ibccmd "github.com/cosmos/cosmos-sdk/modules/ibc/commands"
	noncecmd "github.com/cosmos/cosmos-sdk/modules/nonce/commands"
	rolecmd "github.com/cosmos/cosmos-sdk/modules/roles/commands"
)

// FreeCli - main basecoin client command
var FreeCli = &cobra.Command{
	Use:   "freecli",
	Short: "Light client for Tendermint",
	Long: `Freecli is a certifying light client for the freecoin abci app.

It leverages the power of the tendermint consensus algorithm get full
cryptographic proof of all queries while only syncing a fraction of the
block headers.`,
}

func main() {
	commands.AddBasicFlags(FreeCli)

	// Prepare queries
	query.RootCmd.AddCommand(
		// These are default parsers, but optional in your app (you can remove key)
		query.TxQueryCmd,
		query.KeyQueryCmd,
		coincmd.AccountQueryCmd,
		noncecmd.NonceQueryCmd,
		rolecmd.RoleQueryCmd,
		ibccmd.IBCQueryCmd,
	)

	// these are queries to search for a tx
	search.RootCmd.AddCommand(
		coincmd.SentSearchCmd,
	)

	// set up the middleware
	txcmd.Middleware = txcmd.Wrappers{
		feecmd.FeeWrapper{},
		rolecmd.RoleWrapper{},
		noncecmd.NonceWrapper{},
		basecmd.ChainWrapper{},
		authcmd.SigWrapper{},
	}
	txcmd.Middleware.Register(txcmd.RootCmd.PersistentFlags())

	// you will always want this for the base send command
	txcmd.RootCmd.AddCommand(
		// This is the default transaction, optional in your app
		coincmd.SendTxCmd,
		coincmd.CreditTxCmd,
		// this enables creating roles
		rolecmd.CreateRoleTxCmd,
		// these are for handling ibc
		ibccmd.RegisterChainTxCmd,
		ibccmd.UpdateChainTxCmd,
		ibccmd.PostPacketTxCmd,
	)

	// Set up the various commands to use
	FreeCli.AddCommand(
		commands.InitCmd,
		commands.ResetCmd,
		keys.RootCmd,
		commits.RootCmd,
		rpccmd.RootCmd,
		query.RootCmd,
		search.RootCmd,
		txcmd.RootCmd,
		proxy.RootCmd,
		commands.VersionCmd,
		auto.AutoCompleteCmd,
	)

	cmd := cli.PrepareMainCmd(FreeCli, "FC", os.ExpandEnv("$HOME/.freecli"))
	cmd.Execute()
}
