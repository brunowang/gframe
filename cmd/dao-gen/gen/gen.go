package gen

import (
	"github.com/brunowang/gframe/cmd/dao-gen/gen/helper"
	"github.com/brunowang/gframe/cmd/dao-gen/gen/sqlparser"
	"github.com/brunowang/gframe/cmd/dao-gen/gen/template"
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
	f, err := helper.CreateFile("dao/tbl.go")
	if err != nil {
		return err
	}
	_, err = f.WriteString(tpl.RenderSQL())
	if err != nil {
		return err
	}

	f, err = helper.CreateFile("dao/dao.go")
	if err != nil {
		return err
	}
	_, err = f.WriteString(tpl.Render())
	if err != nil {
		return err
	}
	return nil
}

func GenerateModel(tab sqlparser.TableDef) error {
	tpl := &template.ModelDefTmpl{
		Name: helper.ToCamelCase(tab.Name),
	}
	imports := make(map[string]struct{})
	for _, col := range tab.Cols {
		if col.Type == "time.Time" {
			imports["time"] = struct{}{}
		}
		tpl.Fields = append(tpl.Fields, template.Field{
			Name:    helper.ToCamelCase(col.Name),
			Type:    col.Type,
			ZeroVal: col.ZeroVal,
			ColName: col.Name,
			Comment: col.Comment,
		})
	}
	tpl.Imports = imports
	f, err := helper.CreateFile("dao/" + tab.Name + ".go")
	if err != nil {
		return err
	}
	_, err = f.WriteString(tpl.Render())
	if err != nil {
		return err
	}
	return nil
}

func GenerateCURD(tab sqlparser.TableDef) error {
	tpl := &template.CurdDefTmpl{
		Name:    helper.ToCamelCase(tab.Name),
		TabName: tab.Name,
	}
	for _, col := range tab.Cols {
		tpl.Fields = append(tpl.Fields, template.Field{
			Name:    helper.ToCamelCase(col.Name),
			Type:    col.Type,
			ColName: col.Name,
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
	f, err := helper.CreateFile("dao/" + tpl.TabName + "_dao.go")
	if err != nil {
		return err
	}
	_, err = f.WriteString(tpl.Render())
	if err != nil {
		return err
	}
	return nil
}

func GenerateCache(tab sqlparser.TableDef) error {
	tpl := &template.CacheDefTmpl{}
	{
		f, err := helper.CreateFile("dao/cache_multi.go")
		if err != nil {
			return err
		}
		_, err = f.WriteString(tpl.Render())
		if err != nil {
			return err
		}
	}
	{
		f, err := helper.CreateFile("dao/cache_local.go")
		if err != nil {
			return err
		}
		_, err = f.WriteString(tpl.RenderLocal())
		if err != nil {
			return err
		}
	}
	{
		f, err := helper.CreateFile("dao/cache_redis.go")
		if err != nil {
			return err
		}
		_, err = f.WriteString(tpl.RenderRedis())
		if err != nil {
			return err
		}
	}
	{
		f, err := helper.CreateFile("dao/cache_serial.go")
		if err != nil {
			return err
		}
		_, err = f.WriteString(tpl.RenderSerial())
		if err != nil {
			return err
		}
	}
	return nil
}
