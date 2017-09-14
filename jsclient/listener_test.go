package jsclient

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	l, err := NewListener()
	if err != nil {
		t.Fatalf("Error opening test server: %v", err)
	}
	defer l.Close() // nolint: errcheck
	assert.Equal(t, reflect.TypeOf(l).String(), "*net.TCPListener")
	assert.Equal(t, err, nil)

	_, err = NewListener()
	assert.Equal(t, reflect.TypeOf(err).String(), "*net.OpError")
}

func TestProcessQueue(t *testing.T) {

}

func TestProcessConnections(t *testing.T) {
	// l, err := NewListener()
	// if err != nil {
	// 	t.Fatalf("Error opening test server: %v", err)
	// }
	// go ProcessConnections(l)
	//
	// client, err := net.Dial("tcp", defaultPort)
	// if err != nil {
	// 	t.Fatalf("Error opening test client: %v", err)
	// }
	// defer client.Close()
	//
	// _, err = client.Write([]byte("{'command':['sleep', '4'], 'timeout':1}\n"))
	// if err != nil {
	// 	t.Fatalf("Error writing to client: %v", err)
	// }
	//
	// _, err = client.Write([]byte("{'command':['sleep', '1'], 'timeout':1000}\n"))
	// if err != nil {
	// 	t.Fatalf("Error writing to client: %v", err)
	// }
}
