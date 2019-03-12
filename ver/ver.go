package ver

import "fmt"

var Git, Compile, Date string

func Version() string {
	return fmt.Sprintf(`Git: %s
Compile: %s
Built: %s`, Git, Compile, Date)
}
