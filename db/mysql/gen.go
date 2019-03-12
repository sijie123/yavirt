package mysql

import (
	"fmt"
	"strings"
)

type generator struct{}

func newGenerator() *generator {
	return &generator{}
}

func (g *generator) genInsertx(table string, fields ...string) string {
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		table, g.genPlaces(true, fields...), g.genPlaces(false, fields...))
}

func (g *generator) genPlaces(name bool, fields ...string) string {
	var str strings.Builder
	var leng = len(fields)

	var prefix string
	if !name {
		prefix = ":"
	}

	for _, f := range fields[:leng-1] {
		str.WriteString(fmt.Sprintf("%s%s, ", prefix, f))
	}

	str.WriteString(fmt.Sprintf("%s%s", prefix, fields[leng-1]))

	return str.String()
}
