package site

import (
	"net"
	"net/http"
	"time"
)

var MaxTimeout = time.Millisecond * 1500
var AuthHeader = ""

func GetTransport() http.Transport {
	return http.Transport{
		MaxIdleConns: 100,
		Dial: func(netw, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(netw, addr, MaxTimeout) //设置建立连接超时
			if err != nil {
				return nil, err
			}
			err = conn.SetDeadline(time.Now().Add(MaxTimeout)) //设置发送接受数据超时
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
}

func init() {

}
