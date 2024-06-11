package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
)

// Version is the application version, it can be set at compile time
var Version = "dev"

func main() {
	help, displayVersion, target := parseArgs(os.Args)
	if help {
		displayHelp()
		os.Exit(0)
	}

	if displayVersion {
		fmt.Println("Version:", Version)
		os.Exit(0)
	}

	err := checkTarget(target)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func parseArgs(args []string) (bool, bool, string) {
	help := false
	version := false
	var target string

	for _, arg := range args[1:] {
		switch arg {
		case "-h", "--help":
			help = true
		case "-v", "--version":
			version = true
		default:
			if target != "" {
				fmt.Println("Multiple targets specified")
				os.Exit(1)
			}
			target = arg
		}
	}

	return help, version, target
}

func checkTarget(target string) error {
	u, err := url.Parse(target)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	switch u.Scheme {
	case "http", "https":
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
		resp, err := client.Get(target)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("get %s %d", target, resp.StatusCode)
		}
	case "tcp":
		conn, err := net.Dial("tcp", u.Host)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		conn.Close()
	default:
		return fmt.Errorf("invalid check type")
	}

	return nil
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
