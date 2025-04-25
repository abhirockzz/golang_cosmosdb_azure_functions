// Package common provides shared functionality for processing Cosmos DB documents.
package common

type CosmosDBDocument struct {
	ID            string `json:"id"`
	CustomerNotes string `json:"customerNotes"`
	Rid           string `json:"_rid"`
	Self          string `json:"_self"`
	Etag          string `json:"_etag"`
	Attachments   string `json:"_attachments"`
	Ts            int64  `json:"_ts"`
	Lsn           int    `json:"_lsn"`
}

// Data represents the data field in the Cosmos DB trigger payload.
type Data struct {
	Documents string `json:"documents"`
}

type SysMetadata struct {
	MethodName string `json:"MethodName"`
	UtcNow     string `json:"UtcNow"`
	RandGuid   string `json:"RandGuid"`
}

type Metadata struct {
	Sys SysMetadata `json:"sys"`
}

// CosmosDBTriggerPayload represents the structure of the Cosmos DB trigger payload.
type CosmosDBTriggerPayload struct {
	Data     Data     `json:"Data"`
	Metadata Metadata `json:"Metadata"`
}

// InvokeResponse represents the structure of the response returned by the handler.
type InvokeResponse struct {
	Outputs     map[string]any `json:"outputs"`
	Logs        []string       `json:"logs"`
	ReturnValue any            `json:"returnValue,omitempty"`
}
