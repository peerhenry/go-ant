package chunks

import (
	"fmt"
	"strings"
)

func arrayToString(array interface{}) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(array)), ","), "[]")
}
