package ntcp

import (
	cm "eInfusion/comm"
	"eInfusion/tlogs"
	tf "eInfusion/trsfus"
	"time"
)

// Srv :全局服务器对象
var Srv *TServer

// StartTCPService :启动TCP服务
func StartTCPService() {
	// 生成tcp包前缀
	prefixs := tf.MakePacketHeaderPrefix(0x66)
	h := NewTCPHeader(3, prefixs, 1)
	Srv = NewTCPServer(":9909", 10*time.Minute, 6*time.Hour, h)
	// 发送待发队列内的命令
	loopSendWaitQueue := func(c *Client) {
		for i, od := range Srv.waitQueue {
			if od.ID == cm.GetPureIPAddr(c.conn) {
				c.SendData(od.Data)
				// 册除待发队列相应项
				Srv.waitQueue = append(Srv.waitQueue[:i], Srv.waitQueue[i+1])
			}
		}
	}
	Srv.WhenNewClientConnected(func(c *Client) {
		if c.VerifyLegal() {
			tlogs.DoLog(tlogs.Info, "IP:", cm.GetPureIPAddr(c.conn), " Connected")
			loopSendWaitQueue(c) // 发送待发指令
			// 根据前台需求发送指令
			t := tf.MakeSendOrder(tf.CmdAddDetector, "A0000000", "B0000000", []string{})
			c.SendData(t)
		} else {
			// 非法客户端 TODO:

		}
	})

	Srv.WhenNewDataReceived(func(c *Client, p []byte) {
		// TODO: 解析，落到到具体业务

	})

	Srv.WhenClientConnectionClosed(func(c *Client, err error) {
		tlogs.DoLog(tlogs.Info, "IP:", cm.GetPureIPAddr(c.conn), " Disconnected")
	})
	Srv.Listen()
}
