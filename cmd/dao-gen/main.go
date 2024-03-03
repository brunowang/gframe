package main

import (
	"github.com/brunowang/gframe/cmd/dao-gen/gen"
	"github.com/brunowang/gframe/cmd/dao-gen/gen/sqlparser"
	_ "github.com/pingcap/tidb/pkg/parser/test_driver"
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

	if err := gen.GenerateDAO(tabs); err != nil {
		log.Fatalf("generate dao failed, err: %v", err)
	}

	for _, tab := range tabs {
		if err := gen.GenerateModel(tab); err != nil {
			log.Fatalf("generate model failed, err: %v", err)
		}
		if err := gen.GenerateCURD(tab); err != nil {
			log.Fatalf("generate curd failed, err: %v", err)
		}
		// use gfcache package
		//if err := gen.GenerateCache(tab); err != nil {
		//	log.Fatalf("generate cache failed, err: %v", err)
		//}
	}
}
