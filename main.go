package main

import (
	"bufio"
	"cwrenhold/lsp-from-scratch/rpc"
	"log"
	"os"
)

func main() {
	logger := getLogger("lsp.txt")
	logger.Println("Starting lsp...")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	for scanner.Scan() {
		msg := scanner.Text()
		handleMessage(logger, msg)
	}
}

func handleMessage(logger *log.Logger, msg any) {
	logger.Println("Received message:", msg)
}

func getLogger(filename string) *log.Logger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)

	if err != nil {
		panic("No file was found")
	}

	return log.New(logFile, "[educationallsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
