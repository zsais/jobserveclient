package jsclient

import (
	"bufio"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorOuput(t *testing.T) {
	newTask := &task{nil, &taskRequest{[]string{"ls", "-al"}, 1000}, &taskResult{0, 0, -1, "", timeoutError}}
	out := errorOutput(newTask)

	assert.Equal(t, out["command"], newTask.request.Command)
	assert.Equal(t, out["timeout"], newTask.request.Timeout)
	assert.Equal(t, out["exit_code"], newTask.result.exitCode)
	assert.Equal(t, out["error"], newTask.result.error)
}

func TestSuccessOuput(t *testing.T) {
	newTask := &task{nil, &taskRequest{[]string{"echo", "GOPATH"}, 500}, &taskResult{15, 1, 1, "GOPATH", ""}}
	out := successOutput(newTask)

	assert.Equal(t, out["command"], newTask.request.Command)
	assert.Equal(t, out["timeout"], newTask.request.Timeout)
	assert.Equal(t, out["exit_code"], newTask.result.exitCode)
	assert.Equal(t, out["executed_at"], newTask.result.executedAt)
	assert.Equal(t, out["duration_ms"], newTask.result.durationMs)
}

func TestSendJSON(t *testing.T) {
	l, err := NewListener()
	if err != nil {
		t.Fatalf("Error opening test server: %v", err)
	}
	defer l.Close() // nolint: errcheck

	client, err := net.Dial("tcp", defaultPort)
	if err != nil {
		t.Fatalf("Error opening test client: %v", err)
	}
	defer client.Close() // nolint: errcheck

	conn, err := l.Accept()
	if err != nil {
		t.Fatalf("Error accepting connection: %v", err)
	}

	newTask := &task{conn, &taskRequest{[]string{"ls", "-al"}, 1000}, &taskResult{0, 0, -1, "", timeoutError}}
	out := errorOutput(newTask)
	sendJSON(newTask, out)

	reader := bufio.NewReader(client)
	response, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Error reading output: %v", err)
	}

	assert.Equal(t, response, "{\"command\":[\"ls\",\"-al\"],\"error\":\"timeout exceeded\",\"exit_code\":-1,\"timeout\":1000}\n")

	conn.Close() // nolint: errcheck
	sendJSON(newTask, out)
}

func TestExecuteCommand(t *testing.T) {

}
