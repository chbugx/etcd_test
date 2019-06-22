package main

import (
	"context"
	"fmt"
	etcd "go.etcd.io/etcd/clientv3"
	"log"
	"strconv"
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
		versionModel = string(getResp.Kvs[0].Value)
	}

	getResp, err := kv.Get(context.TODO(), KeyPrefix+ExpUrl)
	if err != nil {
		fmt.Println("get exp init err: ", err)
	} else {
		versionExp = string(getResp.Kvs[0].Value)
	}

	for {
		//lease
		lease := etcd.NewLease(cli)

		grantResp, err := lease.Grant(context.TODO(), 100)

		putResp, err := kv.Put(context.TODO(), KeyPrefix+versionModel+ModelUrl, "chb"+strconv.Itoa(int(time.Now().Unix())), etcd.WithLease(grantResp.ID))
		if err != nil {
			fmt.Println("put model err: ", err)
		} else {
			fmt.Println("put model resp: ", putResp)
			putResp, err := kv.Put(context.TODO(), KeyPrefix+ModelUrl, strconv.Itoa(int(time.Now().Unix())))
			if err != nil {
				fmt.Println("put model src err: ", err)
			} else {
				fmt.Println("put model src resp: ", putResp)
			}
		}

		putResp, err = kv.Put(context.TODO(), KeyPrefix+versionExp+ExpUrl, "chb"+strconv.Itoa(int(time.Now().Unix())), etcd.WithLease(grantResp.ID))
		if err != nil {
			fmt.Println("put exp err: ", err)
		} else {
			fmt.Println("put exp resp: ", putResp)
			putResp, err := kv.Put(context.TODO(), KeyPrefix+ExpUrl, strconv.Itoa(int(time.Now().Unix())))
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
