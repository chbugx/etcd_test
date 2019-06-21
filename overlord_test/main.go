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

	//put

	for {
		putResp, err := kv.Put(context.TODO(), KeyPrefix+ModelUrl, "chb"+strconv.Itoa(int(time.Now().Unix())))
		if err != nil {
			fmt.Println("put model err: ", err)
		} else {
			fmt.Println("put model resp: ", putResp)
		}
		putResp, err = kv.Put(context.TODO(), KeyPrefix+ExpUrl, "chb"+strconv.Itoa(int(time.Now().Unix())))
		if err != nil {
			fmt.Println("put model err: ", err)
		} else {
			fmt.Println("put exp resp: ", putResp)
		}
		time.Sleep(time.Duration(20) * time.Second)
	}

	//lease
	//lease := etcd.NewLease(cli)

	//grantResp, err := lease.Grant(context.TODO(), 10)
	//putResp, err = kv.Put(context.TODO(), "/test/lease", "1105", etcd.WithLease(grantResp.ID))
	//if err != nil {
	//	fmt.Println("put err: ", err)
	//} else {
	//	fmt.Println("putResp: ", putResp)
	//}

	//fmt.Println(grantResp.TTL)
}

func main() {
	test()
}
