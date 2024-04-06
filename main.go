package main

import (
	"bufio"
	"cwrenhold/lsp-from-scratch/analysis"
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

	state := analysis.NewState()

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)

		if err != nil {
			logger.Println("Error decoding message:", err)
			continue
		}

		handleMessage(logger, state, method, contents)
	}
}

func handleMessage(logger *log.Logger, state analysis.State, method string, contents []byte) {
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
	case "textDocument/didOpen":
		var notification lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &notification); err != nil {
			logger.Printf("Error unmarshalling textdocument/didOpen request: %s", err)
		}

		logger.Printf("Opened document %s", notification.Params.TextDocument.URI)
		state.OpenDocument(notification.Params.TextDocument.URI, notification.Params.TextDocument.Text)
	case "textDocument/didChange":
		var notification lsp.TextDocumentDidChangeNotification
		if err := json.Unmarshal(contents, &notification); err != nil {
			logger.Printf("Error unmarshalling textdocument/didChange request: %s", err)
		}

		logger.Printf("Changed: %s", notification.Params.TextDocument.URI)
		for _, change := range notification.Params.ContentChanges {
			state.UpdateDocument(notification.Params.TextDocument.URI, change.Text)
		}
	}
}

func getLogger(filename string) *log.Logger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)

	if err != nil {
		panic("No file was found")
	}

	return log.New(logFile, "[educationallsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
