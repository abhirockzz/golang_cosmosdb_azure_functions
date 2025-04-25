package main

import (
	"cosmosdb_go_function_trigger/common"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const defaultPort = "8080"

func main() {
	addr := ":" + defaultPort

	http.HandleFunc("/processor", processAndLog)

	port := os.Getenv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if port != "" {
		addr = ":" + port
	}
	log.Println("using address", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func processAndLog(w http.ResponseWriter, req *http.Request) {

	// Initialize custom logs
	logs := []string{
		"processor function invoked...",
	}
	payloadBytes, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var triggerPayload common.CosmosDBTriggerPayload
	if err := json.Unmarshal(payloadBytes, &triggerPayload); err != nil {
		http.Error(w, "Failed to parse JSON payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	logs = append(logs, fmt.Sprintf("Raw event payload: %s", triggerPayload))

	// Unmarshal the `documents` field
	var documents []common.CosmosDBDocument
	var documentsRaw string

	// First, unmarshal `triggerPayload.Data.Documents` as a string
	if err := json.Unmarshal([]byte(triggerPayload.Data.Documents), &documentsRaw); err != nil {
		log.Println("error while unmarshaling Documents field as string", err)
		http.Error(w, "Failed to parse documents field: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Then, unmarshal the JSON string into the `documents` slice
	if err := json.Unmarshal([]byte(documentsRaw), &documents); err != nil {
		log.Println("error while unmarshaling string to Document[]", err)
		http.Error(w, "Failed to parse documents: "+err.Error(), http.StatusBadRequest)
		return
	}

	for _, doc := range documents {
		logs = append(logs, fmt.Sprintf("Cosmos DB document: %v", doc))
		//log.Println("Cosmos DB document:", doc.ID)
	}

	// Construct the response with logs
	invokeResponse := common.InvokeResponse{Outputs: nil, Logs: logs, ReturnValue: nil}
	responseJson, _ := json.Marshal(invokeResponse)

	// Set the response headers and write the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJson)
}
