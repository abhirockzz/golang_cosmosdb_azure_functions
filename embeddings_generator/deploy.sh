# Set environment variables
prefix=$RANDOM
export LOCATION=enter the location e.g. eastus
export RG_NAME=${prefix}_cosmosdb_trigger_function_app_golang_rg
export STORAGE_ACC_NAME=${prefix}gotriggerapp
export STORAGE_SKU=Standard_LRS

export FUNCTION_APP_PLAN_NAME=${prefix}_cosmosdb_trigger_functionapp_golang_plan
export FUNCTION_APP_NAME=${prefix}-cosmosdb-embedding-function-app-golang
export FUNCTION_APP_PLAN_SKU=EP1

echo "Resource group $RG_NAME"

# Create resource group and storage account
az group create --name $RG_NAME --location $LOCATION
echo "Resource group $RG_NAME created."  

az storage account create --name $STORAGE_ACC_NAME --location $LOCATION --resource-group $RG_NAME --sku $STORAGE_SKU
echo "Storage account $STORAGE_ACC_NAME created."

# Create function app plan
az functionapp plan create --name $FUNCTION_APP_PLAN_NAME --resource-group $RG_NAME --location $LOCATION --sku $FUNCTION_APP_PLAN_SKU
echo "Function app plan $FUNCTION_APP_PLAN_NAME created."

# Create function app
az functionapp create --name $FUNCTION_APP_NAME --storage-account $STORAGE_ACC_NAME --plan $FUNCTION_APP_PLAN_NAME --resource-group $RG_NAME --functions-version 4 --runtime custom
echo "Function app $FUNCTION_APP_NAME created."

# Build Go binary for Windows
GOOS=windows GOARCH=amd64 go build -o main.exe main.go

# Publish function app (respond "no" to AzureWebJobsStorage overwrite prompt)
func azure functionapp publish $FUNCTION_APP_NAME --publish-local-settings
echo "Function app $FUNCTION_APP_NAME published."

export PRINCIPAL_ID=$(az webapp identity assign --resource-group $RG_NAME --name $FUNCTION_APP_NAME --query "principalId" -o tsv)

export COSMOSDB_ACCOUNT=enter the cosmos db account name
export COSMOSDB_RG_NAME=enter the resource group name of the cosmos db account

export COSMOSDB_ACC_ID=$(az cosmosdb show --name $COSMOSDB_ACCOUNT --resource-group $COSMOSDB_RG_NAME --query "id" -o tsv)

az cosmosdb sql role assignment create -n "Cosmos DB Built-in Data Contributor" -g $COSMOSDB_RG_NAME -a $COSMOSDB_ACCOUNT -p $PRINCIPAL_ID --scope $COSMOSDB_ACC_ID

echo "Cosmos DB role assignment created for $PRINCIPAL_ID"

export AZURE_OPENAI_RESOURCE_NAME=enter the azure openai resource name
export AZURE_OPENAI_RESOURCE_GROUP=enter the resource group name of the azure openai resource

AZURE_OPENAI_ID=$(az cognitiveservices account show --name $AZURE_OPENAI_RESOURCE_NAME --resource-group $AZURE_OPENAI_RESOURCE_GROUP --query "id" -o tsv)

az role assignment create --assignee $PRINCIPAL_ID --role "Cognitive Services OpenAI Contributor" --scope $AZURE_OPENAI_ID

echo "OpenAI role assignment created for $PRINCIPAL_ID"

echo "setup completed in resource group $RG_NAME"