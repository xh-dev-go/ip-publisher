package main

import (
	"flag"
	"github.com/judwhite/go-svc"
	"github.com/xh-dev-go/ip-publisher/IpDetect"
	"github.com/xh-dev-go/xhUtils/flagUtils/flagString"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc/eventlog"
	"log"
	"strings"
)

type program struct {
	ipDetect *IpDetect.IpDetect
}

func main() {
	var loggerName = "ip-publisher"
	err := eventlog.InstallAsEventCreate(loggerName, eventlog.Info|eventlog.Warning|eventlog.Error)

	var msgToDisplayBeforeStart = ""
	if err != nil {
		if strings.HasSuffix(err.Error(), "registry key already exists") {
			msgToDisplayBeforeStart = err.Error()
		} else {
			panic(err)
		}
	}

	wlog, err := eventlog.Open(loggerName)

	if err != nil {
		panic(err)
	}
	IpDetect.Logging = func(msg string) {
		wlog.Info(windows.EVENTLOG_INFORMATION_TYPE, msg)
	}
	IpDetect.LogError = func(err error) {
		wlog.Error(windows.EVENTLOG_ERROR_TYPE, err.Error())
	}

	IpDetect.KeyLog = func(msg string) {
		wlog.Info(800, msg)
	}

	if msgToDisplayBeforeStart != "" {
		IpDetect.KeyLog(msgToDisplayBeforeStart)
	}

	IpDetect.CMD_TopicFlag = flagString.New("topic", "the topic to post").BindCmd()
	IpDetect.CMD_DeviceFlag = flagString.New("device", "the device id or code").BindCmd()
	IpDetect.CMD_UnFlag = flagString.New("username", "the username").BindCmd()
	IpDetect.CMD_PwdFlag = flagString.New("password", "the password").BindCmd()
	IpDetect.CMD_ServersCmd = flagString.New("servers", "servers url").BindCmd()
	flag.Parse()

	prg := &program{}
	prg.ipDetect = &IpDetect.IpDetect{}

	// Call svc.Run to start your program/service.
	if err := svc.Run(prg); err != nil {
		log.Fatal(err)
	}
}

func (p *program) Init(env svc.Environment) error {
	log.Printf("is win service? %v\n", env.IsWindowsService())
	p.ipDetect.Init()
	return nil
}

func (p *program) Start() error {
	p.ipDetect.Start()
	return nil
}

func (p *program) Stop() error {
	IpDetect.KeyLog("Stop service")
	p.ipDetect.Stopping <- struct{}{}

	select {
	case <-p.ipDetect.Stopped:
		IpDetect.KeyLog("Receive stop chan")
	}

	IpDetect.KeyLog("[done] Stop service")
	return nil
}
