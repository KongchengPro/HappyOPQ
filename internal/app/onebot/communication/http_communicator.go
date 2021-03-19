package communication

import (
	"HappyOPQ/pkg/log"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type HTTPCommunicator struct {
	URL string
}

func (c HTTPCommunicator) StartAPIServer() error {
	panic("implement me")
}

func (c HTTPCommunicator) Report(qq int64, event interface{}) error {
	jsonData, err := json.Marshal(event)
	if err != nil {
		log.Error("序列化事件时发生错误：", err)
		return err
	}
	request, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(jsonData))
	if request == nil {
		log.Error("向用户端转发事件时（构造请求时）出现错误：", err)
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Self-ID", strconv.FormatInt(qq, 10))
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Error("向用户端转发事件时（发送请求时）出现错误：", err)
		return err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Error("向用户端转发事件时（关闭响应时）出现错误：", err)
		}
	}()
	// TODO 快速操作
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("向用户端转发事件时（读取响应内容时）出现错误：", err)
		return err
	}
	log.InfoF("已向用户端转发事件：%+v", event)
	log.DebugF("返回内容：%+v", respBody)
	return nil
}
