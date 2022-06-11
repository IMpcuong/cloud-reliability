package main

import (
	"fmt"
	"os/exec"
	"testing"
)

func TestPortAvailable(t *testing.T) {
	host := "localhost"
	port := "3331"
	isEstablished := checkPort(host, port)
	if isEstablished {
		fmt.Printf("Port %s is established with host %s\n", port, host)
	} else {
		t.Errorf("Cannot establish a connection to port %s", port)
	}
}

// Nil pointer dereference when execute command!
func TestOpenPort(t *testing.T) {
	cmd := exec.Command("python", "-m", "http.server", "3331")
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestOpenConn(t *testing.T) {
	node := Node{"localhost:3331"}
	conn, err := openConn(node)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(conn)
	if conn != nil {
		defer conn.Close()
	}
}
