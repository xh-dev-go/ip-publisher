package main

import (
	"fmt"
	"github.com/xh-dev-go/ip-publisher/IpDetect"
	"github.com/xh-dev-go/ip-publisher/info"
)

func initLog() {
	IpDetect.KeyLog("===== IP Publisher started =====")
	IpDetect.KeyLog(fmt.Sprintf("Build at: %s", info.BuildTime))
	IpDetect.KeyLog(fmt.Sprintf("Build cimmit: %s", info.CommitHash))
	IpDetect.Logging("")
	IpDetect.Logging("")
	IpDetect.KeyLog(fmt.Sprintf("Device name: %s", IpDetect.CMD_DeviceFlag.Value()))
	IpDetect.KeyLog(fmt.Sprintf("User name: %s", IpDetect.CMD_UnFlag.Value()))
	IpDetect.KeyLog(fmt.Sprintf("Server hosts: %s", IpDetect.CMD_ServersCmd.Value()))
	IpDetect.KeyLog(fmt.Sprintf("Topic: %s", IpDetect.CMD_TopicFlag.Value()))

	var abFunc = func(testResult bool, trueCase string, falseCase string) string {
		if testResult {
			return trueCase
		} else {
			return falseCase
		}
	}

	IpDetect.KeyLog(fmt.Sprintf("Detection: %d %s", IpDetect.CMD_DetectionPeriod.Value(),
		abFunc(IpDetect.CMD_DetectionPeriod.Value() > 1, "minutes", "minute")))
	IpDetect.KeyLog(fmt.Sprintf("Detection cache count: %d %s", IpDetect.CMD_DetectionCacheCount.Value(),
		abFunc(IpDetect.CMD_DetectionCacheCount.Value() > 1, "times", "time")))

	IpDetect.Logging("")
	IpDetect.Logging("")
	IpDetect.Logging("")
}
