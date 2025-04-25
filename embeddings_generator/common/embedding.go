// Package common provides shared functionality for processing Cosmos DB documents.
package common

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

var client *azopenai.Client

func init() {
	var err error
	client, err = getOpenAIClient()
	if err != nil {
		log.Fatalf("Failed to create OpenAI client: %v", err)
	}
}

// CreateEmbedding generates an embedding for the given input text using Azure OpenAI.
func CreateEmbedding(input string) ([]float32, error) {
	modelDeploymentID := os.Getenv("OPENAI_DEPLOYMENT_NAME")
	if modelDeploymentID == "" {
		return nil, errors.New("OPENAI_DEPLOYMENT_NAME environment variable not set")
	}

	resp, err := client.GetEmbeddings(context.Background(), azopenai.EmbeddingsOptions{
		Input:          []string{input},
		DeploymentName: &modelDeploymentID,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, errors.New("no embedding data received from OpenAI")
	}

	return resp.Data[0].Embedding, nil
}

// getOpenAIClient creates and returns an Azure OpenAI client using default Azure credentials.
func getOpenAIClient() (*azopenai.Client, error) {
	azureOpenAIEndpoint := os.Getenv("OPENAI_ENDPOINT")
	modelDeploymentID := os.Getenv("OPENAI_DEPLOYMENT_NAME")

	if modelDeploymentID == "" || azureOpenAIEndpoint == "" {
		return nil, errors.New("required environment variables OPENAI_ENDPOINT and OPENAI_DEPLOYMENT_NAME must be set")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create default Azure credential: %w", err)
	}

	client, err := azopenai.NewClient(azureOpenAIEndpoint, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI client: %w", err)
	}

	return client, nil
}
