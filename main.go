package main

import (
	"flag"
	"fmt"
	etcd "github.com/asim/go-micro/plugins/registry/etcd/v4"
	"go-micro.dev/v4/registry"
	"strconv"
	"time"
)

func main() {
	var addr string
	var user string
	var password string
	var serviceName string
	flag.StringVar(&addr, "addr", "127.0.0.1:2379", "etcd地址，格式类似127.0.0.1:2379")
	flag.StringVar(&user, "user", "", "etcd用户名")
	flag.StringVar(&password, "password", "", "etcd用户密码")
	flag.StringVar(&serviceName, "name", "", "要查询的服务名称")
	flag.Parse()

	var svrlist []*registry.Service
	var err error
	reg := GetRegistry(user, password, addr)
	if serviceName == "" {
		svrlist, err = reg.ListServices()
	} else {
		svrlist, err = reg.GetService(serviceName)
	}
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	mapAddress := make(map[string]string)
	for i, v := range svrlist {
		for _, node := range v.Nodes {
			mapAddress[v.Name+"\t"+strconv.Itoa(i)] = node.Address
		}
	}
	for k, v := range mapAddress {
		fmt.Println(k, "\t\t", v)
	}
}

func GetRegistry(Username string, Password string, Address ...string) registry.Registry {
	var addrs []string
	var name, password string
	if len(Address) == 0 {
		addrs = []string{"127.0.0.1:2379"}
		name = ""
		password = ""
	} else {
		addrs = Address
		name = Username
		password = Password
	}
	reg := etcd.NewRegistry(
		registry.Addrs(addrs...),
		registry.Timeout(5*time.Second),
		etcd.Auth(name, password),
	)
	return reg
}
