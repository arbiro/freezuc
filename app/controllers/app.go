package controllers

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/commands"
	"github.com/cosmos/cosmos-sdk/client/commands/query"
	"github.com/cosmos/cosmos-sdk/modules/coin"
	"github.com/cosmos/cosmos-sdk/stack"
	"github.com/revel/revel"
	"github.com/spf13/viper"
	"github.com/tendermint/tmlibs/cli"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Hello(myName string) revel.Result {
	c.Validation.Required(myName).Message("Your name is required!")
	c.Validation.MinSize(myName, 3).Message("Your name is not long enough!")

	homeDir := "C:\\Users\\Andrei\\.freecli"
	viper.Set(cli.HomeFlag, homeDir)

	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(homeDir)  // search root directory
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// stderr, so if we redirect output to json file, this doesn't appear
		// fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		// we ignore not found error, only parse error
		// stderr, so if we redirect output to json file, this doesn't appear
		fmt.Fprintf(os.Stderr, "%#v", err)
	}

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(App.Index)
	}
	log := c.Log.New("id", 2311)
	log.Debug("Reading the output")

	act, err := commands.ParseActor(myName)
	if err != nil {
		fmt.Print("Failed to load actor!!!")
		log.Errorf("Failed to load actor :%s", err.Error())
		return c.Redirect(App.Index)
	}

	act = coin.ChainAddr(act)
	key := stack.PrefixedKey(coin.NameCoin, act.Bytes())

	acc := coin.Account{}
	prove := !viper.GetBool(commands.FlagTrustNode)

	height, err := query.GetParsed(key, &acc, query.GetHeight(), prove)
	fmt.Printf("\nGetshere %d\n", height)
	if client.IsNoDataErr(err) {
		fmt.Printf("\nFUUUUUUUUUCJK %d\n", height)
		return c.Redirect(App.Index)
	} else if err != nil {
		return c.Redirect(App.Index)
	}
	fmt.Printf("\nGetshere %d\n", height)
	query.OutputProof(acc, height)

	return c.Render(myName)
}
