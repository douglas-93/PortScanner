package main

import (
	"fmt"
	"net"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	result := make(chan string)

	fmt.Println("#### Scan Started!  ####")

	// ports := []int{80, 443, 8080}
	for port := 1; port <= 65535; port++ {
		wg.Add(1)
		go scanPort("192.168.1.1", port, &wg, result)
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	var openPorts []string
	var closedPorts []string

	for msg := range result {
		if strings.Contains(msg, "open") {
			openPorts = append(openPorts, msg)
		} else {
			closedPorts = append(closedPorts, msg)
		}
	}
	fmt.Println("#### Scan Finished! ####")
	fmt.Println("Open Ports:")
	slices.Sort(openPorts)
	for _, port := range openPorts {
		fmt.Println(port)
	}
	fmt.Println("Closed Ports:")
	slices.Sort(closedPorts)
	for _, port := range closedPorts {
		fmt.Println(port)
	}
}

func scanPort(host string, port int, wg *sync.WaitGroup, result chan string) {
	defer wg.Done()
	connectionType := "tcp"
	address := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout(connectionType, address, 3*time.Second)
	if err != nil {
		result <- fmt.Sprintf("Port %5d is closed or filtered", port)
		return
	}
	defer conn.Close()
	result <- fmt.Sprintf("Port %5d is open", port)
}
