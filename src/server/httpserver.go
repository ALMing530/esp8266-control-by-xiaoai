package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func Serve() {
	go createSocket()
	http.HandleFunc("/ledOn", handleLed)
	err := http.ListenAndServeTLS(":53000", "1_alming.cn_bundle.crt", "2_alming.cn.key", nil)
	// err := http.ListenAndServe(":53000", nil)
	if err != nil {
		fmt.Println(err)
	}
}
func handleLed(w http.ResponseWriter, r *http.Request) {

	//if !isXiaoMiServer(r) {
	//	return
	//}
	ch = make(chan string, 1)
	if connectStatus {
		rev := parseReceive(r)
		msg := handleReceive(rev)

		res, _ := json.Marshal(msg)
		_, _ = w.Write(res)
	} else {
		message := defaultMessage()
		message.openMic(false)
		message.setTalkContent("客户端未连接")
		res, _ := json.Marshal(message)
		_, _ = w.Write(res)
	}
}
func isXiaoMiServer(r *http.Request) bool {
	method := r.Method
	url := r.RequestURI
	params := ""
	x_xiomi_date := r.Header.Get("X-Xiaomi-Date")
	host := "alming.cn:53000"
	contentType := r.Header.Get("Content-Type")
	ContentMd5 := r.Header.Get("Content-MD5")
	sigStr := method + "\n" +
		url + "\n" +
		params + "\n" +
		x_xiomi_date + "\n" +
		host + "\n" +
		contentType + "\n" +
		ContentMd5 + "\n"
	signature := computeSignature(sigStr)
	authorization := r.Header.Get("Authorization")
	recSignature := strings.Split(authorization, ":")
	return signature == recSignature[2]
}
func parseReceive(r *http.Request) *receive {
	body := r.Body
	buf := make([]byte, 1024)
	var reqContent []byte
	for {
		readLen, err := body.Read(buf)
		if readLen == 0 && err == io.EOF {
			break
		}
		reqContent = append(reqContent, buf[0:readLen]...)
	}
	rev := new(receive)
	err := json.Unmarshal(reqContent, rev)
	if err != nil {
		log.Println("parse json error")
		return nil
	}
	log.Println(rev.Query)
	return rev
}
func handleReceive(rev *receive) message {
	param := rev.Query
	switch param {
	case "打开LED":
		sendMessage("ledon")
		return ledFeedback()
	case "关闭LED":
		sendMessage("ledoff")
		return ledFeedback()
	case "退下":
		msg := defaultMessage()
		msg.openMic(false)
		msg.setTalkContent("那小爱先行告退喽")
		return msg
	case "进入LED控制系统":
		msg := defaultMessage()
		msg.setTalkContent("好的")
		return msg
	case "打开LED控制系统":
		msg := defaultMessage()
		msg.setTalkContent("好的")
		return msg
	default:
		msg := defaultMessage()
		msg.setTalkContent("你说啥")
		return msg

	}
}
func ledFeedback() message {
	msg := defaultMessage()
	feedback := "未知错误"
	//在另一个goroutines启动两秒定时器，(小爱请求超时时间为2.5s)两秒内未收到开发版回复由该goroutines触发超时动作
	timer := time.NewTimer(2 * time.Second)
	go func() {
		<-timer.C
		ch <- "time out"
	}()

	cmd := <-ch
	switch cmd {
	case "ledon":
		feedback = "LED已开启"
		timer.Stop()
	case "ledoff":
		feedback = "LED已关闭"
		timer.Stop()
	case "time out":
		feedback = "连接超时"
	case "err":
		feedback = "设备异常断开"
		timer.Stop()
	}
	
	close(ch)
	msg.setTalkContent(feedback)
	return msg
}
func computeSignature(content string) string {
	secret := "your secret"
	temp, _ := base64.StdEncoding.DecodeString(secret)
	h := hmac.New(sha256.New, temp)
	h.Write([]byte(content))
	signature := hex.EncodeToString(h.Sum(nil))
	return signature
}
