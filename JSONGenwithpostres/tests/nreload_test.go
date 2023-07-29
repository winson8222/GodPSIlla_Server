package test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"nupdate"
	"os"
	"os/exec"
	"runtime"
	"testing"
)

func isPortOpen(port int) bool {
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return false // Port is not in use, no process is running on it.
	}
	conn.Close()
	return true // Port is in use, a process is running on it.
}

// Test to see if the nginx is reloaded with same nignx conf on port 80
func TestReload(t *testing.T) {
	err := os.Chdir("../nstart")
	if err != nil {
		log.Fatal("Error occur")
	}

	var cmd *exec.Cmd
	var cmd2 *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("./nstart.exe")
		cmd2 = exec.Command("./nstop.exe")
	}

	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		cmd = exec.Command("./nstartmac")
		cmd2 = exec.Command("./nstopmac")
	}

	err = cmd.Start()
	if err != nil {
		log.Fatalf("Failed to reload NGINX server: %s", err)
	}

	nupdate.NReload()

	port80 := isPortOpen(80)

	if port80 {
		err := os.Chdir("../nstop")
		if err != nil {
			log.Fatal("Error occur")
		}
		err = cmd2.Run()
		if err != nil {
			fmt.Print("gateway stopping failed")
		}
		t.Log("Port is up")
	} else {
		t.Error("Port is not up")
	}

}

// Test if the nignx server success successfully update with new config
func TestReload2(t *testing.T) {

	err := os.Chdir("../nstart")
	if err != nil {
		log.Fatal("Error occur")
	}

	var cmd *exec.Cmd
	var cmd2 *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("./nstart.exe")
		cmd2 = exec.Command("./nstop.exe")
	}

	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		cmd = exec.Command("./nstartmac")
		cmd2 = exec.Command("./nstopmac")
	}

	err = cmd.Start()
	if err != nil {
		log.Fatalf("Failed to reload NGINX server: %s", err)
	}

	nginxConfPath := "../nginx/conf/nginx.conf"
	secondConfPath := "../testfiles/nginx/second.conf"

	// Replace nginx.conf with the content of second.conf.
	err = backupAndReplaceFile(secondConfPath, nginxConfPath)
	if err != nil {
		log.Fatal("Error replacing nginx.conf with second.conf:", err)
	}
	log.Println("nginx.conf has been replaced with second.conf successfully.")

	defer nupdate.NReload()

	url := "http://0.0.0.0:80/Echo/Echo"

	// Create an HTTP client.
	client := http.Client{}

	// Send an HTTP GET request to the URL.
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Error connecting to %s: %s\n", url, err)
		return
	}
	defer resp.Body.Close()

	// Check the response status code.
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusBadGateway {
		fmt.Printf("Successfully connected to %s\n", url)
	} else {
		err := os.Chdir("../nstop")
		if err != nil {
			log.Fatal("Error occur")
		}
		err = cmd2.Run()
		if err != nil {
			fmt.Print("gateway stopping failed")
		}
		t.Errorf("Failed to connect to %s. Status Code: %d\n", url, resp.StatusCode)
	}

	// Restore nginx.conf back to the original file.
	backupPath := nginxConfPath + ".backup"
	err = backupAndReplaceFile(backupPath, nginxConfPath)
	if err != nil {
		log.Fatal("Error restoring nginx.conf to the original file:", err)
	}
	log.Println("nginx.conf has been restored to the original file successfully.")

	port80 := isPortOpen(80)

	if port80 {
		err := os.Chdir("../nstop")
		if err != nil {
			log.Fatal("Error occur")
		}
		err = cmd2.Run()
		if err != nil {
			fmt.Print("gateway stopping failed")
		}
		t.Log("Port is up")
	} else {
		t.Error("Port is not up")
	}

}

func backupAndReplaceFile(sourcePath string, destinationPath string) error {
	// Check if the backup file exists.
	backupPath := destinationPath + ".backup"
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		// If the backup file doesn't exist, create a new one.
		originalContent, err := ioutil.ReadFile(destinationPath)
		if err != nil {
			log.Fatal("Error reading nginx.conf:", err)
			return err

		}

		err = ioutil.WriteFile(backupPath, originalContent, 0644)
		if err != nil {
			log.Fatal("Error creating backup:", err)
			return err

		}
	}

	// Read the content from the source file (second.conf).
	sourceContent, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		log.Fatal("Error reading second.conf:", err)
		return err
	}

	// Replace the destination file (nginx.conf) with the content of the source file (second.conf).
	err = ioutil.WriteFile(destinationPath, sourceContent, 0644)
	if err != nil {
		log.Fatal("Error replacing nginx.conf with second.conf:", err)
		return err
	}

	return err
}
