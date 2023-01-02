package cache

import (
	"crypto/md5"
	"fmt"
)

type CacheEngine struct {
	cacheFor int
	maxCount int
	hashVal  string
	Logging  func(msg string)
}

func NewDefault(maxCount int, Logging func(msg string)) CacheEngine {
	val := CacheEngine{}
	defVal := 0
	val.cacheFor = defVal
	val.maxCount = maxCount
	defStr := ""
	val.hashVal = defStr
	val.Logging = Logging
	return val
}

func (cacheEngine *CacheEngine) CacheInternal(msg string, handler func(msg string)) {

	val := fmt.Sprint("%n\n", md5.Sum([]byte(msg)))
	if cacheEngine.hashVal == "" {
		cacheEngine.hashVal = val
		cacheEngine.cacheFor += 1
		cacheEngine.Logging("init engine")
		handler(msg)
	} else if cacheEngine.hashVal == val {
		if cacheEngine.cacheFor < cacheEngine.maxCount {
			cacheEngine.cacheFor += 1
			cacheEngine.Logging("ip address cached")
		} else {
			//go SendToKafka(CMD_TopicFlag.Value(), CMD_DeviceFlag.Value(), val, doneSendMessage)
			cacheEngine.Logging("ip address send")
			handler(msg)
			cacheEngine.cacheFor = 0
		}
	} else {
		cacheEngine.cacheFor = 1
		cacheEngine.Logging("ip address changed")
		//go SendToKafka(CMD_TopicFlag.Value(), CMD_DeviceFlag.Value(), val, doneSendMessage)
		handler(msg)
	}
	cacheEngine.Logging(fmt.Sprintf("cached for %d", cacheEngine.cacheFor))

}
