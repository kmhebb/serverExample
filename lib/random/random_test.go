// +build !integration

package random_test

import (
	"testing"

	"github.com/kmhebb/serverExample/lib/random"
)

func TestWords(t *testing.T) {
	var prev string

	for i := 0; i < 10; i++ {
		curr := random.Words(random.PhraseLength)
		if curr == prev {
			t.Errorf("probably shouldn't get same words twice: %s == %s\n", curr, prev)
		} else {
			t.Logf("Curr: %s\n", curr)
		}
	}
}
