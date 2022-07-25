package main

import (
	"os"

	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/core"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type ASwapUserCliFlags struct {
	Chain                 *cli.StringFlag
	TokenAddress          *cli.StringFlag
	RecipientAddress      *cli.StringFlag
	ContributionAmount    *cli.IntFlag
	MinContributionAmount *cli.IntFlag
	MinZilAmount          *cli.IntFlag
	ZilAmount             *cli.IntFlag
	TokenAmount           *cli.IntFlag
	MaxTokenAmount        *cli.IntFlag
	MinTokenAmount        *cli.IntFlag
	DeadlineBlock         *cli.IntFlag
}

type ASwapUserCli struct {
	sdk    *AvelySDK
	config *Config
	chain  string
	aswap  *ASwap
}

var log *Log
var aswapUserCliFlags = ASwapUserCliFlags{
	Chain:                 &cli.StringFlag{Name: "chain", Usage: "Chain", Value: "local"},
	TokenAddress:          &cli.StringFlag{Name: "token_address", Usage: "Token address, StZil by default"},
	RecipientAddress:      &cli.StringFlag{Name: "recipient_address", Usage: "Recipient address", Required: true},
	ContributionAmount:    &cli.IntFlag{Name: "contribution_amount", Usage: "Contribution amount", Required: true},
	MinContributionAmount: &cli.IntFlag{Name: "min_contribution_amount", Usage: "Minimum contribution amount", Required: true},
	MinZilAmount:          &cli.IntFlag{Name: "min_zil_amount", Usage: "Minimum ZIL amount", Required: true},
	ZilAmount:             &cli.IntFlag{Name: "zil_amount", Usage: "Native ZIL _amount to send", Required: true},
	TokenAmount:           &cli.IntFlag{Name: "token_amount", Usage: "Token amount", Required: true},
	MaxTokenAmount:        &cli.IntFlag{Name: "max_token_amount", Usage: "Maximum token amount", Required: true},
	MinTokenAmount:        &cli.IntFlag{Name: "min_token_amount", Usage: "Minimum token amount", Required: true},
	DeadlineBlock:         &cli.IntFlag{Name: "deadline_block", Usage: "Deadline block", Required: true},
}

func main() {

	logrus.SetLevel(logrus.DebugLevel)

	app := &cli.App{
		EnableBashCompletion: true,
		Name:                 "aswap-user-cmd",
		Usage:                "ASwap user commands",
		Commands: []*cli.Command{
			addLiquidity(),
			removeLiquidity(),
			swapExactZilForTokens(),
			swapExactTokensForZil(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.WithError(err).Fatal("could not run application")
	}
}

func addLiquidity() *cli.Command {
	return &cli.Command{
		Name:  "add_liquidity",
		Usage: "Adds liquidity",
		Flags: []cli.Flag{
			aswapUserCliFlags.Chain,
			aswapUserCliFlags.ZilAmount,
			aswapUserCliFlags.TokenAddress,
			aswapUserCliFlags.MinContributionAmount,
			aswapUserCliFlags.MaxTokenAmount,
			aswapUserCliFlags.DeadlineBlock,
		},
		Action: func(ctx *cli.Context) error {
			aswapUser, err := NewASwapUserCli(ctx)
			if err != nil {
				log.WithError(err).Fatal("Can't initialize ASwap CLI")
				return err
			}
			tokenAddress := ctx.String(aswapUserCliFlags.TokenAddress.Name)
			if tokenAddress == "" {
				tokenAddress = aswapUser.config.StZilAddr
			}
			_, err = aswapUser.aswap.AddLiquidity(
				ctx.String(aswapUserCliFlags.ZilAmount.Name),
				tokenAddress,
				ctx.String(aswapUserCliFlags.MinContributionAmount.Name),
				ctx.String(aswapUserCliFlags.MaxTokenAmount.Name),
				ctx.Int(aswapUserCliFlags.DeadlineBlock.Name),
			)
			if err != nil {
				log.WithError(err).WithFields(logrus.Fields{"token_address": tokenAddress}).Fatal("AddLiquidity error")
				return err
			}
			log.WithFields(logrus.Fields{"token_address": tokenAddress}).Info("AddLiquidity succeed")
			return nil
		},
	}
}

func removeLiquidity() *cli.Command {
	return &cli.Command{
		Name:  "remove_liquidity",
		Usage: "Removes liquidity",
		Flags: []cli.Flag{
			aswapUserCliFlags.Chain,
			aswapUserCliFlags.TokenAddress,
			aswapUserCliFlags.ContributionAmount,
			aswapUserCliFlags.MinZilAmount,
			aswapUserCliFlags.MinTokenAmount,
			aswapUserCliFlags.DeadlineBlock,
		},
		Action: func(ctx *cli.Context) error {
			aswapUser, err := NewASwapUserCli(ctx)
			if err != nil {
				log.WithError(err).Fatal("Can't initialize ASwap CLI")
				return err
			}
			tokenAddress := ctx.String(aswapUserCliFlags.TokenAddress.Name)
			if tokenAddress == "" {
				tokenAddress = aswapUser.config.StZilAddr
			}
			_, err = aswapUser.aswap.RemoveLiquidity(
				tokenAddress,
				ctx.String(aswapUserCliFlags.ContributionAmount.Name),
				ctx.String(aswapUserCliFlags.MinZilAmount.Name),
				ctx.String(aswapUserCliFlags.MinTokenAmount.Name),
				ctx.Int(aswapUserCliFlags.DeadlineBlock.Name),
			)
			if err != nil {
				log.WithError(err).WithFields(logrus.Fields{"token_address": tokenAddress}).Fatal("RemoveLiquidity error")
				return err
			}
			log.WithFields(logrus.Fields{"token_address": tokenAddress}).Info("RemoveLiquidity succeed")
			return nil
		},
	}
}

func swapExactZilForTokens() *cli.Command {
	return &cli.Command{
		Name:  "swap_exact_zil_for_tokens",
		Usage: "Swaps exact zil for tokens",
		Flags: []cli.Flag{
			aswapUserCliFlags.Chain,
			aswapUserCliFlags.ZilAmount,
			aswapUserCliFlags.TokenAddress,
			aswapUserCliFlags.MinTokenAmount,
			aswapUserCliFlags.DeadlineBlock,
			aswapUserCliFlags.RecipientAddress,
		},
		Action: func(ctx *cli.Context) error {
			aswapUser, err := NewASwapUserCli(ctx)
			if err != nil {
				log.WithError(err).Fatal("Can't initialize ASwap CLI")
				return err
			}
			tokenAddress := ctx.String(aswapUserCliFlags.TokenAddress.Name)
			if tokenAddress == "" {
				tokenAddress = aswapUser.config.StZilAddr
			}
			_, err = aswapUser.aswap.SwapExactZILForTokens(
				ctx.String(aswapUserCliFlags.ZilAmount.Name),
				tokenAddress,
				ctx.String(aswapUserCliFlags.MinTokenAmount.Name),
				ctx.String(aswapUserCliFlags.RecipientAddress.Name),
				ctx.Int(aswapUserCliFlags.DeadlineBlock.Name),
			)
			if err != nil {
				log.WithError(err).WithFields(logrus.Fields{"token_address": tokenAddress}).Fatal("SwapExactZILForTokens error")
				return err
			}
			log.WithFields(logrus.Fields{"token_address": tokenAddress}).Info("SwapExactZILForTokens succeed")
			return nil
		},
	}
}

func swapExactTokensForZil() *cli.Command {
	return &cli.Command{
		Name:  "swap_exact_tokens_for_zil",
		Usage: "Swaps exact tokens for zil",
		Flags: []cli.Flag{
			aswapUserCliFlags.Chain,
			aswapUserCliFlags.TokenAddress,
			aswapUserCliFlags.TokenAmount,
			aswapUserCliFlags.MinZilAmount,
			aswapUserCliFlags.DeadlineBlock,
			aswapUserCliFlags.RecipientAddress,
		},
		Action: func(ctx *cli.Context) error {
			aswapUser, err := NewASwapUserCli(ctx)
			if err != nil {
				log.WithError(err).Fatal("Can't initialize ASwap CLI")
				return err
			}
			tokenAddress := ctx.String(aswapUserCliFlags.TokenAddress.Name)
			if tokenAddress == "" {
				tokenAddress = aswapUser.config.StZilAddr
			}
			_, err = aswapUser.aswap.SwapExactTokensForZIL(
				tokenAddress,
				ctx.String(aswapUserCliFlags.TokenAmount.Name),
				ctx.String(aswapUserCliFlags.MinZilAmount.Name),
				ctx.String(aswapUserCliFlags.RecipientAddress.Name),
				ctx.Int(aswapUserCliFlags.DeadlineBlock.Name),
			)
			if err != nil {
				log.WithError(err).WithFields(logrus.Fields{"token_address": tokenAddress}).Fatal("SwapExactTokensForZIL error")
				return err
			}
			log.WithFields(logrus.Fields{"token_address": tokenAddress}).Info("SwapExactTokensForZIL succeed")
			return nil
		},
	}
}

func NewASwapUserCli(ctx *cli.Context) (*ASwapUserCli, error) {
	chain := ctx.String(aswapUserCliFlags.Chain.Name)

	log = NewLog()
	config := NewConfig(chain)

	aswapUserCli := &ASwapUserCli{
		sdk:    NewAvelySDK(*config),
		config: config,
		chain:  chain,
	}
	aswapUserCli.aswap = aswapUserCli.restoreASwap()

	return aswapUserCli, nil
}

func (usercli *ASwapUserCli) restoreASwap() *ASwap {
	aswap, err := RestoreASwap(usercli.sdk, usercli.config.ASwapAddr, "")
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{}).Fatal("Can't restore ASwap contract")
	}
	log.WithFields(logrus.Fields{}).Info("ASwap contract restored")
	aswap.UpdateWallet(usercli.config.OwnerKey)
	return aswap
}
