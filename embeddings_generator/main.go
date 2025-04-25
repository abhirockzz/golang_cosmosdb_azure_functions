package main

import (
	"crypto/sha256"
	"embeddings_generator_function/common"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"os"
)

const defaultPort = "8080"

var (
	cosmosVectorPropertyName        string
	cosmosVectorPropertyToEmbedName string
	cosmosHashPropertyName          string
	logs                            []string
)

var keysToRemove = []string{
	cosmosVectorPropertyName,
	//cosmosHashPropertyName,
	"_rid",
	"_self",
	"_etag",
	"_attachments",
	"_lsn",
	"_ts",
}

func init() {
	logs = []string{}
	cosmosVectorPropertyName = os.Getenv("COSMOS_VECTOR_PROPERTY")
	cosmosVectorPropertyToEmbedName = os.Getenv("COSMOS_PROPERTY_TO_EMBED")
	cosmosHashPropertyName = os.Getenv("COSMOS_HASH_PROPERTY")
}

func main() {
	addr := ":" + defaultPort
	http.HandleFunc("/cosmosdbprocessor", EmbeddingHandler)

	if port := os.Getenv("FUNCTIONS_CUSTOMHANDLER_PORT"); port != "" {
		addr = ":" + port
	}
	log.Printf("Server starting on address %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// EmbeddingHandler processes incoming Cosmos DB documents and generates embeddings for them.
func EmbeddingHandler(w http.ResponseWriter, req *http.Request) {
	logs = []string{"function invoked"}
	logs = append(logs,
		fmt.Sprintf("cosmosVectorPropertyName: %s", cosmosVectorPropertyName),
		fmt.Sprintf("cosmosVectorPropertyToEmbedName: %s", cosmosVectorPropertyToEmbedName),
		fmt.Sprintf("cosmosHashPropertyName: %s", cosmosHashPropertyName),
	)

	payloadBytes, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	documents, err := common.Parse(payloadBytes)
	if err != nil {
		log.Printf("Failed to parse payload: %v", err)
		http.Error(w, fmt.Sprintf("Failed to parse payload: %v", err), http.StatusBadRequest)
		return
	}

	logs = append(logs, fmt.Sprintf("Processing %d documents", len(documents)))

	var outputDocuments []map[string]any
	for _, doc := range documents {
		docID := doc["id"].(string)
		logs = append(logs,
			fmt.Sprintf("Processing document ID: %s", docID),
			fmt.Sprintf("Document data: %s", doc[cosmosVectorPropertyToEmbedName].(string)),
		)

		isNew, hashValue := isDocumentNewOrModified(doc, cosmosHashPropertyName, cosmosVectorPropertyToEmbedName)
		logs = append(logs, fmt.Sprintf("Document modification status: %t, hash: %s", isNew, hashValue))

		if isNew {
			// Cleanse the document of system properties
			doc = cleanse(doc, keysToRemove)

			docWithEmbedding, err := process(doc, hashValue, common.CreateEmbedding)
			if err != nil {
				log.Printf("Failed to process document %s: %v", docID, err)
				http.Error(w, fmt.Sprintf("Failed to process document: %v", err), http.StatusInternalServerError)
				return
			}

			outputDocuments = append(outputDocuments, docWithEmbedding)
		}
	}

	output := map[string]any{}
	if len(outputDocuments) > 0 {
		logs = append(logs, fmt.Sprintf("Adding %d documents with embeddings", len(outputDocuments)))
		output["outputData"] = outputDocuments
		logs = append(logs, "Added enriched documents to binding output")
	}

	response := common.InvokeResponse{
		Outputs:     output,
		Logs:        logs,
		ReturnValue: nil,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		http.Error(w, "Failed to generate response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

// process generates embeddings for a document and adds them along with a hash value.
func process(doc map[string]any, hashValue string, createEmbedding func(input string) ([]float32, error)) (map[string]any, error) {
	result := maps.Clone(doc)

	embedding, err := createEmbedding(doc[cosmosVectorPropertyToEmbedName].(string))
	if err != nil {
		return nil, fmt.Errorf("failed to create embedding: %w", err)
	}

	logs = append(logs, fmt.Sprintf("Created embedding for document: %v", doc))
	result[cosmosVectorPropertyName] = embedding
	result[cosmosHashPropertyName] = hashValue

	return result, nil
}

// cleanse removes specified keys from a document.
func cleanse(doc map[string]any, keysToRemove []string) map[string]any {
	for _, key := range keysToRemove {
		delete(doc, key)
	}
	return doc
}

// isDocumentNewOrModified checks if a document is new or has been modified.
func isDocumentNewOrModified(doc map[string]any, hashPropertyName, propertyToEmbedName string) (bool, string) {
	if _, exists := doc[hashPropertyName]; !exists {
		newHash := computeJSONHash(doc, propertyToEmbedName)
		logs = append(logs, "New document detected, generated hash: "+newHash)
		return true, newHash
	}

	existingHash, ok := doc[hashPropertyName].(string)
	if !ok {
		logs = append(logs, "Invalid hash property in document")
		return false, ""
	}

	hash := computeJSONHash(doc, propertyToEmbedName)
	if hash != existingHash {
		logs = append(logs, fmt.Sprintf("Document modified - old hash: %s, new hash: %s", existingHash, hash))
		return true, hash
	}

	logs = append(logs, "Document unchanged, hash: "+existingHash)
	return false, ""
}

// computeJSONHash generates a SHA256 hash of a document's specified property.
func computeJSONHash(doc map[string]any, propertyName string) string {
	property, ok := doc[propertyName].(string)
	if !ok {
		return ""
	}

	hash := sha256.Sum256([]byte(property))
	return hex.EncodeToString(hash[:])
}
