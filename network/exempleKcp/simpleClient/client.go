package main

import (
	"fmt"
	"github.com/xtaci/kcp-go"
)

func main() {
	kcpconn, err := kcp.DialWithOptions("localhost:10000", nil, 10, 3)
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			b := make([]byte, 1024)
			n, err := kcpconn.Read(b)
			if err != nil {
				fmt.Println("Read err:", err, n)
			}
			fmt.Printf("Read:%s n:%v \n", b[:n], n)
		}
	}()

	go func() {
		for {
			var a string
			n, err := fmt.Scanf("%s--\n", &a)
			kcpconn.Write([]byte(a))
			if err != nil {
				fmt.Println("Write err:", err, n)
				//break
			}
		}
	}()

	select {}
}
