package sqlparser

import (
	"fmt"
	"github.com/brunowang/gframe/cmd/dao-gen/gen/template"
	"github.com/brunowang/gframe/cmd/protoc-gen-go-gframe/gen/helper"
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/test_driver"
	"reflect"
)

type ColumnDef struct {
	Name    string
	Type    string
	Comment string
}

type IndexDef struct {
	Type string
	Cols []string
}

type TableDef struct {
	Name    string
	Cols    []ColumnDef
	Idxs    []IndexDef
	allType map[string]struct{}
}

func (v *TableDef) Enter(in ast.Node) (ast.Node, bool) {
	switch node := in.(type) {
	case *ast.TableName:
		v.Name = node.Name.O
	case *ast.ColumnDef:
		col := ColumnDef{
			Name: node.Name.Name.O,
			Type: typeMap[node.Tp.GetType()],
		}
		for _, o := range node.Options {
			if o.Tp == ast.ColumnOptionComment {
				if expr, ok := o.Expr.(*test_driver.ValueExpr); ok {
					col.Comment = expr.Datum.GetString()
				}
			}
		}
		v.Cols = append(v.Cols, col)
	case *ast.Constraint:
		// TODO handle index
		fmt.Printf("Constraint: %+v", node)
	default:
		v.allType[fmt.Sprintf("%+v", reflect.TypeOf(node))] = struct{}{}
	}
	return in, false
}

func (v *TableDef) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func ParseTable(sql string) ([]*template.StructDefTmpl, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}
	tpls := make([]*template.StructDefTmpl, 0, len(stmtNodes))
	for _, stmt := range stmtNodes {
		tab := &TableDef{allType: make(map[string]struct{})}
		stmt.Accept(tab)
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
		tpls = append(tpls, tpl)
	}
	return tpls, nil
}
