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
	var versionModel string
	var versionExp string
	//put
	getResp, err := kv.Get(context.TODO(), KeyPrefix+ModelUrl)
	if err != nil {
		fmt.Println("get model init err: ", err)
	} else {
		fmt.Println("get model init resp: ", getResp)
		versionModel = string(getResp.Kvs[0].Value)
		fmt.Println("get model version info: ", versionModel)
	}

	getResp, err = kv.Get(context.TODO(), KeyPrefix+ExpUrl)
	if err != nil {
		fmt.Println("get exp init err: ", err)
	} else {
		fmt.Println("get exp init resp: ", getResp)
		versionExp = string(getResp.Kvs[0].Value)
		fmt.Println("get exp version info: ", versionExp)
	}

	for {
		//lease
		lease := etcd.NewLease(cli)

		grantResp, err := lease.Grant(context.TODO(), 100)
		versionNew := "/" + time.Unix(time.Now().Unix(), 0).Format("20060102150405")
		putResp, err := kv.Put(context.TODO(), KeyPrefix+versionNew+ModelUrl, "chb"+versionNew, etcd.WithLease(grantResp.ID))
		if err != nil {
			fmt.Println("put model err: ", err)
		} else {
			fmt.Println("put model resp: ", putResp)
			putResp, err := kv.Put(context.TODO(), KeyPrefix+ModelUrl, versionNew)
			if err != nil {
				fmt.Println("put model src err: ", err)
			} else {
				fmt.Println("put model src resp: ", putResp)
			}
		}

		putResp, err = kv.Put(context.TODO(), KeyPrefix+versionNew+ExpUrl, "chb"+versionNew, etcd.WithLease(grantResp.ID))
		if err != nil {
			fmt.Println("put exp err: ", err)
		} else {
			fmt.Println("put exp resp: ", putResp)
			putResp, err := kv.Put(context.TODO(), KeyPrefix+ExpUrl, versionNew)
			if err != nil {
				fmt.Println("put exp src err: ", err)
			} else {
				fmt.Println("put exp src resp: ", putResp)
			}
		}
		time.Sleep(time.Duration(20) * time.Second)
	}
}

func main() {
	test()
}
