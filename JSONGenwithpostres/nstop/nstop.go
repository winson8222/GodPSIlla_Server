package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func main() {

	osType := runtime.GOOS

	if osType == "windows" {
		opath, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			return
		}

		npath, exist := os.LookupEnv("NGINX_PATH")

		if !exist {
			log.Fatal("Nginx not installed Or folder not set to NGINX_PATH")
			return
		}

		err = os.Chdir(npath)
		if err != nil {
			log.Fatal("move to nginx folder failed")
		}

		err = os.Chdir(opath)
		if err != nil {
			log.Fatal("return to original location failed")
		}

		fmt.Println("Nginx is stopped...")
		// Additional code here

		// Allow some time for the NGINX server to start up
		time.Sleep(2 * time.Second)

		cmd := exec.Command("nginx", "-s", "stop")

		err = cmd.Start()

		if err != nil {
			// Print the captured logs
			fmt.Println("Error running Nginx:", err)
			return
		} else {
			fmt.Println("Error running Nginx: OS not found")
		}
	} else if osType == "darwin" {
		cmd := exec.Command("nginx", "-s", "quit")

		err := cmd.Start()
		if err != nil {
			log.Fatalf("Failed to stop NGINX server: %s", err)
		}

		fmt.Println("NGINX server stopped!")

		// Allow some time for the NGINX server to start up
		time.Sleep(2 * time.Second)
	}
}

// //for Windows
// func main() {

// 	err := os.Chdir("../nginx")
// 	if err != nil {
// 		log.Fatal("cannot move into gateway folder")
// 	}

// 	cmd := exec.Command("./nginx", "-s", "quit")

// 	err = cmd.Start()
// 	if err != nil {
// 		log.Fatalf("Failed to stop NGINX server: %s", err)
// 	}

// 	fmt.Println("NGINX server stopped!")

// 	// Allow some time for the NGINX server to start up
// 	time.Sleep(2 * time.Second)
// }
