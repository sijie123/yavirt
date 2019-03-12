package template

import (
	"bytes"
	"strings"
	"sync"
	text "text/template"

	"github.com/juju/errors"

	"github.com/projecteru2/yavirt/util"
)

func Render(filepath string, args interface{}) ([]byte, error) {
	var tmpl, err = templates.get(filepath)
	if err != nil {
		return nil, errors.Errorf("no such template: %s", filepath)
	}

	var wr bytes.Buffer
	if err := tmpl.Execute(&wr, args); err != nil {
		return nil, errors.Trace(err)
	}

	return wr.Bytes(), nil
}

var templates Templates

type Templates struct {
	sync.Map
}

func (t *Templates) get(fpth string) (tmpl *text.Template, err error) {
	var val, ok = t.Load(fpth)
	if ok {
		return val.(*text.Template), nil
	}

	tmpl, err = t.parse(fpth)
	if err != nil {
		return nil, errors.Trace(err)
	}

	t.Store(fpth, tmpl)

	return tmpl, nil
}

func (t *Templates) parse(fpth string) (*text.Template, error) {
	if !t.isTempl(fpth) {
		return nil, errors.Errorf("%s is not a template file", fpth)
	}

	buf, err := util.ReadAll(fpth)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return text.New(fpth).Parse(string(buf))
}

func (t *Templates) isTempl(fpth string) bool {
	return strings.HasSuffix(fpth, ".xml")
}
