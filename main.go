package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"path/filepath"
	"time"
)

var (
	dataFile   string
	outputFile string
	host       string
	port       int
	timeout    int
	bufferSize int
	testTimes  int
)

func init() {

	short := " (shorthand)"

	fileUsage := "the file to send to tcp server"
	flag.StringVar(&dataFile, "file", "", fileUsage)
	flag.StringVar(&dataFile, "f", "", fileUsage+short)

	hostUsage := "server ip"
	flag.StringVar(&host, "ip", "127.0.0.1", hostUsage)
	flag.StringVar(&host, "i", "127.0.0.1", hostUsage+short)

	hostPortUsage := "server port"
	flag.IntVar(&port, "port", 6788, hostPortUsage)
	flag.IntVar(&port, "p", 6788, hostPortUsage+short)

	timeoutUsage := "timeout to wait server response"
	flag.IntVar(&timeout, "timeout", 100, timeoutUsage)
	flag.IntVar(&timeout, "t", 100, timeoutUsage+short)

	bufferSizeUsage := "the buffer size"
	flag.IntVar(&bufferSize, "buffer", 1024, bufferSizeUsage)
	flag.IntVar(&bufferSize, "b", 1024, bufferSizeUsage+short)

	testTimesUsage := "test times"
	flag.IntVar(&testTimes, "times", 1, testTimesUsage)
	flag.IntVar(&testTimes, "n", 1, testTimesUsage+short)

	outputFileUsage := "the content response form server will write this file"
	flag.StringVar(&outputFile, "output", "", outputFileUsage)
	flag.StringVar(&outputFile, "o", "", outputFileUsage+short)
}

func main() {
	flag.Parse()
	if dataFile == "" {
		err := fmt.Errorf("Error: file not set.")
		log.Fatalln(err)
		flag.Usage()
		return
	}

	if bufferSize <= 0 {
		bufferSize = 1024
	}

	strTcpAddrress := fmt.Sprintf("%s:%d", host, port)

	var tcpAddr *net.TCPAddr
	if addr, e := net.ResolveTCPAddr("tcp4", strTcpAddrress); e != nil {
		log.Fatalln(e)
		return
	} else {
		tcpAddr = addr
	}

	var fileData []byte
	if data, e := ioutil.ReadFile(dataFile); e != nil {
		log.Fatalln(e)
		return
	} else {
		fileData = data
	}

	durTimeout := time.Duration(timeout) * time.Millisecond

	for i := 0; i < testTimes; i++ {
		var tcpConn *net.TCPConn
		if conn, e := net.DialTCP("tcp", nil, tcpAddr); e != nil {
			log.Fatalln(e)
			return
		} else {
			tcpConn = conn
			tcpConn.SetKeepAlive(false)
			defer tcpConn.Close()
		}

		tcpConn.SetReadDeadline(time.Now().Add(durTimeout))
		tcpConn.Write(fileData)

		var buf []byte
		buf = make([]byte, bufferSize)

		index := 0

		for {
			if size, e := tcpConn.Read(buf[index:]); e != nil {
				break
			} else {
				index += size
				if index > bufferSize {
					log.Fatalln("buffer size is too small")
					return
				}
			}
		}

		writeData := buf[0:index]

		fileName, ext := filename_and_ext(outputFile)
		if testTimes != 1 {
			fileName = fmt.Sprintf("%s_%d%s", fileName, i, ext)
		} else {
			fileName = outputFile
		}

		if outputFile != "" {
			if e := ioutil.WriteFile(fileName, writeData, 0666); e != nil {
				log.Fatalln("write response to file failed")
				return
			}
		} else {
			fmt.Printf("[%05d]: %s\n", i+1, string(writeData))
		}
	}
}

func filename_and_ext(fileName string) (name string, ext string) {
	ext = filepath.Ext(fileName)
	name = string([]byte(fileName)[0 : len(fileName)-len(ext)])
	return
}
