package main

import (
	"context"
	"fmt"
	etcd "go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

const (
	KeyPrefix = "test/overlord"
	ExpUrl    = "/exp"
	ModelUrl  = "/model"
)

func test() {
	//client init
	cli, err := etcd.New(etcd.Config{
		Endpoints: []string{"192.168.91.7:2379"},
		//Endpoints:   []string{"localhost:2379"},
		DialTimeout: 10 * time.Second,
	})
	if err != nil {
		log.Fatalln("etcd new err: ", err)
	}
	kv := etcd.NewKV(cli)

	fmt.Println("cli begin")

	//get
	getResp, err := kv.Get(context.TODO(), KeyPrefix+ExpUrl, etcd.WithPrefix())
	if err != nil {
		fmt.Println("get err: ", err)
	} else {
		fmt.Println("getResp: ", getResp)
	}

	wt := etcd.NewWatcher(cli)
	wtModel := wt.Watch(context.Background(), KeyPrefix+ModelUrl)
	go func() {
		for {
			_ = <-wtModel
			getResp, err := kv.Get(context.TODO(), KeyPrefix+ModelUrl)
			if err != nil {
				fmt.Println("update model src err: ", err)
			} else {
				fmt.Println("update model src resp: ", getResp)
				versionModel := string(getResp.Kvs[0].Value)
				fmt.Println("key: ", KeyPrefix+versionModel+ExpUrl)
				getResp, err := kv.Get(context.TODO(), KeyPrefix+versionModel+ModelUrl)
				if err != nil {
					fmt.Println("update model err: err")
				} else {
					fmt.Println("update model resp: ", getResp)
					//todo update model conf
				}
			}
		}
	}()
	wtExp := wt.Watch(context.Background(), KeyPrefix+ExpUrl)
	go func() {
		for {
			_ = <-wtExp
			getResp, err := kv.Get(context.TODO(), KeyPrefix+ExpUrl)
			if err != nil {
				fmt.Println("update exp src err: ", err)
			} else {
				fmt.Println("update exp src resp: ", getResp)
				versionExp := string(getResp.Kvs[0].Value)
				fmt.Println("key: ", KeyPrefix+versionExp+ExpUrl)
				getResp, err := kv.Get(context.TODO(), KeyPrefix+versionExp+ExpUrl)
				if err != nil {
					fmt.Println("update exp err: err")
				} else {
					fmt.Println("update exp resp: ", getResp)
					//todo update exp conf
				}
			}
		}
	}()
	time.Sleep(time.Duration(10) * time.Hour)
}

func main() {
	test()
}
