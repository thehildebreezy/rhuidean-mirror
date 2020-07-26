package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"sync"

	"github.com/webview/webview"
)

// some constants
const tcpClientPort = 50998
const tcpServerPort = 50999
const host = "stoneoftear"

// interval to update
const updateInterval = 10 * 60 * 1000

// byte definitions
const (
	startByte       byte = 0x00
	manetherenByte  byte = 0xFA
	serveWeather    byte = 0x00
	serveForecast   byte = 0x01
	serveQuote      byte = 0x02
	serveTime       byte = 0x03
	serveCalendar   byte = 0x04
	serveTasks      byte = 0x05
	serveConfig     byte = 0x06
	serveOther      byte = 0x07
	requestWeather  byte = 0x08
	requestForecast byte = 0x09
	requestQuote    byte = 0x0a
	requestTime     byte = 0x0b
	requestCalendar byte = 0x0c
	requestTasks    byte = 0x0d
	requestConfig   byte = 0x0e
	requestOther    byte = 0x0f
)

func requestUpdate() {
	requestTime()
	requestWeather()
	requestForecast()
}

func requestTime() {
	sendMessage(requestTime, "?v=1")
}

func requestWeather() {
	sendMessage(requestWeather, "?v=1")
}

func requestForecast() {
	sendMessage(requestForecast, "?v=0&other=simple")
}

func sendMessage(messageType byte, message string) {
	c, err := net.Dial("tcp", host+":"+tcpServerPort)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()
	n, err := c.Write(formatMessage(messageType, message))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func formatMessage(messageType byte, message string) []byte {
	buf := make([]byte, 7+len(message))
	sizeBuf := make([]byte, 4)

	// start the message up
	buf[0] = startByte
	buf[1] = manetherenByte

	// save the size
	binary.BigEndian.PutUint32(sizeBuf, uint32(len(message)))
	copy(buf[2:6], sizeBuf)

	// set request type
	buf[6] = messageType

	// copy the message
	copy(buf[7:], message)

	return buf
}

func startServer(wg WaitGroup, w WebView) {
	defer wg.Done()

	l, err := net.Listen("tcp", host+":"+tcpClientPort)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	// accept new connections
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go acceptConnection(c, w)
	}

}

func acceptConnection(c Conn, w WebView) {
	defer c.Close()

	data, err := ioutil.ReadAll(c)
	if err != nil {
		fmt.Println("Read from TCP connection failed")
		fmt.Println(err)
		return
	}
	fmt.Println("From TCP:")
	fmt.Println(string(data))

	evalString := "updateDisplay('" + string(data) + "');"
	w.Eval(evalString)
}

func main() {
	debug := true
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("Minimal webview example")
	w.SetSize(800, 600, webview.HintNone)

	w.Navigate("http://localhost/mirror/index.html")
	w.Run()

	// create some bindings
	w.Bind("serviceUpdate", func() {
		requestUpdate()
	})

	// now spawn a new go routine to
	wg := sync.WaitGroup

	// add new routine to wg
	wg.Add(1)
	go startServer(wg, w)
	wg.Wait()
}
