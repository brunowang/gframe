package template

import (
	_ "embed"
	"github.com/brunowang/gframe/cmd/dao-gen/gen/helper"
)

type MgrDefTmpl struct {
	Stmts []CreateStmt
}

type CreateStmt struct {
	Name     string
	StmtName string
	OrigStmt string
}

//go:embed mgrdef.tmpl
var mgrdeftmpl string

func (t *MgrDefTmpl) Render() string {
	return helper.NewTmplRenderer("mgrdef.tmpl").
		Text(mgrdeftmpl).Data(t).RenderTmpl()
}

const sqldeftmpl = `package dao
{{range .Stmts}}
const {{.StmtName}} = ` +
	"`\n{{.OrigStmt}}\n`" + `
{{end}}`

func (t *MgrDefTmpl) RenderSQL() string {
	return helper.NewTmplRenderer("sqldef.tmpl").
		Text(sqldeftmpl).Data(t).RenderTmpl()
}
