package sqlparser

import (
	"fmt"
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/test_driver"
	"reflect"
	"strings"
)

type ColumnDef struct {
	Name    string
	Type    string
	ZeroVal string
	Comment string
}

type IndexDef struct {
	Uniq bool
	Cols []string
}

type TableDef struct {
	Name    string
	Stmt    string
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
			Name:    node.Name.Name.O,
			Type:    typeMap.GetType(node.Tp.GetType()),
			ZeroVal: typeMap.GetZeroVal(node.Tp.GetType()),
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
		idx := IndexDef{}
		if node.Tp == ast.ConstraintPrimaryKey || node.Tp == ast.ConstraintUniq ||
			node.Tp == ast.ConstraintUniqKey || node.Tp == ast.ConstraintUniqIndex {
			idx.Uniq = true
		}
		idx.Cols = make([]string, 0, len(node.Keys))
		for _, key := range node.Keys {
			idx.Cols = append(idx.Cols, key.Column.Name.O)
		}
		v.Idxs = append(v.Idxs, idx)
	default:
		v.allType[fmt.Sprintf("%+v", reflect.TypeOf(node))] = struct{}{}
	}
	return in, false
}

func (v *TableDef) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func ParseTable(sql string) ([]TableDef, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}
	tabs := make([]TableDef, 0, len(stmtNodes))
	for _, stmt := range stmtNodes {
		if _, ok := stmt.(*ast.CreateTableStmt); !ok {
			continue
		}
		tab := TableDef{
			Stmt:    strings.TrimSpace(stmt.Text()),
			allType: make(map[string]struct{}),
		}
		if !strings.HasSuffix(tab.Stmt, ";") {
			tab.Stmt += ";"
		}
		stmt.Accept(&tab)
		tabs = append(tabs, tab)
	}
	return tabs, nil
}
