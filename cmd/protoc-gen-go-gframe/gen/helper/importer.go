package helper

import (
	"sort"
	"strings"
)

type Importer struct {
	data  []importPathAndAlias
	cache map[string]struct{}
}

type importPathAndAlias struct {
	path  string
	alias string
}

func NewImporter() *Importer {
	return &Importer{cache: make(map[string]struct{})}
}

func (i *Importer) Import(importPath string) *Importer {
	if _, has := i.cache[importPath]; has {
		return i
	}
	i.data = append(i.data,
		importPathAndAlias{path: importPath},
	)
	i.cache[importPath] = struct{}{}
	return i
}

func (i *Importer) ImportWithAlias(importPath string, alias string) *Importer {
	if _, has := i.cache[importPath]; has {
		return i
	}
	i.data = append(i.data,
		importPathAndAlias{
			path:  importPath,
			alias: alias,
		},
	)
	i.cache[importPath] = struct{}{}
	return i
}

func (i *Importer) String() string {
	sort.Slice(i.data, func(a, b int) bool {
		return i.data[a].path < i.data[b].path
	})
	var buf strings.Builder
	buf.WriteString("import (\n")
	for _, im := range i.data {
		if im.alias != "" {
			buf.WriteRune('\t')
			buf.WriteString(im.alias)
			buf.WriteRune(' ')
		}
		buf.WriteString("\t\"")
		buf.WriteString(im.path)
		buf.WriteString("\"\n")
	}
	buf.WriteString(")\n")
	return buf.String()
}
