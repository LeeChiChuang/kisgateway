package main

import (
	"fmt"
	"io/ioutil"
	"kisgateway/proxyhttp"
	"kisgateway/serverlib/conf"
	"kisgateway/serverlib/logx"
	"kisgateway/serverlib/store"
)

type (
	Config struct {
		Etcd struct {
			Hosts []string
			Key   string
		}
		DataSource string
		Table      string
		Cache      []NodeConf
	}

	NodeConf struct {
		RedisConf
		Weight int `json:",default=100"`
	}

	RedisConf struct {
		Host string
		Type string `json:",default=node,options=node|cluster"`
		Pass string `json:",optional"`
	}
)

func loadFile(path string) ([]byte, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}

	return content, nil
}

func main()  {
	/*
	c, err := conf.ReadConf("/Users/simon/code/go/src/kisgateway/cmd/php.ini")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v \n", c)
	 */
	conf.ReadConf("/Users/simon/code/go/src/kisgateway/conf/kisgateway.ini")
	logx.SetUp()
	logx.Info(fmt.Sprintf("%v", conf.ConfigInfo))

	store.InitStore("test")

	s := proxyhttp.New()
	s.Start()
}

