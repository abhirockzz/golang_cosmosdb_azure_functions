// Package common provides shared functionality for processing Cosmos DB documents.
package common

import (
	"encoding/json"
	"fmt"
	"log"
)

// Parse unmarshals the Cosmos DB trigger payload and extracts the documents.
// It performs a two-step unmarshaling process due to the nested JSON structure.
// This generic function allows you to specify the type T that the documents should be unmarshaled to.
func Parse[T any](payloadBytes []byte) ([]T, error) {
	var triggerPayload CosmosDBTriggerPayload
	if err := json.Unmarshal(payloadBytes, &triggerPayload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal trigger payload: %w", err)
	}

	// First unmarshal step: convert the Documents field from string to []byte
	var documentsRaw string
	if err := json.Unmarshal([]byte(triggerPayload.Data.Documents), &documentsRaw); err != nil {
		log.Printf("Failed to unmarshal Documents field as string: %v", err)
		return nil, fmt.Errorf("failed to unmarshal Documents field: %w", err)
	}

	// Second unmarshal step: convert the JSON string to []T
	var documents []T
	if err := json.Unmarshal([]byte(documentsRaw), &documents); err != nil {
		log.Printf("Failed to unmarshal documents array: %v", err)
		return nil, fmt.Errorf("failed to unmarshal documents array: %w", err)
	}

	return documents, nil
}

// // Parse unmarshals the Cosmos DB trigger payload and extracts the documents.
// // It performs a two-step unmarshaling process due to the nested JSON structure.
// func Parse(payloadBytes []byte) ([]map[string]any, error) {
// 	var triggerPayload CosmosDBTriggerPayload
// 	if err := json.Unmarshal(payloadBytes, &triggerPayload); err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal trigger payload: %w", err)
// 	}

// 	// First unmarshal step: convert the Documents field from string to []byte
// 	var documentsRaw string
// 	if err := json.Unmarshal([]byte(triggerPayload.Data.Documents), &documentsRaw); err != nil {
// 		log.Printf("Failed to unmarshal Documents field as string: %v", err)
// 		return nil, fmt.Errorf("failed to unmarshal Documents field: %w", err)
// 	}

// 	// Second unmarshal step: convert the JSON string to []map[string]any
// 	var documents []map[string]any
// 	if err := json.Unmarshal([]byte(documentsRaw), &documents); err != nil {
// 		log.Printf("Failed to unmarshal documents array: %v", err)
// 		return nil, fmt.Errorf("failed to unmarshal documents array: %w", err)
// 	}

// 	return documents, nil
// }
