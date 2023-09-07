package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"server/model"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

func main() {
	configFile, err := os.Open("config.txt")
	if err != nil {
		fmt.Printf("Error opening config file: %v\n", err)
		os.Exit(1)
	}
	defer configFile.Close()

	// Create a new scanner and scan the first line
	scanner := bufio.NewScanner(configFile)
	scanner.Scan()
	GPRC_PORT := scanner.Text()
	GPRC_PORT = strings.TrimSpace(GPRC_PORT)

	// Convert the string to an integer
	port, err := strconv.Atoi(GPRC_PORT)
	if err != nil {
		fmt.Printf("Error converting port to integer: %v\n", err)
		os.Exit(1)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	model.RegisterMiddlemanServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
