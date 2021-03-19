package communication

import (
	"HappyOPQ/internal/app/common"
	"HappyOPQ/internal/app/opqbot/protocol/events"
	"HappyOPQ/pkg/log"
	sio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

func Connect(host string, port int, eventCh chan<- common.Convertible, flagCh chan<- int) *sio.Client {
	client, err := sio.Dial(
		sio.GetUrl(host, port, false),
		transport.GetDefaultWebsocketTransport(),
	)
	if err != nil {
		log.Error("无法与 OPQBot 连接，因为：", err)
		flagCh <- ConnectionTerminate
	}
	_ = client.On(sio.OnConnection, func(c *sio.Channel) {
		log.Info("成功与 OPQBot 连接，ID 为", client.Id())
		flagCh <- ConnectionSucceed
	})
	_ = client.On(sio.OnDisconnection, func(c *sio.Channel) {
		log.Error("与 OPQBot 断开连接")
		flagCh <- ConnectionTerminate
	})
	_ = client.On(sio.OnError, func(c *sio.Channel) {
		log.Error("与 OPQBot 的连接发生错误")
		flagCh <- ConnectionTerminate
	})
	_ = client.On("OnFriendMsgs", func(conn *sio.Channel, msg events.FriendMessage) {
		log.InfoF("收到 OnFriendMsgs 事件：%+v", msg)
		eventCh <- &msg
	})
	_ = client.On("OnGroupMsgs", func(c *sio.Channel, msg interface{}) {
		log.InfoF("收到 OnGroupMsgs 事件：%+v", msg)
	})
	_ = client.On("OnEvents", func(conn *sio.Channel, msg interface{}) {
		log.InfoF("收到 OnEvents 事件：%+v", msg)
	})
	return client
}
