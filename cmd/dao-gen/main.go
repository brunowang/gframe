package main

import (
	"github.com/brunowang/gframe/cmd/dao-gen/gen/sqlparser"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
)

var (
	file = kingpin.Flag("file", "sql file path").Short('f').Required().String()
)

func main() {
	kingpin.Parse()

	bs, err := os.ReadFile(*file)
	if err != nil {
		log.Fatalf("read file failed, err: %v", err)
	}

	tabs, err := sqlparser.ParseTable(string(bs))
	if err != nil {
		log.Fatalf("parse table failed, err: %v", err)
	}
	for _, v := range tabs {
		f, err := os.Create(v.TabName + ".go")
		if err != nil {
			log.Fatalf("create file failed, err: %v", err)
		}
		_, err = f.WriteString(v.Render() + "\n")
		if err != nil {
			log.Fatalf("write file failed, err: %v", err)
		}
	}
}
