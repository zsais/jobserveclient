package jsclient

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"os/exec"
	"syscall"
	"time"
)

var (
	concurrencyError = "cannot allow concurrent executions"
	timeoutError     = "timeout exceeded"
	errorCode        = -1
)

type task struct {
	conn    net.Conn
	request *taskRequest
	result  *taskResult
}

type taskRequest struct {
	Command []string `json:"command"`
	Timeout int      `json:"timeout"`
}

type taskResult struct {
	executedAt int64
	durationMs int64
	exitCode   int
	output     string
	error      string
}

func errorOutput(task *task) map[string]interface{} {
	out := map[string]interface{}{"command": task.request.Command,
		"timeout":   task.request.Timeout,
		"exit_code": task.result.exitCode,
		"error":     task.result.error}
	return out
}

func successOutput(task *task) map[string]interface{} {
	out := map[string]interface{}{"command": task.request.Command,
		"timeout":     task.request.Timeout,
		"executed_at": task.result.executedAt,
		"duration_ms": task.result.durationMs,
		"output":      task.result.output,
		"exit_code":   task.result.exitCode}
	return out
}

func sendJSON(task *task, output map[string]interface{}) {
	enc := json.NewEncoder(task.conn)
	err := enc.Encode(output)
	if err != nil {
		log.Printf("response could not be sent: %v\n", err)
	}
}

func executeCommand(task *task) {
	var response map[string]interface{}
	var ctx context.Context
	var cancel context.CancelFunc
	if task.request.Timeout != 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(task.request.Timeout)*time.Millisecond)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	cmd := exec.CommandContext(ctx, task.request.Command[0], task.request.Command[1:len(task.request.Command)]...)
	execStart := time.Now()
	out, err := cmd.Output()
	execElapsed := time.Since(execStart)

	if ctx.Err() == context.DeadlineExceeded {
		task.result.error = timeoutError
		task.result.exitCode = errorCode
		response = errorOutput(task)
	} else {
		if err != nil {
			task.result.error = err.Error()
			task.result.exitCode = errorCode
			response = errorOutput(task)
		} else {
			task.result.executedAt = execStart.Unix()
			task.result.exitCode = cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
			task.result.durationMs = int64(execElapsed) / 1000000
			task.result.output = string(out[:])
			response = successOutput(task)
		}
	}
	sendJSON(task, response)
}
