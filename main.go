package main

import (
	"bufio"
	"cwrenhold/lsp-from-scratch/analysis"
	"cwrenhold/lsp-from-scratch/lsp"
	"cwrenhold/lsp-from-scratch/rpc"
	"encoding/json"
	"io"
	"log"
	"os"
)

func main() {
	logger := getLogger("lsp.txt")
	logger.Println("Starting lsp...")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	state := analysis.NewState()
	writer := os.Stdout

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)

		if err != nil {
			logger.Println("Error decoding message:", err)
			continue
		}

		handleMessage(logger, writer, state, method, contents)
	}
}

func handleMessage(logger *log.Logger, writer io.Writer, state analysis.State, method string, contents []byte) {
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

		writeResponse(writer, msg)

		logger.Print("Sent initialize response")
	case "textDocument/didOpen":
		var notification lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &notification); err != nil {
			logger.Printf("Error unmarshalling textdocument/didOpen request: %s", err)
		}

		logger.Printf("Opened document %s", notification.Params.TextDocument.URI)
		diagnostics := state.OpenDocument(notification.Params.TextDocument.URI, notification.Params.TextDocument.Text)

		if len(diagnostics) > 0 {
			writeResponse(writer, lsp.PublishDiagnosticsNotification{
				Notification: lsp.Notification{
					RPC:    "2.0",
					Method: "textDocument/publishDiagnostics",
				},
				Params: lsp.PublishDiagnosticsParams{
					URI:         notification.Params.TextDocument.URI,
					Diagnostics: diagnostics,
				},
			})
		}
	case "textDocument/didChange":
		var notification lsp.TextDocumentDidChangeNotification
		if err := json.Unmarshal(contents, &notification); err != nil {
			logger.Printf("Error unmarshalling textdocument/didChange request: %s", err)
		}

		logger.Printf("Changed: %s", notification.Params.TextDocument.URI)
		for _, change := range notification.Params.ContentChanges {
			diagnostics := state.UpdateDocument(notification.Params.TextDocument.URI, change.Text)

			if len(diagnostics) > 0 {
				writeResponse(writer, lsp.PublishDiagnosticsNotification{
					Notification: lsp.Notification{
						RPC:    "2.0",
						Method: "textDocument/publishDiagnostics",
					},
					Params: lsp.PublishDiagnosticsParams{
						URI:         notification.Params.TextDocument.URI,
						Diagnostics: diagnostics,
					},
				})
			}
		}
	case "textDocument/hover":
		var request lsp.HoverRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Error unmarshalling hover request: %s", err)
		}

		// Create response
		response := state.Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position)

		writeResponse(writer, response)
	// Triggered by :lua vim.lsp.buf.definition(), this is because some plugins seem to prevent this working out of the box?
	case "textDocument/definition":
		var request lsp.DefinitionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Error unmarshalling definition request: %s", err)
		}

		// Create response
		response := state.Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position)

		writeResponse(writer, response)
	// Triggered by :lua vim.lsp.buf.code_action()
	case "textDocument/codeAction":
		var request lsp.CodeActionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Error unmarshalling code action request: %s", err)
		}

		response := state.TextDocumentCodeAction(request.ID, request.Params.TextDocument.URI)

		writeResponse(writer, response)
	case "textDocument/completion":
		var request lsp.CompletionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Error unmarshalling completion request: %s", err)
		}

		response := state.TextDocumentCompletion(request.ID, request.Params.TextDocument.URI)

		writeResponse(writer, response)
	}
}

func writeResponse(writer io.Writer, msg any) {
	reply := rpc.EncodeMessage(msg)
	writer.Write([]byte(reply))
}

func getLogger(filename string) *log.Logger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)

	if err != nil {
		panic("No file was found")
	}

	return log.New(logFile, "[educationallsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
