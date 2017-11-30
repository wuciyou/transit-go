package main

import (
	"net"
	"fmt"
	"sync"
)

func main()  {

	w := &sync.WaitGroup{}
	w.Add(1)
	start()
	w.Wait()
}

func start(){
	conn,err := net.ListenUDP("udp4",&net.UDPAddr{
		IP:net.IPv4(0,0,0,0),
		Port:2549,
	})

	if err != nil{
		panic(err)
	}

	go newConn(conn)
}

func newConn(conn *net.UDPConn){
	var data = make([]byte,1024)
	for{
		n,e := conn.Read(data)
		if e != nil{
			panic(e)
		}
		fmt.Printf("data:%s \n",string(data[:n]))
	}
}