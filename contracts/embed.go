package contracts

import "embed"

//go:embed source/*
var contractFs embed.FS

func GetContractFs() embed.FS {
	return contractFs
}
