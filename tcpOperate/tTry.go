package tcpOperate

import (
	"eInfusion/comm"
	"eInfusion/logs"
	ep "eInfusion/protocol"
	"net"
	//	"sync"
)

func TryTcpServer() {
	//	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":8989")
	//	checkError(err)
	netListen, err := net.Listen("tcp", ":"+c_TcpServer_Port)
	defer netListen.Close()
	//	系统开始运行时log记录时间
	logs.LogMain.Info(c_Msg_ServerStart + "（" + comm.GetCurrentDate() + "）")
	if err != nil {
		logs.LogMain.Critical("监听TCP出错", err)
		panic(err)
	}
	comm.Msg(comm.SprtLin(60))
	comm.Msg("TCP Port:" + c_TcpServer_Port)
	//	最大连接数不能超过规定数
	if len(g_Conns) <= c_MaxConnectionAmount {
		for {
			conn, err := netListen.Accept()
			if err != nil {
				continue
			}
			////////////////临时Tconn连接对象////////////////////////////////
			var c TConn
			c.ID = conn.RemoteAddr().String()
			c.IsAlive = true
			c.Conn = conn
			c.IPAddr = conn.RemoteAddr().(*net.TCPAddr).IP.String()
			g_Conns = append(g_Conns, c)
			///////////////////////////////////////////////////////////////
			comm.Msg(comm.SprtLin(60))
			logs.LogMain.Info("客户端：" + conn.RemoteAddr().String() + " 连接!")
			go tryreceiveData(c.Conn)
			//	time.Sleep(time.Second * 2)
			///////////////////////////////////////////////////////////////
		}
	} else {
		//超出连接数则不再接收连接
		logs.LogMain.Warn(c_MaxConnectionAmount)
	}
}

func tryreceiveData(conn net.Conn) {
	for {
		//	指定接收数据包头的帧长
		recDataHeader := make([]byte, ep.GetDataHeaderLength())
		_, err := conn.Read(recDataHeader)
		if err != nil {
			comm.Msg(comm.SprtLin(60))
			comm.Msg(conn.RemoteAddr(), " 客户端连接丢失!")
			return
		}
		// 数据包数据内容长度记录变量
		var intPckContentLength int
		// 判断包头是否正确，如果正确，获取长度
		if !ep.DecodeHeader(recDataHeader, &intPckContentLength) {
			comm.Msg("调试信息：数据包头不正确")
			continue
		}
		// 如果包头接收
		recDataContent := make([]byte, intPckContentLength)
		_, err = conn.Read(recDataContent)
		if !comm.CkErr("接收报文出错", err) {
			// 处理报文数据内容
			ep.DecodeRcvData(recDataContent, conn.RemoteAddr().(*net.TCPAddr).IP.String())
		}
	}
}

func trysendData(conn net.Conn, packetData []byte) {
	_, err := conn.Write(packetData) // don't care about return value
	defer conn.Close()
	if err != nil {
		comm.Msg(c_Msg_SendDataErr, err)
		logs.LogMain.Critical(c_Msg_SendDataErr, err)
		return
	}
}
