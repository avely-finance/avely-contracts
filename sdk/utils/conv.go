package utils

import (
	"strconv"
)

func ArrayAtoi(input []string) []int {
	out := []int{}
	for _, val := range input {
		res, _ := strconv.Atoi(val)
		out = append(out, res)
	}
	return out
}

func ArrayItoa(input []int) []string {
	out := []string{}
	for _, val := range input {
		out = append(out, strconv.Itoa(val))
	}
	return out
}
