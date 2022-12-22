package gen

import (
	"github.com/brunowang/gframe/cmd/dao-gen/gen/helper"
	"github.com/brunowang/gframe/cmd/dao-gen/gen/sqlparser"
	"github.com/brunowang/gframe/cmd/dao-gen/gen/template"
	"os"
	"strings"
)

func GenerateDAO(tabs []sqlparser.TableDef) error {
	tpl := &template.MgrDefTmpl{
		Stmts: make([]template.CreateStmt, 0, len(tabs)),
	}
	for _, tab := range tabs {
		stmt := template.CreateStmt{
			Name:     helper.ToCamelCase(tab.Name),
			StmtName: helper.UnTitle(helper.ToCamelCase(tab.Name)) + "Table",
			OrigStmt: strings.ReplaceAll(tab.Stmt, "`", ""),
		}
		tpl.Stmts = append(tpl.Stmts, stmt)
	}
	f, err := os.Create("tbl.go")
	if err != nil {
		return err
	}
	_, err = f.WriteString(tpl.RenderSQL())
	if err != nil {
		return err
	}

	f, err = os.Create("dao.go")
	if err != nil {
		return err
	}
	_, err = f.WriteString(tpl.Render() + "\n")
	if err != nil {
		return err
	}
	return nil
}

func GenerateCURD(tab sqlparser.TableDef) error {
	tpl := &template.StructDefTmpl{
		Name:    helper.ToCamelCase(tab.Name),
		TabName: tab.Name,
	}
	for _, col := range tab.Cols {
		tpl.Fields = append(tpl.Fields, template.Field{
			Name:    helper.ToCamelCase(col.Name),
			Type:    col.Type,
			ColName: col.Name,
			Comment: col.Comment,
		})
	}
	for _, idx := range tab.Idxs {
		fields := make([]template.Field, 0, len(idx.Cols))
		for _, col := range idx.Cols {
			field, err := tpl.Fields.Find(helper.ToCamelCase(col))
			if err != nil {
				field = template.Field{
					Name:    "<unknown>",
					Type:    "<unknown>",
					ColName: "<unknown>",
				}
			}
			field.Name = helper.UnTitle(field.Name)
			fields = append(fields, field)
		}
		tpl.Indexes = append(tpl.Indexes, template.Index{
			Uniq: idx.Uniq,
			Cols: fields,
		})
	}
	f, err := os.Create(tpl.TabName + ".go")
	if err != nil {
		return err
	}
	_, err = f.WriteString(tpl.Render())
	if err != nil {
		return err
	}
	return nil
}
