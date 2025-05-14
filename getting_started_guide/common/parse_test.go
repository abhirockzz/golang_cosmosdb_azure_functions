package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	payload := `{"Data":{"documents":"\"[{\\\"id\\\":\\\"dfa26d32-f876-44a3-b107-369f1f48c689\\\",\\\"customerNotes\\\":\\\"this is a great product\\\",\\\"_rid\\\":\\\"lV8dAK7u9cCUAAAAAAAAAA==\\\",\\\"_self\\\":\\\"dbs/lV8dAA==/colls/lV8dAK7u9cA=/docs/lV8dAK7u9cCUAAAAAAAAAA==/\\\",\\\"_etag\\\":\\\"\\\\\\\"0f007efc-0000-0800-0000-67f5fb920000\\\\\\\"\\\",\\\"_attachments\\\":\\\"attachments/\\\",\\\"_ts\\\":1744173970,\\\"_lsn\\\":160}]\""},"Metadata":{"sys":{"MethodName":"cosmosdbprocessor","UtcNow":"2025-04-09T04:46:10.723203Z","RandGuid":"0d00378b-6426-4af1-9fc0-0793f4ce3745"}}}`

	result, err := Parse[CosmosDBDocument]([]byte(payload))
	assert.NoError(t, err, "Parse function returned an error")
	assert.Len(t, result, 1, "expected 1 document")

	doc := result[0]
	assert.Equal(t, "dfa26d32-f876-44a3-b107-369f1f48c689", doc.ID, "expected id to match")
	assert.Equal(t, "this is a great product", doc.CustomerNotes, "expected customerNotes to match")

}

func TestParseMultipleDocuments(t *testing.T) {
	payload := `{"Data":{"documents":"\"[{\\\"id\\\":\\\"51e0c1b0-87d3-4611-ac41-7ac3e77d9920\\\",\\\"customerNotes\\\":\\\"Schedule team meeting\\\",\\\"_rid\\\":\\\"lV8dAK7u9cCVAAAAAAAAAA==\\\",\\\"_self\\\":\\\"dbs/lV8dAA==/colls/lV8dAK7u9cA=/docs/lV8dAK7u9cCVAAAAAAAAAA==/\\\",\\\"_etag\\\":\\\"\\\\\\\"0f00a3fd-0000-0800-0000-67f5fc640000\\\\\\\"\\\",\\\"_attachments\\\":\\\"attachments/\\\",\\\"_ts\\\":1744174180,\\\"_lsn\\\":161},{\\\"id\\\":\\\"cfbf42b9-48e8-449b-9cff-17c6fbd00f83\\\",\\\"customerNotes\\\":\\\"Update dependencies\\\",\\\"_rid\\\":\\\"lV8dAK7u9cCWAAAAAAAAAA==\\\",\\\"_self\\\":\\\"dbs/lV8dAA==/colls/lV8dAK7u9cA=/docs/lV8dAK7u9cCWAAAAAAAAAA==/\\\",\\\"_etag\\\":\\\"\\\\\\\"0f00a9fd-0000-0800-0000-67f5fc670000\\\\\\\"\\\",\\\"_attachments\\\":\\\"attachments/\\\",\\\"_ts\\\":1744174183,\\\"_lsn\\\":162}]\""},"Metadata":{"sys":{"MethodName":"cosmosdbprocessor","UtcNow":"2025-04-09T04:49:45.157601Z","RandGuid":"304980d9-584d-4323-98e3-b46bb1eebded"}}}`

	// Test using the helper function
	result, err := Parse[CosmosDBDocument]([]byte(payload))
	assert.NoError(t, err, "Parse function returned an error")
	assert.Len(t, result, 2, "expected 2 documents")

	doc1 := result[0]
	assert.Equal(t, "51e0c1b0-87d3-4611-ac41-7ac3e77d9920", doc1.ID, "expected id of first document to match")
	assert.Equal(t, "Schedule team meeting", doc1.CustomerNotes, "expected customerNotes of first document to match")

	doc2 := result[1]
	assert.Equal(t, "cfbf42b9-48e8-449b-9cff-17c6fbd00f83", doc2.ID, "expected id of second document to match")
	assert.Equal(t, "Update dependencies", doc2.CustomerNotes, "expected customerNotes of second document to match")
}

func TestParseToMapSlice(t *testing.T) {
	payload := `{"Data":{"documents":"\"[{\\\"id\\\":\\\"dfa26d32-f876-44a3-b107-369f1f48c689\\\",\\\"customerNotes\\\":\\\"this is a great product\\\",\\\"customField\\\":\\\"custom value\\\",\\\"_rid\\\":\\\"lV8dAK7u9cCUAAAAAAAAAA==\\\"}]\""},"Metadata":{"sys":{"MethodName":"cosmosdbprocessor","UtcNow":"2025-04-09T04:46:10.723203Z","RandGuid":"0d00378b-6426-4af1-9fc0-0793f4ce3745"}}}`

	result, err := Parse[map[string]any]([]byte(payload))
	assert.NoError(t, err, "Parse function returned an error")
	assert.Len(t, result, 1, "expected 1 document")

	doc := result[0]
	assert.Equal(t, "dfa26d32-f876-44a3-b107-369f1f48c689", doc["id"], "expected id to match")
	assert.Equal(t, "this is a great product", doc["customerNotes"], "expected customerNotes to match")
	assert.Equal(t, "custom value", doc["customField"], "expected customField to match")

}

func TestParseToMapSliceMultipleDocuments(t *testing.T) {
	payload := `{"Data":{"documents":"\"[{\\\"id\\\":\\\"51e0c1b0-87d3-4611-ac41-7ac3e77d9920\\\",\\\"customerNotes\\\":\\\"Schedule team meeting\\\",\\\"_rid\\\":\\\"lV8dAK7u9cCVAAAAAAAAAA==\\\",\\\"_self\\\":\\\"dbs/lV8dAA==/colls/lV8dAK7u9cA=/docs/lV8dAK7u9cCVAAAAAAAAAA==/\\\",\\\"_etag\\\":\\\"\\\\\\\"0f00a3fd-0000-0800-0000-67f5fc640000\\\\\\\"\\\",\\\"_attachments\\\":\\\"attachments/\\\",\\\"_ts\\\":1744174180,\\\"_lsn\\\":161},{\\\"id\\\":\\\"cfbf42b9-48e8-449b-9cff-17c6fbd00f83\\\",\\\"customerNotes\\\":\\\"Update dependencies\\\",\\\"_rid\\\":\\\"lV8dAK7u9cCWAAAAAAAAAA==\\\",\\\"_self\\\":\\\"dbs/lV8dAA==/colls/lV8dAK7u9cA=/docs/lV8dAK7u9cCWAAAAAAAAAA==/\\\",\\\"_etag\\\":\\\"\\\\\\\"0f00a9fd-0000-0800-0000-67f5fc670000\\\\\\\"\\\",\\\"_attachments\\\":\\\"attachments/\\\",\\\"_ts\\\":1744174183,\\\"_lsn\\\":162}]\""},"Metadata":{"sys":{"MethodName":"cosmosdbprocessor","UtcNow":"2025-04-09T04:49:45.157601Z","RandGuid":"304980d9-584d-4323-98e3-b46bb1eebded"}}}`

	// Test using the helper function
	result, err := Parse[map[string]any]([]byte(payload))
	assert.NoError(t, err, "Parse function returned an error")
	assert.Len(t, result, 2, "expected 2 documents")

	doc1 := result[0]
	assert.Equal(t, "51e0c1b0-87d3-4611-ac41-7ac3e77d9920", doc1["id"], "expected id of first document to match")
	assert.Equal(t, "Schedule team meeting", doc1["customerNotes"], "expected customerNotes of first document to match")

	doc2 := result[1]
	assert.Equal(t, "cfbf42b9-48e8-449b-9cff-17c6fbd00f83", doc2["id"], "expected id of second document to match")
	assert.Equal(t, "Update dependencies", doc2["customerNotes"], "expected customerNotes of second document to match")
}
