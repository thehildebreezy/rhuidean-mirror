package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/webview/webview"
)

// some constants
const tcpClientPort = "50998"
const tcpServerPort = "50999"
const host = "manetheren"

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

// hold our message queue
var messageQueue []string

func popMessage() string {
	if len(messageQueue) == 0 {
		return ""
	}
	var x string
	x, messageQueue = messageQueue[0], messageQueue[1:]
	return x
}

func requestUpdate() {
	sendRequestTime()
	sendRequestWeather()
	sendRequestForecast()
}

func sendRequestTime() {
	sendMessage(requestTime, "?v=1")
}

func sendRequestWeather() {
	sendMessage(requestWeather, "?v=1")
}

func sendRequestForecast() {
	sendMessage(requestForecast, "?v=0&other=simple")
}

func sendMessage(messageType byte, message string) {
	c, err := net.Dial("tcp" /*host+*/, "192.168.1.27:"+tcpServerPort)
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
	fmt.Printf("%d Bytes written to tcp\n", n)
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

func startServer(w webview.WebView) {

	l, err := net.Listen("tcp", ":"+tcpClientPort)
	if err != nil {
		fmt.Println(err)
		return
	}
	go awaitConnections(l, w)
}

func awaitConnections(l net.Listener, w webview.WebView) {
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

func acceptConnection(c net.Conn, w webview.WebView) {
	defer c.Close()

	data, err := ioutil.ReadAll(c)
	if err != nil {
		fmt.Println("Read from TCP connection failed")
		fmt.Println(err)
		return
	}
	fmt.Println("From TCP:")
	fmt.Println(string(data))

	//evalString := "updateDisplay('" + string(data) + "');"
	//w.Eval(evalString)
	messageQueue = append(messageQueue, string(data))
}

func main() {
	debug := true
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("Minimal webview example")
	w.SetSize(800, 600, webview.HintNone)

	// create some bindings
	fmt.Println("Binding service update")
	w.Bind("serviceUpdate", func() {
		fmt.Println("received serviceUpdate")
		requestUpdate()
	})

	w.Bind("pollMessage", func() string {
		fmt.Println("received poll")
		return popMessage()
	})

	startServer(w)

	w.Navigate("http://" + host + "/mirror/index.html")
	w.Run()

}
