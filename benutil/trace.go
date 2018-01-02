package benutil

import (
	"fmt"
)

func Debug(format string, values ...interface{}) {
	fmt.Printf("+++<ben>:"+format, values...)
}
