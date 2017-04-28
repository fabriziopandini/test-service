package main // import "github.com/fabriziopandini/test-service"

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"io"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", doHostname)
	router.HandleFunc("/echo", doEcho)
	router.HandleFunc("/echoheaders", doEchoheaders)
	router.HandleFunc("/hostname", doHostname)
	router.HandleFunc("/fqdn", doFqdn)
	router.HandleFunc("/ip", doIp)
	router.HandleFunc("/env", doEnv)
	router.HandleFunc("/exit/{exitCode:[0-9]+}", doExit)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func doEcho(w http.ResponseWriter, r *http.Request) {
	io.Copy(w, r.Body)
}

func doEchoheaders(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		fmt.Fprintf(w, "%s=%s\n", k, v)
	}
}

func doHostname(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(w, "Error: %s!\n", err)
		return
	}
	fmt.Fprintf(w, "%s\n", hostname)
}

func doEnv(w http.ResponseWriter, r *http.Request) {
	for _, env := range os.Environ() {
		fmt.Fprintf(w, "%s\n", env)
	}
}

func doIp(w http.ResponseWriter, r *http.Request) {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Fprintf(w, "Error: %s!\n", err)
		return
	}

	// handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Fprintf(w, "Error: %s!\n", err)
			return
		}
		// handle err
		for _, addr := range addrs {
			var ipAddr net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ipAddr = v.IP
			case *net.IPAddr:
				ipAddr = v.IP
			}
			fmt.Fprintf(w, "%s\n", ipAddr)
		}
	}
}

func doFqdn(w http.ResponseWriter, r *http.Request) {
	// from https://github.com/ShowMax/go-fqdn/blob/master/fqdn.go
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(w, "Error: %s!\n", err)
		return
	}

	addrs, err := net.LookupIP(hostname)
	if err != nil {
		fmt.Fprintf(w, "%s\n", hostname)
		return
	}

	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			ip, err := ipv4.MarshalText()
			if err != nil {
				fmt.Fprintf(w, "%s\n", hostname)
				return
			}
			hosts, err := net.LookupAddr(string(ip))
			if err != nil || len(hosts) == 0 {
				fmt.Fprintf(w, "%s\n", hostname)
				return
			}
			fqdn := hosts[0]
			fmt.Fprintf(w, "%s\n", strings.TrimSuffix(fqdn, ".")) // return fqdn without trailing dot
			return
		}
	}
	fmt.Fprintf(w, "%s\n", hostname)
}

func doExit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	exitCode, _ := strconv.Atoi(vars["exitCode"])

	os.Exit(exitCode)
}
