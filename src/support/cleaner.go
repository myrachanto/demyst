package support

import (
	"os"
)

var Clean clean

type clean struct{}

func (c clean) Cleaner(filename string) {
	if filename != "" && filename != "undefined" {
		os.Remove("./src/public" + filename)
	}
}
