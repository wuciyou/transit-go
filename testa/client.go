package main

import (
	"net"
)

func main()  {

	start()
}

func start(){
	conn,err := net.DialUDP("udp4",nil,&net.UDPAddr{
		IP:net.IPv4(255,255,255,255),
		Port:2549,
	})

	if err != nil{
		panic(err)
	}

	conn.Write([]byte("hello wuciyou"))
}
