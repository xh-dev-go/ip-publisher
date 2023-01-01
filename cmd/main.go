package main

import (
	"flag"
	"github.com/xh-dev-go/ip-publisher/IpDetect"
	"github.com/xh-dev-go/xhUtils/flagUtils/flagString"
	"log"
)

func main() {
	IpDetect.CMD_TopicFlag = flagString.New("topic", "the topic to post").BindCmd()
	IpDetect.CMD_DeviceFlag = flagString.New("device", "the device id or code").BindCmd()
	IpDetect.CMD_UnFlag = flagString.New("username", "the username").BindCmd()
	IpDetect.CMD_PwdFlag = flagString.New("password", "the password").BindCmd()
	IpDetect.CMD_ServersCmd = flagString.New("servers", "servers url").BindCmd()

	IpDetect.Logging = func(msg string) {
		log.Println(msg)
	}
	IpDetect.LogError = func(err error) {
		log.Fatal(err)
	}
	IpDetect.KeyLog = func(msg string) {
		log.Println(msg)
	}

	flag.Parse()

	publisher := IpDetect.IpDetect{}
	publisher.Init()
	publisher.Start()

	select {
	case <-publisher.Stopped:
	}
	IpDetect.Logging("Done")
}
