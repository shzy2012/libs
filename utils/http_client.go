package utils

import (
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

//newClient 初始化http客户端
func newClient(timeout, maxIdelConns, maxConnsPerHost int) *http.Client {
	client := &http.Client{
		Timeout: time.Minute * time.Duration(timeout), //设置超时时间,默认0不设置超时时间
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second, //限制建立TCP连接的时间
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second, //限制 TLS握手的时间
			ResponseHeaderTimeout: 10 * time.Second, //限制读取response header的时间
			ExpectContinueTimeout: 1 * time.Second,  //限制client在发送包含 Expect: 100-continue的header到收到继续发送body的response之间的时间等待。
			MaxIdleConns:          maxIdelConns,     //连接池对所有host的最大连接数量，默认无穷大
			MaxConnsPerHost:       maxConnsPerHost,  //连接池对每个host的最大连接数量。
			IdleConnTimeout:       30 * time.Minute, //how long an idle connection is kept in the connection pool.
		},
	}

	return client
}

//NewRequest 发起post请求
func NewRequest(url string, data []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		log.Println("[newRequest]=> new request failed.")
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json") //设置Content-Type
	client := newClient(1, 100, 100)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[newRequest]=> dial tcp failed.")
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[newRequest]=> read response body faild.")
		return nil, err
	}

	return bytes, nil
}
