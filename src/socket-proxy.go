package main

import (
	"time"
	"os"
	"fmt"
	"sync"
	"net"

	"github.com/farsightsec/golang-framestream"
)

var wgOut sync.WaitGroup
var wg sync.WaitGroup

func run() {
	fname := *sourceSocketDec
	os.Remove(fname)

	socket, err := net.Listen("unix", fname)
	if err != nil {
		fmt.Println(err, "Could not create socket!")
		fmt.Println(VERSION)
		os.Exit(1)
	}
	defer socket.Close()
	_ = os.Chmod(fname, 0777)

	go checkSocket()

	broker := NewBroker()
	go broker.Start()

	for _, val := range destSocket {
		wgOut.Add(1)
		go outConn(val, broker)
	}

	for {
		conn, err := socket.Accept()
		if err != nil {
			time.Sleep(RETRY_DELAY * time.Second)
			continue
		}
		wg.Add(1)
		go handleConn(conn, broker)
	}
	wg.Wait()
	broker.Stop()
	wgOut.Wait()
}

func handleConn(conn net.Conn, broker *Broker) {
	defer wg.Done()
	defer conn.Close()

	var FSContentType = []byte("protobuf:dnstap.Dnstap")
	bi := true
	timeout := time.Second

	readerOptions := framestream.ReaderOptions{
		ContentTypes:	[][]byte{FSContentType},
		Bidirectional:	bi,
		Timeout:	timeout,
	}
	fs, err := framestream.NewReader(conn, &readerOptions)
	if err != nil {
		return
	}
	buf := make([]byte, BUFFER_SIZE * KILOBYTE)
	for {
		n, err := fs.ReadFrame(buf)
		if err == framestream.EOF {
			break
		}
		if err != nil {
			continue
		}
		frame := make([]byte, n)
		copy(frame, buf[:n])
		broker.Publish(frame)
	}

}

func outConn(socket string, broker *Broker) {
	defer wgOut.Done()
	for {
		conn, err := net.Dial("unix", socket)
		if err != nil {
			time.Sleep(RETRY_DELAY * time.Second)
			continue
		}
		writeFrames(conn, broker)
	}
}

func writeFrames(conn net.Conn, broker *Broker) {
	var FSContentType = []byte("protobuf:dnstap.Dnstap")
	bi := true
	timeout := time.Second
	writerOptions := framestream.WriterOptions{
		ContentTypes:	[][]byte{FSContentType},
		Bidirectional:	bi,
		Timeout:	timeout,
	}
	fwriter, err := framestream.NewWriter(conn, &writerOptions)
	if err != nil {
		return
	}
	defer fwriter.Close()
	var count int
	for frame := range broker.Subscribe() {
		count++
		_, err = fwriter.WriteFrame(frame)
		if err != nil {
			return
		}
		err := fwriter.Flush()
		if err != nil {
			return
		}
	}
}

func handlepanic() {
	if a := recover(); a != nil {
		run()
	}
}

func checkSocket() {
	defer handlepanic()
	for {
		time.Sleep(5 * time.Second)
		socket := *sourceSocketDec
		_, err := os.Stat(socket)
		if err != nil {
			panic("Socket Missing!")
		}
	}
}

type Broker struct {
	stopCh	chan struct{}
	publishCh chan []byte
	subCh	 chan chan []byte
	unsubCh   chan chan []byte
}

func NewBroker() *Broker {
	return &Broker{
		stopCh:	make(chan struct{}),
		publishCh: make(chan []byte, 5000),
		subCh:	 make(chan chan []byte, 5000),
		unsubCh:   make(chan chan []byte, 5000),
	}
}

func (b *Broker) Start() {
	subs := map[chan []byte]struct{}{}
	for {
		select {
		case <-b.stopCh:
			for msgCh := range subs {
				close(msgCh)
			}
			return
		case msgCh := <-b.subCh:
			subs[msgCh] = struct{}{}
		case msgCh := <-b.unsubCh:
			delete(subs, msgCh)
		case msg := <-b.publishCh:
			for msgCh := range subs {
				// msgCh is buffered, use non-blocking send to protect the broker:
				select {
				case msgCh <- msg:
				default:
				}
			}
		}
	}
}

func (b *Broker) Stop() {
	close(b.stopCh)
}

func (b *Broker) Subscribe() chan []byte {
	msgCh := make(chan []byte, 5000)
	b.subCh <- msgCh
	return msgCh
}

func (b *Broker) Unsubscribe(msgCh chan []byte) {
	b.unsubCh <- msgCh
}

func (b *Broker) Publish(msg []byte) {
	b.publishCh <- msg
}
