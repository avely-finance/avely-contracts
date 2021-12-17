package transitions

import (
	"Azil/test/helpers"
	"fmt"
)

const AZIL_SSN_ADDRESS = "0x166862bdd5d76b3a4775d2494820179d582acac5"
const AZIL_SSN_REWARD_SHARE_PERCENT = "50"

const HOLDER_INITIAL_DELEGATE_ZIL = 1000

// owner
const adminKey = "d96e9eb5b782a80ea153c937fa83e5948485fbfc8b7e7c069d7b914dbc350aba"
const admin = "381f4008505e940ad7681ec3468a719060caf796"

const verifierKey = "5430365143ce0154b682301d0ab731897221906a7054bbf5bd83c7663a6cbc40"
const verifier = "10200e3da08ee88729469d6eabc055cb225821e7"

const key1 = "1080d2cca18ace8225354ac021f9977404cee46f1d12e9981af8c36322eac1a4"
const addr1 = "ac941274c3b6a50203cc5e7939b7dad9f32a0c12"
const key2 = "254d9924fc1dcdca44ce92d80255c6a0bb690f867abde80e626fbfef4d357004"
const addr2 = "ec902fe17d90203d0bddd943d97b29576ece3177"
const key3 = "b8fc4e270594d87d3f728d0873a38fb0896ea83bd6f96b4f3c9ff0a29122efe4"
const addr3 = "c2035715831ab100ec42e562ce341b834bed1f4c"
const key4 = "b87f4ba7dcd6e60f2cca8352c89904e3993c5b2b0b608d255002edcda6374de4"
const addr4 = "6cd3667ba79310837e33f0aecbe13688a6cbca32"

const qa = "000000000000"

var shortcuts = map[string]string{
	"azilssn":  AZIL_SSN_ADDRESS,
	"addr1":    "0x" + addr1,
	"addr2":    "0x" + addr2,
	"addr3":    "0x" + addr3,
	"addr4":    "0x" + addr4,
	"admin":    "0x" + admin,
	"verifier": "0x" + verifier,
}
var t = helpers.NewTesting(shortcuts)

func zil(amount int) string {
	if amount == 0 {
		return "0"
	}
	return fmt.Sprintf("%d%s", amount, qa)
}

func azil(amount int) string {
	if amount == 0 {
		return "0"
	}
	return fmt.Sprintf("%d%s", amount, qa)
}

type Transitions struct {
}

func NewTransitions() *Transitions {
	return &Transitions{}
}
