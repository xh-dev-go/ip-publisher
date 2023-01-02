//go:build windows

package main

import (
	"flag"
	"fmt"
	"github.com/judwhite/go-svc"
	"github.com/xh-dev-go/ip-publisher/IpDetect"
	"github.com/xh-dev-go/xhUtils/flagUtils/flagBool"
	"github.com/xh-dev-go/xhUtils/flagUtils/flagInt"
	"github.com/xh-dev-go/xhUtils/flagUtils/flagString"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc/eventlog"
	"log"
	"os"
	"strings"
)

func main() {
	var root = os.Args[0]
	root = root[:strings.LastIndex(root, "\\")]
	log.Println(root)

	f, err := os.Create(fmt.Sprintf("%s\\error.log", root))
	f.WriteString(root)
	f.WriteString("\n")
	if err != nil {
		panic(err)
	}

	var loggerName = "ip-publisher"
	err = eventlog.InstallAsEventCreate(loggerName, eventlog.Info|eventlog.Warning|eventlog.Error)
	var wLogAvailable = false

	var msgToDisplayBeforeStart = ""
	if err != nil {
		if strings.HasSuffix(err.Error(), "registry key already exists") {
			msgToDisplayBeforeStart = err.Error()
			wLogAvailable = true
		} else if strings.HasSuffix(err.Error(), "Access is denied.") {
			f.WriteString(err.Error())
			f.WriteString("\n")
			wLogAvailable = false
		} else {
			f.WriteString(err.Error())
			f.WriteString("\n")
			panic(err)
		}
	}

	var wlog *eventlog.Log
	if wLogAvailable {
		wlog, err = eventlog.Open(loggerName)

		if err != nil {
			f.WriteString(err.Error())
			f.WriteString("\n")
			panic(err)
		}
	}

	var outFile = true
	IpDetect.Logging = func(msg string) {
		log.Println(msg)
		if wLogAvailable {
			wlog.Info(windows.EVENTLOG_INFORMATION_TYPE, msg)
		}
		if outFile {
			f.WriteString(fmt.Sprintf(msg))
			f.WriteString("\n")
		}
	}
	IpDetect.LogError = func(err error) {

		if wLogAvailable {
			wlog.Error(windows.EVENTLOG_ERROR_TYPE, err.Error())
		}
		if outFile {
			f.WriteString(fmt.Sprint(err))
			f.WriteString("\n")
		}
	}

	IpDetect.KeyLog = func(msg string) {
		log.Println(msg)
		if wLogAvailable {
			wlog.Info(800, msg)
		}
		if outFile {
			f.WriteString(fmt.Sprintf(msg))
			f.WriteString("\n")
		}
	}

	if msgToDisplayBeforeStart != "" {
		IpDetect.KeyLog(msgToDisplayBeforeStart)
	}

	IpDetect.CMD_DetectionPeriod = flagInt.NewDefault("detection-period", "the detection period in minutes", 1).BindCmd()
	IpDetect.CMD_DetectionCacheCount = flagInt.NewDefault("detection-cache-count", "the detection cache period count. \n The app prevent ip publish when detection cache period reach or the ip address changed", 1).BindCmd()

	IpDetect.CMD_TopicFlag = flagString.New("topic", "the topic to post").BindCmd()
	IpDetect.CMD_DeviceFlag = flagString.New("device", "the device id or code").BindCmd()
	IpDetect.CMD_UnFlag = flagString.New("username", "the username").BindCmd()
	IpDetect.CMD_PwdFlag = flagString.New("password", "the password").BindCmd()
	IpDetect.CMD_ServersCmd = flagString.New("servers", "servers url").BindCmd()

	outToFile := flagBool.NewDefault("out-file", "also output to default file", false).BindCmd()
	flag.Parse()

	if !outToFile.Value() {
		outFile = outToFile.Value()
		f.Close()
	}

	initLog()

	prg := &program{}
	prg.ipDetect = &IpDetect.IpDetect{}

	// Call svc.Run to start your program/service.
	if err := svc.Run(prg); err != nil {
		log.Fatal(err)
	}

}

type program struct {
	ipDetect *IpDetect.IpDetect
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

	IpDetect.KeyLog("===== IP Publisher ended =====")
	return nil
}
