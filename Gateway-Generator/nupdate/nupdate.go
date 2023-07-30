package nupdate

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// Runs Nginx Reload
func NReload() {
	if runtime.GOOS == "windows" {
		var cmd *exec.Cmd
		opath, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			return
		}

		err = os.Chdir("../nginx/conf/")
		if err != nil {
			log.Fatal("move to folder failed")
		}

		path, err := os.Getwd()
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

		cmd = exec.Command("./nginx.exe", "-s", "reload", "-c", path+"/nginx.conf")

		err = cmd.Start()
		if err != nil {
			// Print the captured logs
			fmt.Println("Error running Nginx:", err)
			return
		}

		err = os.Chdir(opath)
		if err != nil {
			log.Fatal("return to original location failed")
		}

		fmt.Println("NGINX server reloaded!")

		// Allow some time for the NGINX server to start up
		time.Sleep(2 * time.Second)

	} else if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		var cmd *exec.Cmd

		relativePath := "../nginx/conf/nginx.conf"
		absolutePath, err := filepath.Abs(relativePath)
		if err != nil {
			fmt.Println("Error getting absolute path:", err)
			return
		}
		fmt.Println("Absolute path:", absolutePath)

		cmd = exec.Command("nginx", "-s", "reload", "-c", absolutePath)

		err = cmd.Start()
		if err != nil {
			log.Fatalf("Failed to reload NGINX server: %s", err)
		}

		fmt.Println("NGINX server reloaded!")

		// Allow some time for the NGINX server to start up
		time.Sleep(2 * time.Second)
	}
}
