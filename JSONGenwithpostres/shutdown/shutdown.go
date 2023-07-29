package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Process struct {
	Id   string
	Name string
}

// for windows
func main() {
	if len(os.Args) < 2 {
		log.Fatal("Port numbers are required.")
	}

	osType := runtime.GOOS

	if osType == "windows" {
		var wg sync.WaitGroup
		processesToKill := make(map[string]Process) // Mapping of port:process
		var hasError bool = false
		// Collect the processes
		for _, portStr := range os.Args[1:] {
			port, err := strconv.Atoi(portStr)
			if err != nil {
				log.Fatalf("Invalid port number: %s", portStr)
			}

			wg.Add(1)
			go func(p int) {
				defer wg.Done()

				// Get the processes
				output, err := exec.Command("powershell", "-Command", fmt.Sprintf("netstat -ano | findstr :%d", p)).CombinedOutput()

				if err != nil {
					hasError = true
					log.Printf("Error getting process on port %d: %v. Output: %s", p, err, string(output))
					return
				}

				if len(output) > 0 {
					lines := strings.Split(strings.TrimSpace(string(output)), "\n")
					for _, line := range lines {
						parts := strings.Fields(line)
						processId := parts[len(parts)-1] // The process id is the last item on the line

						// Get the process name
						nameOutput, err := exec.Command("powershell", "-Command", fmt.Sprintf("(Get-Process -Id %s).Name", processId)).CombinedOutput()
						if err != nil {
							log.Printf("Error getting process name for PID %s on port %d: %v. Output: %s", processId, p, err, string(nameOutput))
							return
						}
						processName := strings.TrimSpace(string(nameOutput))

						key := fmt.Sprintf("%d:%s", p, processId) // Key is port:process
						processesToKill[key] = Process{Id: processId, Name: processName}
					}
				} else {
					hasError = true
					log.Printf("No process found running on port %d", p)
				}
			}(port)
		}

		wg.Wait()

		// Now kill the processes
		for key, process := range processesToKill {
			err := exec.Command("powershell", "-Command", fmt.Sprintf("Stop-Process -Id %s -Force", process.Id)).Run()
			if err != nil {
				hasError = true
				log.Printf("Error terminating process PID %s (name: %s) on port %s: %v", process.Id, process.Name, strings.Split(key, ":")[0], err)
			} else {
				log.Printf("Terminated process PID %s (name: %s) on port %s", process.Id, process.Name, strings.Split(key, ":")[0])
			}
		}

		if hasError {
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	} else if osType == "darwin" {
		// For each port number arg, get the pid and kill
		for _, portStr := range os.Args[1:] {
			port, err := strconv.Atoi(portStr)
			if err != nil {
				log.Fatalf("Invalid port number: %s", portStr)
			}

			// Construct the command and arguments
			cmdLsof := exec.Command("lsof", "-ti", fmt.Sprintf(":%d", port))
			cmdKill := exec.Command("xargs", "kill")

			// Pipe the output of lsof to the input of kill
			cmdKill.Stdin, _ = cmdLsof.StdoutPipe()

			// Start the lsof command
			if err := cmdLsof.Start(); err != nil {
				fmt.Printf("Error starting lsof command: %v\n", err)
				return
			}

			// Start the kill command and wait for it to finish
			if err := cmdKill.Run(); err != nil {
				fmt.Printf("Error running kill command: %v\n", err)
				return
			}

			// Wait for the lsof command to finish
			if err := cmdLsof.Wait(); err != nil {
				fmt.Printf("Error waiting for lsof command: %v\n", err)
				return
			}

			fmt.Println("Processes on pid " + portStr + " killed successfully.")
			time.Sleep(1 * time.Second)

		}
	}
}

// macos version
// func main() {
// 	if len(os.Args) < 2 {
// 		log.Fatal("Port numbers are required.")
// 	}

// 	var wg sync.WaitGroup
// 	processesToKill := make(map[string]Process) // Mapping of port:process

// 	// Collect the processes
// 	for _, portStr := range os.Args[1:] {
// 		port, err := strconv.Atoi(portStr)
// 		if err != nil {
// 			log.Fatalf("Invalid port number: %s", portStr)
// 		}

// 		wg.Add(1)
// 		go func(p int) {
// 			defer wg.Done()

// 			// Get the processes
// 			output, err := exec.Command("powershell", "-Command", fmt.Sprintf("netstat -ano | findstr :%d", p)).CombinedOutput()

// 			if err != nil {
// 				log.Printf("Error getting process on port %d: %v. Output: %s", p, err, string(output))
// 				return
// 			}

// 			if len(output) > 0 {
// 				lines := strings.Split(strings.TrimSpace(string(output)), "\n")
// 				for _, line := range lines {
// 					parts := strings.Fields(line)
// 					processId := parts[len(parts)-1] // The process id is the last item on the line

// 					// Get the process name
// 					nameOutput, err := exec.Command("powershell", "-Command", fmt.Sprintf("(Get-Process -Id %s).Name", processId)).CombinedOutput()
// 					if err != nil {
// 						log.Printf("Error getting process name for PID %s on port %d: %v. Output: %s", processId, p, err, string(nameOutput))
// 						return
// 					}
// 					processName := strings.TrimSpace(string(nameOutput))

// 					key := fmt.Sprintf("%d:%s", p, processId) // Key is port:process
// 					processesToKill[key] = Process{Id: processId, Name: processName}
// 				}
// 			} else {
// 				log.Printf("No process found running on port %d", p)
// 			}
// 		}(port)
// 	}

// 	wg.Wait()

// 	// Now kill the processes
// 	for key, process := range processesToKill {
// 		err := exec.Command("powershell", "-Command", fmt.Sprintf("Stop-Process -Id %s -Force", process.Id)).Run()
// 		if err != nil {
// 			log.Printf("Error terminating process PID %s (name: %s) on port %s: %v", process.Id, process.Name, strings.Split(key, ":")[0], err)
// 		} else {
// 			log.Printf("Terminated process PID %s (name: %s) on port %s", process.Id, process.Name, strings.Split(key, ":")[0])
// 		}
// 	}
// }
