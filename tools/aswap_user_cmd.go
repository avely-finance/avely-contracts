package main

import (
	"os"

	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/core"
	. "github.com/avely-finance/avely-contracts/sdk/utils"

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
	DeadlineBlock:         &cli.IntFlag{Name: "deadline_block", Usage: "Deadline block for localnet, gap for other", Required: true},
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

			//token address
			tokenAddress := ctx.String(aswapUserCliFlags.TokenAddress.Name)
			if tokenAddress == "" {
				tokenAddress = aswapUser.config.StZilAddr
			}

			//deadline block
			deadlineBlock := ctx.Int(aswapUserCliFlags.DeadlineBlock.Name)
			if aswapUser.chain != "local" {
				//consider it as gap for testnet/mainnet
				deadlineBlock += aswapUser.sdk.GetBlockHeight()
			}

			_, err = aswapUser.aswap.AddLiquidity(
				ToQA(ctx.Int(aswapUserCliFlags.ZilAmount.Name)),
				tokenAddress,
				ToQA(ctx.Int(aswapUserCliFlags.MinContributionAmount.Name)),
				ToQA(ctx.Int(aswapUserCliFlags.MaxTokenAmount.Name)),
				deadlineBlock,
			)
			if err != nil {
				log.WithError(err).WithFields(logrus.Fields{"token_address": tokenAddress, "deadline_block": deadlineBlock}).Fatal("AddLiquidity error")
				return err
			}
			log.WithFields(logrus.Fields{"token_address": tokenAddress, "deadline_block": deadlineBlock}).Info("AddLiquidity succeed")
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

			//token address
			tokenAddress := ctx.String(aswapUserCliFlags.TokenAddress.Name)
			if tokenAddress == "" {
				tokenAddress = aswapUser.config.StZilAddr
			}

			//deadline block
			deadlineBlock := ctx.Int(aswapUserCliFlags.DeadlineBlock.Name)
			if aswapUser.chain != "local" {
				//consider it as gap for testnet/mainnet
				deadlineBlock += aswapUser.sdk.GetBlockHeight()
			}
			_, err = aswapUser.aswap.RemoveLiquidity(
				tokenAddress,
				ToQA(ctx.Int(aswapUserCliFlags.ContributionAmount.Name)),
				ToQA(ctx.Int(aswapUserCliFlags.MinZilAmount.Name)),
				ToQA(ctx.Int(aswapUserCliFlags.MinTokenAmount.Name)),
				deadlineBlock,
			)
			if err != nil {
				log.WithError(err).WithFields(logrus.Fields{"token_address": tokenAddress, "deadline_block": deadlineBlock}).Fatal("RemoveLiquidity error")
				return err
			}
			log.WithFields(logrus.Fields{"token_address": tokenAddress, "deadline_block": deadlineBlock}).Info("RemoveLiquidity succeed")
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

			//token address
			tokenAddress := ctx.String(aswapUserCliFlags.TokenAddress.Name)
			if tokenAddress == "" {
				tokenAddress = aswapUser.config.StZilAddr
			}

			//deadline block
			deadlineBlock := ctx.Int(aswapUserCliFlags.DeadlineBlock.Name)
			if aswapUser.chain != "local" {
				//consider it as gap for testnet/mainnet
				deadlineBlock += aswapUser.sdk.GetBlockHeight()
			}

			_, err = aswapUser.aswap.SwapExactZILForTokens(
				ToQA(ctx.Int(aswapUserCliFlags.ZilAmount.Name)),
				tokenAddress,
				ToQA(ctx.Int(aswapUserCliFlags.MinTokenAmount.Name)),
				ctx.String(aswapUserCliFlags.RecipientAddress.Name),
				deadlineBlock,
			)
			if err != nil {
				log.WithError(err).WithFields(logrus.Fields{"token_address": tokenAddress, "deadline_block": deadlineBlock}).Fatal("SwapExactZILForTokens error")
				return err
			}
			log.WithFields(logrus.Fields{"token_address": tokenAddress, "deadline_block": deadlineBlock}).Info("SwapExactZILForTokens succeed")
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

			//token address
			tokenAddress := ctx.String(aswapUserCliFlags.TokenAddress.Name)
			if tokenAddress == "" {
				tokenAddress = aswapUser.config.StZilAddr
			}

			//deadline block
			deadlineBlock := ctx.Int(aswapUserCliFlags.DeadlineBlock.Name)
			if aswapUser.chain != "local" {
				//consider it as gap for testnet/mainnet
				deadlineBlock += aswapUser.sdk.GetBlockHeight()
			}

			_, err = aswapUser.aswap.SwapExactTokensForZIL(
				tokenAddress,
				ToQA(ctx.Int(aswapUserCliFlags.TokenAmount.Name)),
				ToQA(ctx.Int(aswapUserCliFlags.MinZilAmount.Name)),
				ctx.String(aswapUserCliFlags.RecipientAddress.Name),
				deadlineBlock,
			)
			if err != nil {
				log.WithError(err).WithFields(logrus.Fields{"token_address": tokenAddress, "deadline_block": deadlineBlock}).Fatal("SwapExactTokensForZIL error")
				return err
			}
			log.WithFields(logrus.Fields{"token_address": tokenAddress, "deadline_block": deadlineBlock}).Info("SwapExactTokensForZIL succeed")
			return nil
		},
	}
}

func NewASwapUserCli(ctx *cli.Context) (*ASwapUserCli, error) {
	chain := ctx.String(aswapUserCliFlags.Chain.Name)

	log = NewLog()
	config := NewConfig(chain)
	sdk := NewAvelySDK(*config)

	aswapUserCli := &ASwapUserCli{
		sdk:    sdk,
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
