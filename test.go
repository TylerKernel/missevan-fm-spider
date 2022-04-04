package main

import (
	"fmt"
	"log"
	"maoer-fm-spider/util"
	"net"
)

func main() {
	client := util.NewClient()
	addr, err := net.ResolveTCPAddr("tcp4", "192.168.0.106:0")
	if err != nil {
		log.Fatal(err)
	}
	client.AddChannel(util.NewChannelWithLocalAddr(addr))
	resp, err := client.Get("https://fm.missevan.com/api/v2/chatroom/open/list?p=2&type=0")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Time       :", resp.Time())
	fmt.Println("  Cookie     :", resp.Cookies())
}
