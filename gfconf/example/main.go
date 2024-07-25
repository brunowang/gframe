package main

import (
	"fmt"
	"github.com/brunowang/gframe/gfconf"
	"github.com/brunowang/gframe/gfconf/example/internal"
)

func main() {
	main1()
	main2()
}

// use nacos remote
func main1() {
	nacosConf := gfconf.NacosConf{
		Enabled:  true,
		Host:     "39.105.40.37",
		Port:     8848,
		Username: "nacos",
		Password: "ApQKpsb2HdYFq6@Q",
		Paths: []string{
			"/develop/example/app.yaml",
			"/develop/example/whitelist.yaml",
		},
	}
	gfconf.MustInitNacos(nacosConf, internal.HotConfRegister{})
	hotConf := internal.GetHotConf()
	fmt.Printf("nacos example app start at: %+v\n\tlog conf: %+v\n\twhitelist: %+v\n",
		hotConf.Application.App, hotConf.Application.Log, hotConf.Whitelist)
}

// use local file
func main2() {
	nacosConf := gfconf.NacosConf{
		Enabled: false,
		Paths: []string{
			"./gfconf/example/testyaml/app.yaml",
			"./gfconf/example/testyaml/whitelist.yaml",
		},
	}
	gfconf.MustInitNacos(nacosConf, internal.HotConfRegister{})
	hotConf := internal.GetHotConf()
	fmt.Printf("local example app start at: %+v\n\tlog conf: %+v\n\twhitelist: %+v\n",
		hotConf.Application.App, hotConf.Application.Log, hotConf.Whitelist)
}
