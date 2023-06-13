package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
)

var (
	host     string
	port     int
	bindPort int
)

func init() {
	flag.StringVar(&host, "host", "", "remote host")
	flag.IntVar(&port, "port", 0, "remote port")
	flag.IntVar(&bindPort, "l", 0, "bind port")
	flag.Parse()

}

func bind(port int) (net.Listener, error) {
	return net.Listen("tcp", fmt.Sprintf(":%d", bindPort))
}
func main() {
	if bindPort != 0 {
		listener, err := bind(bindPort)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
		handleTcp(conn)
		os.Exit(0)
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
	handleInput(conn)
}

func handleInput(conn net.Conn) {
	var buffer = make([]byte, 2048)
	var rate = new(RateWriter)
	writer := io.MultiWriter(conn, rate)
	for {
		n, err := os.Stdin.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
		if n == 0 {
			break
		}
		_, err = writer.Write(buffer[:n])
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
		rate.count = rate.count + int64(len(buffer))
		rate.endTime = time.Now()
		fmt.Fprintf(os.Stderr, "\r%s", strings.Repeat(" ", 80))

		fmt.Fprintf(os.Stderr, "\r%s", rate)
	}
}
func handleTcp(conn net.Conn) {
	var buf = make([]byte, 2048)
	var rate = new(RateWriter)
	writer := io.MultiWriter(os.Stdout, rate)
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
		if n == 0 {
			break
		}
		_, err = writer.Write(buf[:n])
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
		rate.count = rate.count + int64(len(buf))
		rate.endTime = time.Now()
		fmt.Fprintf(os.Stderr, "\r%s", rate)
	}
}

type RateWriter struct {
	startTime time.Time
	endTime   time.Time
	count     int64
}

func (rate RateWriter) String() string {
	speed := math.Round(float64(rate.count) / rate.endTime.Sub(rate.startTime).Seconds())
	return humanize.Bytes(uint64(rate.count)) + " " + "[ " + humanize.Bytes(uint64(speed)) + "/s ]"

}
func (rate *RateWriter) Write(buf []byte) (int, error) {
	if rate.startTime.IsZero() {
		rate.startTime = time.Now()
		return len(buf), nil
	}
	rate.count = rate.count + int64(len(buf))
	rate.endTime = time.Now()
	return len(buf), nil
}
