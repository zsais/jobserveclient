package jsclient

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

var (
	defaultPort     = ":3000"
	timeoutDuration = 10 * time.Second
)

func processQueue(tasks chan *task) {
	for {
		currTask := <-tasks
		executeCommand(currTask)
	}
}

// ProcessConnections processes new connections made to the listener
func ProcessConnections(l net.Listener) {
	tasks := make(chan *task)
	go processQueue(tasks)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("cannot accept conn, err: %v\n", err)
			continue
		}

		go func() {
			defer conn.Close() // nolint: errcheck
			reader := bufio.NewReader(conn)
			for {
				conn.SetDeadline(time.Now().Add(timeoutDuration)) // nolint: errcheck
				input, err := reader.ReadString('\n')
				if err != nil {
					break
				}

				var request *taskRequest
				err = json.Unmarshal([]byte(input), &request)
				if err != nil {
					break
				}

				currTask := &task{conn, request, &taskResult{0, 0, 0, "", ""}}

				select {
				case tasks <- currTask: // add new task to channel if empty
				default: // respond with concurrencyError
					currTask.result.exitCode = errorCode
					currTask.result.error = concurrencyError
					sendJSON(currTask, errorOutput(currTask))
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
