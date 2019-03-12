package cmd

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
)

func Prof(port int) {
	var enable = strings.ToLower(os.Getenv("YAVIRTD_PPROF"))
	switch enable {
	case "":
		fallthrough
	case "0":
		fallthrough
	case "false":
		fallthrough
	case "off":
		http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
	default:
		return
	}
}
