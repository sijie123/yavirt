package mysql

import (
	"fmt"
	"strings"
)

type Dsn struct {
	User, Password, Addr, Proto, DB string
}

func (d *Dsn) String() string {
	var str strings.Builder

	if len(d.User) > 0 {
		str.WriteString(fmt.Sprintf("%s:%s", d.User, d.Password))
	}

	str.WriteString(d.Proto)

	if len(d.Addr) > 0 {
		str.WriteString(fmt.Sprintf("(%s)", d.Addr))
	}

	str.WriteString(fmt.Sprintf("/%s", d.DB))

	return str.String()
}
