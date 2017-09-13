package jsclient

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"
)

var (
	concurrencyError = "cannot allow concurrent executions"
	timeoutError     = "timeout exceeded"
	exitCode         = -1
	nilTask          = &task{nil, &taskRequest{"", 0}, &taskResult{0, 0, 0, "", ""}}
)

type task struct {
	conn    net.Conn
	request *taskRequest
	result  *taskResult
}

type taskRequest struct {
	Command string `json:"command"`
	Timeout int    `json:"timeout"`
}

type taskResult struct {
	executedAt int64
	durationMs int64
	exitCode   int
	output     string
	error      string
}

func sendJSON(task *task, output map[string]interface{}) {
	enc := json.NewEncoder(task.conn)
	enc.Encode(output)
}

func errorOutput(task *task) {
	out := map[string]interface{}{"command": task.request.Command,
		"timeout":   task.request.Timeout,
		"exit_code": task.result.exitCode,
		"error":     task.result.error}

	sendJSON(task, out)
}

func successOutput(task *task) {
	out := map[string]interface{}{"command": task.request.Command,
		"timeout":     task.request.Timeout,
		"executed_at": task.result.executedAt,
		"duration_ms": task.result.durationMs,
		"output":      task.result.output,
		"exit_code":   task.result.exitCode}

	sendJSON(task, out)
}

func executeCommand(task *task) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(task.request.Timeout)*time.Second)
	defer cancel()
	arr := strings.Fields(task.request.Command)
	cmd := exec.CommandContext(ctx, arr[0], arr[1:len(arr)]...)

	execStart := time.Now()
	out, err := cmd.Output()
	execElapsed := time.Since(execStart)

	if ctx.Err() == context.DeadlineExceeded {
		task.result.error = timeoutError
		task.result.exitCode = exitCode
		errorOutput(task)
		return
	}

	if err != nil {
		log.Println(err)
	} else {
		task.result.executedAt = execStart.Unix()
		task.result.durationMs = int64(execElapsed) / 1000000
		task.result.output = string(out[:])
		successOutput(task)
	}
}
