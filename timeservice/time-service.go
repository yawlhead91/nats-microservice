package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats"

	"github.com/yawlhead91/nats-microservices/transport"
)

// We use globals because it's a small application demonstrating NATS.
var nc *nats.Conn

func replyWithTime(m *nats.Msg) {
	curTime := Transport.Time{time.Now().Format(time.RFC3339)}
	data, err := proto.Marshal(&curTime)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Replying to ", m.Reply)
	nc.Publish(m.Reply, data)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments. Need NATS server address.")
		return
	}
	var err error
	nc, err = nats.Connect(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Connected to ", os.Args[1])

	nc.QueueSubscribe("TimeTeller", "TimeTellers", replyWithTime)
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
