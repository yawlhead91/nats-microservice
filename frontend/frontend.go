package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats"

	"github.com/yawlhead91/nats-microservices/transport"
)

var nc *nats.Conn

func handleUserWithTime(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	myUser := Transport.User{Id: vars["id"]}
	curTime := Transport.Time{}
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		data, err := proto.Marshal(&myUser)
		if err != nil || len(myUser.Id) == 0 {
			fmt.Println(err)
			w.WriteHeader(500)
			fmt.Println("Problem with parsing the user Id.")
			return
		}

		msg, err := nc.Request("UserNameById", data, 100*time.Millisecond)
		if err == nil && msg != nil {
			myUserWithName := Transport.User{}
			err := proto.Unmarshal(msg.Data, &myUserWithName)
			if err == nil {
				myUser = myUserWithName
			}
		}
		wg.Done()
	}()

	go func() {
		msg, err := nc.Request("TimeTeller", nil, 100*time.Millisecond)
		if err == nil && msg != nil {
			receivedTime := Transport.Time{}
			err := proto.Unmarshal(msg.Data, &receivedTime)
			if err == nil {
				curTime = receivedTime
			}
		}
		wg.Done()
	}()

	wg.Wait()
	fmt.Fprintln(w, "Hello ", myUser.Name, " with id ", myUser.Id, ", the time is ", curTime.Time, ".")
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

	m := mux.NewRouter()
	m.HandleFunc("/{id}", handleUserWithTime)
	http.ListenAndServe(":3000", m)
}
