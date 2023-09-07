package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"time"
)

func main() {

	fmt.Printf("Starting middle server")

	// cmd := exec.Command("nohup", "./gateway", arg)
	// cmd := exec.Command("bash", "-c", "nohup ./server "+" > /dev/null 2>&1 & disown")
	cmd := exec.Command("./server")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	// Start separate goroutines to read stdout and stderr
	go readOutput(stdout)
	go readOutput(stderr)

	// Start the command and do not wait for it to finish

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)
	fmt.Printf("Middleman Server Started")

}

// readOutput reads data from the given reader and writes it to the console
func readOutput(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
