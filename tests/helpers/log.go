package helpers

import (
	"fmt"
	"log"
)

func (t *Testing) LogStart(tag string) {
	log.Printf("⚙️  === Start to test %s === \n", tag)
}

func (t *Testing) LogEnd() {
	log.Println("🏁 TESTS PASSED SUCCESSFULLY")
}

func (t *Testing) LogError(tag string, err error) {
	log.Fatalf("🔴 Failed at %s, err = %s\n", tag, err.Error())
}

func (t *Testing) LogState(state string) {
	result := t.HighlightShortcuts(state)
	fmt.Println(result)
}
