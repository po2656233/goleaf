package main

import (
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"github.com/xtaci/kcp-go"
	"io"
	"net"
)

func main() {
	fmt.Println("kcp listens on 10000")
	lis, err := kcp.ListenWithOptions(":10000", nil, 13, 3)
	if err != nil {
		panic(err)
	}
	for {
		conn, e := lis.AcceptKCP()
		if e != nil {
			panic(e)
		}
		go func(conn net.Conn) {
			var buffer = make([]byte, 1024, 1024)
			for {
				n, e := conn.Read(buffer)
				if e != nil {
					if e == io.EOF {
						fmt.Println("EOF")
						break
					}
					fmt.Println(errorx.Wrap(e))
					break
				}
				conn.Write([]byte("ok"))
				//fmt.Println("receive from client:", buffer[:n])
				fmt.Printf("receive from client:%s\n", buffer[:n])
			}
		}(conn)
	}
}
