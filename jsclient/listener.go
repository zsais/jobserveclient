package jsclient

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

var (
	defaultPort     = ":3030"
	timeoutDuration = 10 * time.Second
)

func processQueue(tasks chan *task) {
	for {
		currTask := <-tasks
		executeCommand(currTask)

		_ = <-tasks
	}
}

// ProcessConnections processes new connections made to the listener
func ProcessConnections(l net.Listener) {
	tasks := make(chan *task)
	go processQueue(tasks)
	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Printf("cant accept conn, err: %v\n", err)
			continue
		}
		go func() {
			defer conn.Close()
			reader := bufio.NewReader(conn)
			for {
				conn.SetDeadline(time.Now().Add(timeoutDuration))
				input, err := reader.ReadString('\n')
				println(input)
				if err != nil {
					break
				}
				var request *taskRequest

				err = json.Unmarshal([]byte(input), &request)
				if err != nil {
					break
				}
				currTask := &task{conn, request, &taskResult{0, 0, 0, "", ""}}

				if len(tasks) == 0 {
					tasks <- currTask
					tasks <- nilTask
				} else {
					currTask.result.exitCode = exitCode
					currTask.result.error = concurrencyError
					errorOutput(currTask)
				}
			}
		}()
	}
}

// NewListener creates a new tcp server on port 3030
func NewListener() (net.Listener, error) {
	l, err := net.Listen("tcp", defaultPort)
	return l, err
}
