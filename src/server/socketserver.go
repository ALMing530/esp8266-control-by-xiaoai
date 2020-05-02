package server

import (
	"fmt"
	"net"
	"sync"
)

var con net.Conn
var ch chan string
var connectStatus bool = false
var lock sync.Mutex

func createSocket() {

	addr, err := net.ResolveTCPAddr("", ":53055")
	if err != nil {
		fmt.Println("resolve ip address fail")
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println("listen fail")
	}
	if listener != nil {
		for {
			
			con, err = listener.Accept()
			go handleAccept(con)
			if err != nil {
				fmt.Println("accept fail ")
			}
		}

	}
}
func handleAccept(connection net.Conn) {
	fmt.Println("一个设备以连接")
	defer func() {
		err := recover()
		if err != nil {
			connectStatus = false
		}
	}()
	connectStatus = true
	buf := make([]byte, 128)
	for {
		redLen, err := connection.Read(buf)
		if err != nil {
			fmt.Println("read error")
			_ = connection.Close()
			connectStatus = false
			break
		}
		str := string(buf[0 : redLen])
		if str == "ledon" {
			ch <- "ledon"
		} else if str == "ledoff" {
			ch <- "ledoff"
		}
	}

}
func sendMessage(message string) {
	//这个锁暂时是没用的，可去掉
	lock.Lock()
	if con != nil {
		_, err := con.Write([]byte(message))
		fmt.Println(err)
		if err != nil {
			ch <- "err"
			_ = con.Close()
			connectStatus = false
			fmt.Println("write data fail")
		}
	}
	lock.Unlock()
}
