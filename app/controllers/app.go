package controllers

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/commands"
	"github.com/cosmos/cosmos-sdk/client/commands/query"
	"github.com/cosmos/cosmos-sdk/modules/coin"
	"github.com/cosmos/cosmos-sdk/stack"
	"github.com/revel/revel"
	"github.com/spf13/viper"
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

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(App.Index)
	}
	log := c.Log.New("id", 2311)
	log.Debug("Reading the output")

	act, err := commands.ParseActor(myName)
	if err != nil {
		log.Errorf("Failed to load actor :%s", err.Error())
		return c.Redirect(App.Index)
	}
	act = coin.ChainAddr(act)
	key := stack.PrefixedKey(coin.NameCoin, act.Bytes())

	acc := coin.Account{}
	prove := !viper.GetBool(commands.FlagTrustNode)
	height, err := query.GetParsed(key, &acc, query.GetHeight(), prove)
	if client.IsNoDataErr(err) {
		return c.Redirect(App.Index)
	} else if err != nil {
		return c.Redirect(App.Index)
	}

	query.OutputProof(acc, height)

	return c.Render(myName)
}
