package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
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

	if len(os.Args) < 2 {
		displayHelp()
		os.Exit(1)
	}

	target := os.Args[1]
	u, err := url.Parse(target)
	if err != nil {
		fmt.Println("Error: Invalid target URL")
		os.Exit(1)
	}

	switch u.Scheme {
	case "http":
		resp, err := http.Get(target)
		if err != nil || resp.StatusCode != http.StatusOK {
			fmt.Println("Error:", err)
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
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	case "tcp":
		conn, err := net.Dial("tcp", u.Host)
		if err != nil {
			fmt.Println("Error:", err)
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
   check [url]

   Example:
   check tcp://example.com:2222
   check http://example.com:8080
   check https://example.com:8443

Version:
   ` + Version + `

Description:
   check is a high performance check tool whose command can be started
   by using this command. If none of the *http*, *https*, *tcp*, or *udp* commands
   are specified, the default action of the **check** command is to display this help.

Global Options:
   --help         show help
   --version      print the version`)
}
