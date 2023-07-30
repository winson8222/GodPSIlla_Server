package create

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// Modify the original IDL files to include the proper Routes (ServiceName/MethodNames)
func CreateIDL(constants Constants) {
	log.Printf("Opening input file...%s", constants.FilepathToService)
	inputFile, err := os.Open(constants.FilepathToService)
	if err != nil {
		log.Println("Error opening input file:", err)
		os.Exit(1)
	}

	log.Println("Creating temporary file...")
	tempFile, err := os.Create("../idl/temp.thrift")
	if err != nil {
		log.Println("Error creating temporary file:", err)
		inputFile.Close() // Close the inputFile if the tempFile creation fails
		os.Exit(1)
	}

	scanner := bufio.NewScanner(inputFile)
	writer := bufio.NewWriter(tempFile)

	log.Println("Starting to scan input file...")

	serviceFlag := false
	methodIndex := 0

	log.Println("Starting to scan input file...")
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		if strings.HasPrefix(trimmedLine, "service "+constants.ServiceName) {
			serviceFlag = true
			methodIndex = 0
			fmt.Fprintln(writer, line)
			continue
		}

		if serviceFlag && strings.HasPrefix(trimmedLine, "}") {
			serviceFlag = false
		}

		if serviceFlag && trimmedLine != "" {
			// Get the method name for the current methodIndex
			method := constants.Methods[methodIndex]

			// Remove the semicolon at the end of the line
			line = strings.TrimRight(line, ";")

			// Add the (api.post="/ServiceName/MethodName") annotation with a semicolon
			line = line + fmt.Sprintf(" (api.post=\"/%s/%s\");", constants.ServiceName, method.MethodName)

			// Increment the methodIndex for the next line
			methodIndex = (methodIndex + 1) % len(constants.Methods)
		}

		fmt.Fprintln(writer, line)
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error scanning input file:", err)
		inputFile.Close()
		tempFile.Close()
		os.Exit(1)
	}

	log.Println("Flushing the writer...")
	if err := writer.Flush(); err != nil {
		log.Println("Error flushing writer:", err)
		inputFile.Close()
		tempFile.Close()
		os.Exit(1)
	}

	// Close the files before renaming them
	log.Println("Closing input file...")
	if err := inputFile.Close(); err != nil {
		log.Println("Error closing input file:", err)
		os.Exit(1)
	}

	log.Println("Closing temporary file...")
	if err := tempFile.Close(); err != nil {
		log.Println("Error closing temporary file:", err)
		os.Exit(1)
	}

	log.Println("Removing original file...")
	if err := os.Remove(constants.FilepathToService); err != nil {
		log.Println("Error removing original file:", err)
		os.Exit(1)
	}

	log.Println("Renaming temporary file...")
	if err := os.Rename("../idl/temp.thrift", constants.FilepathToService); err != nil {
		log.Println("Error renaming temporary file:", err)
		os.Exit(1)
	}

	log.Println("File has been updated.")
}
