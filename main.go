package main

import (
	"context"
	"fmt"
	etcd "go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

const (
	KeyPrefix = "test"
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

	//put

	putResp, err := kv.Put(context.TODO(), KeyPrefix+"/c", "chb")
	putResp, err = kv.Put(context.TODO(), KeyPrefix+"/b", "chb")
	putResp, err = kv.Put(context.TODO(), KeyPrefix+"/a", "chb")
	if err != nil {
		fmt.Println("put err: ", err)
	} else {
		fmt.Println("putResp: ", putResp)
	}

	//get
	getResp, err := kv.Get(context.TODO(), KeyPrefix+"/test/a")
	if err != nil {
		fmt.Println("put err: ", err)
	} else {
		fmt.Println("getResp: ", getResp)
	}

	//get
	getResp, err = kv.Get(context.TODO(), KeyPrefix+"/test", etcd.WithPrefix())
	if err != nil {
		fmt.Println("put err: ", err)
	} else {
		fmt.Println("getResp: ", getResp)
	}

	//lease
	lease := etcd.NewLease(cli)

	grantResp, err := lease.Grant(context.TODO(), 10)
	putResp, err = kv.Put(context.TODO(), "/test/lease", "1105", etcd.WithLease(grantResp.ID))
	if err != nil {
		fmt.Println("put err: ", err)
	} else {
		fmt.Println("putResp: ", putResp)
	}

	fmt.Println(grantResp.TTL)
}

func main() {
	test()
}
