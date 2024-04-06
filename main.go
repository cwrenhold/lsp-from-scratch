package main

import (
	"bufio"
	"cwrenhold/lsp-from-scratch/lsp"
	"cwrenhold/lsp-from-scratch/rpc"
	"encoding/json"
	"log"
	"os"
)

func main() {
	logger := getLogger("lsp.txt")
	logger.Println("Starting lsp...")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)

		if err != nil {
			logger.Println("Error decoding message:", err)
			continue
		}

		handleMessage(logger, method, contents)
	}
}

func handleMessage(logger *log.Logger, method string, contents []byte) {
	logger.Printf("Received message with method: %s", method)

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Error unmarshalling initialize request: %s", err)
		}

		logger.Printf("Connected to %s %s",
			request.Params.ClientInfo.Name,
			request.Params.ClientInfo.Version)

		// Initialise a response
		msg := lsp.NewInitializeResponse(request.ID)
		reply := rpc.EncodeMessage(msg)

		writer := os.Stdout
		writer.Write([]byte(reply))

		logger.Print("Sent initialize response")
	}
}

func getLogger(filename string) *log.Logger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)

	if err != nil {
		panic("No file was found")
	}

	return log.New(logFile, "[educationallsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
