package IpDetect

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/xh-dev-go/ip-publisher/IpDetect/cache"
	"github.com/xh-dev-go/xhUtils/flagUtils/flagInt"
	"github.com/xh-dev-go/xhUtils/flagUtils/flagString"
	"github.com/xh-dev-go/xhUtils/xhKafka/KHeader"
	"io"
	"net/http"
	"strings"
	"time"
)

var Logging func(msg string)
var KeyLog func(msg string)
var LogError func(err error)

var CMD_DetectionPeriod *flagInt.IntParam
var CMD_DetectionCacheCount *flagInt.IntParam

var CMD_TopicFlag *flagString.StringParam
var CMD_DeviceFlag *flagString.StringParam
var CMD_UnFlag *flagString.StringParam
var CMD_PwdFlag *flagString.StringParam
var CMD_ServersCmd *flagString.StringParam

type GetIpResponse struct {
	Time  time.Time
	Value string
	Err   error
}

func (resp GetIpResponse) HasError() bool {
	return resp.Err != nil
}

func WithError(now time.Time, err error) GetIpResponse {
	return GetIpResponse{
		now,
		"",
		err,
	}
}
func NoError(now time.Time, msg string) GetIpResponse {
	return GetIpResponse{
		now,
		msg,
		nil,
	}
}

func GetIp(respChannel chan GetIpResponse) {
	resp, err := http.Get("https://api.myip.com")
	now := time.Now()
	defer resp.Body.Close()

	if err != nil {
		respChannel <- WithError(now, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		respChannel <- WithError(now, err)
	}
	respChannel <- NoError(now, string(body))
}

var w *kafka.Writer = nil

func initWriter(servers []string, un, pwd string) {

	mechanism, err := scram.Mechanism(scram.SHA256, un, pwd)
	if err != nil {
		LogError(err)
	}

	dialer := &kafka.Dialer{
		SASLMechanism: mechanism,
		TLS:           &tls.Config{},
	}

	w = kafka.NewWriter(kafka.WriterConfig{
		Brokers: servers,
		Dialer:  dialer,
	})
}

func SendToKafka(topic, key string, msg string, completeChan chan bool) {
	Logging("Send message to kafka")
	var headers KHeader.KafkaHeaders

	err := w.WriteMessages(context.Background(), kafka.Message{
		Topic:   topic,
		Key:     []byte(key),
		Value:   []byte(msg),
		Headers: headers.ToKafkaHeaders(),
	})

	if err != nil {
		LogError(err)
		return
	}

	Logging("[done] Send message to kafka")
	completeChan <- true

}

func (ipDetect *IpDetect) Start() {
	var servers []string
	for _, val := range strings.Split(CMD_ServersCmd.Value(), ",") {
		server := strings.TrimSpace(val)
		if server == "" {
			panic(errors.New("server is not allow empty"))
		}
		servers = append(servers, server)
	}

	if len(servers) == 0 {
		err := errors.New("no server is passed in")
		LogError(err)
		panic(err)
	}

	initWriter(servers, CMD_UnFlag.Value(), CMD_PwdFlag.Value())

	chanOfGetIp := make(chan GetIpResponse)
	doneSendMessage := make(chan bool)
	ticker := time.NewTicker(time.Duration(CMD_DetectionPeriod.Value()) * time.Minute)

	go GetIp(chanOfGetIp)

	var quiting = false
	go func() {
		defer w.Close()

		//var cachedFor = 0
		//var lastRecord = ""
		//var detectionCacheCount = CMD_DetectionCacheCount.Value()

		for {
			select {
			case resp := <-chanOfGetIp:
				KeyLog("start getting ip")

				if resp.HasError() {
					println(resp.Err)
				} else {
					ipDetect.cacheEngine.CacheInternal(resp.Value, func(msg string) {
						go SendToKafka(CMD_TopicFlag.Value(), CMD_DeviceFlag.Value(), msg, doneSendMessage)
					})

					//val := fmt.Sprint("%n\n", md5.Sum([]byte(resp.Value)))
					//Logging(fmt.Sprintf("Cached for %d", cachedFor))
					//if lastRecord == "" {
					//	lastRecord = val
					//	cachedFor += 1
					//	go SendToKafka(CMD_TopicFlag.Value(), CMD_DeviceFlag.Value(), val, doneSendMessage)
					//} else if lastRecord == val {
					//	if cachedFor < detectionCacheCount {
					//		cachedFor += 1
					//		Logging("ip address cached")
					//	} else {
					//		go SendToKafka(CMD_TopicFlag.Value(), CMD_DeviceFlag.Value(), val, doneSendMessage)
					//		cachedFor = 0
					//	}
					//} else {
					//	cachedFor = 1
					//	go SendToKafka(CMD_TopicFlag.Value(), CMD_DeviceFlag.Value(), val, doneSendMessage)
					//}
				}
				Logging("[done] get ip")
			case <-doneSendMessage:
				KeyLog("[done] received complete send message")
			case <-ticker.C:
				go GetIp(chanOfGetIp)
				Logging("[done] received ticket")
			case <-ipDetect.Stopping:
				KeyLog("Received message of quiting")
				quiting = true
			}
			if quiting {
				break
			}
		}

		Logging("[Done] Received message of quiting")
		ipDetect.Stopped <- struct{}{}
	}()
}

func (ipDetect *IpDetect) Init() {
	ipDetect.Stopped = make(chan struct{})
	ipDetect.Stopping = make(chan struct{})
	if KeyLog == nil {
		panic(errors.New("KeyLog function not init"))
	}
	if Logging == nil {
		panic(errors.New("Logging function not init"))
	}
	if LogError == nil {
		panic(errors.New("LogError function not init"))
	}
	ipDetect.cacheEngine = cache.NewDefault(CMD_DetectionCacheCount.Value(), Logging)
}

type IpDetect struct {
	Stopping    chan struct{}
	Stopped     chan struct{}
	cacheEngine cache.CacheEngine
}
