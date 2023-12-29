package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

// Version is the application version, it can be set at compile time
var Version string

func main() {
	help := flag.Bool("help", false, "Display Help")
	displayVersion := flag.Bool("version", false, "Display Version")
	flag.Parse()

	if *help {
		displayHelp()
		os.Exit(0)
	}

	if *displayVersion {
		fmt.Println("Version:", Version)
		os.Exit(0)
	}

	if len(os.Args) < 3 {
		displayHelp()
		os.Exit(1)
	}

	checkType := os.Args[1]
	target := os.Args[2]

	switch strings.ToLower(checkType) {
	case "http":
		resp, err := http.Get(target)
		if err != nil || resp.StatusCode != http.StatusOK {
			fmt.Println("Error: Failed to connect to the target")
			os.Exit(1)
		}
	case "https":
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
		resp, err := client.Get(target)
		if err != nil || resp.StatusCode != http.StatusOK {
			fmt.Println("Error: Failed to connect to the target")
			os.Exit(1)
		}
	case "tcp":
		conn, err := net.Dial("tcp", target)
		if err != nil {
			fmt.Println("Error: Failed to connect to the target")
			os.Exit(1)
		}
		conn.Close()
	case "udp":
		conn, err := net.Dial("udp", target)
		if err != nil {
			fmt.Println("Error: Failed to connect to the target")
			os.Exit(1)
		}
		conn.Close()
	default:
		fmt.Println("Error: Invalid check type")
		displayHelp()
		os.Exit(1)
	}
}

func displayHelp() {
	fmt.Println(`Usage:
   check command [arguments...]

   Example:
   check tcp example.com:2222
   check udp example.com:5353
   check http http://example.com:8080
   check https https://example.com:8443

Version:
   ` + Version + `

Description:
   check is a high performance check tool whose command can be started
   by using this command. If none of the *http*, *https*, *tcp*, or *udp* commands
   are specified, the default action of the **check** command is to display this help.

Commands:
   http    check http target
   https   check https target
   tcp     check tcp target
   udp     check udp target
   help    Shows a list of commands or help for one command

Global Options:
   --help         show help
   --version      print the version`)
}
