package scaffold

import (
	"fmt"
	"github.com/azure/azure-dev/cli/azd/internal"
	"strings"
)

type InfraSpec struct {
	Parameters []Parameter
	Services   []ServiceSpec

	// Databases to create
	DbPostgres    *DatabasePostgres
	DbMySql       *DatabaseMySql
	DbRedis       *DatabaseRedis
	DbCosmosMongo *DatabaseCosmosMongo
	DbCosmos      *DatabaseCosmosAccount

	// ai models
	AIModels []AIModel

	AzureServiceBus     *AzureDepServiceBus
	AzureEventHubs      *AzureDepEventHubs
	AzureStorageAccount *AzureDepStorageAccount
}

type Parameter struct {
	Name   string
	Value  any
	Type   string
	Secret bool
}

type DatabasePostgres struct {
	DatabaseUser string
	DatabaseName string
	AuthType     internal.AuthType
}

type DatabaseMySql struct {
	DatabaseUser string
	DatabaseName string
	AuthType     internal.AuthType
}

type CosmosSqlDatabaseContainer struct {
	ContainerName     string
	PartitionKeyPaths []string
}

type DatabaseCosmosAccount struct {
	DatabaseName string
	Containers   []CosmosSqlDatabaseContainer
}

type DatabaseCosmosMongo struct {
	DatabaseName string
}

type DatabaseRedis struct {
}

// AIModel represents a deployed, ready to use AI model.
type AIModel struct {
	Name  string
	Model AIModelModel
}

// AIModelModel represents a model that backs the AIModel.
type AIModelModel struct {
	// The name of the underlying model.
	Name string
	// The version of the underlying model.
	Version string
}

type AzureDepServiceBus struct {
	Queues                 []string
	TopicsAndSubscriptions map[string][]string
	AuthType               internal.AuthType
	IsJms                  bool
}

type AzureDepEventHubs struct {
	EventHubNames     []string
	AuthType          internal.AuthType
	UseKafka          bool
	SpringBootVersion string
}

type AzureDepStorageAccount struct {
	ContainerNames []string
	AuthType       internal.AuthType
}

type ServiceSpec struct {
	Name string
	Port int

	EnvironmentVariableInformation EnvironmentVariableInformation // todo delete this
	Envs                           []Env
	ResourceConnectionEnvs         []ResourceConnectionEnv

	// Front-end properties.
	Frontend *Frontend

	// Back-end properties
	Backend *Backend
}

type Env struct {
	Name  string
	Value string
}

type ResourceConnectionEnv struct {
	Name               string
	ResourceType       ResourceType
	ConnectionInfoType ConnectionInfoType
}

func (env *ResourceConnectionEnv) ToString() string {
	return fmt.Sprintf("ResourceConnectionEnv(Name=%s, ResourceType=%s, ConnectionInfoType=%s)",
		env.Name, env.ResourceType, env.ConnectionInfoType)
}

// todo merge ResourceType and project.ResourceType
// Not use project.ResourceType because it will cause cycle import.
// Not merge it in current PR to avoid conflict with upstream main branch.
// Solution proposal: define a ResourceType in lower level that can be used both in scaffold and project package.

type ResourceType string

const (
	ResourceTypeDbRedis             ResourceType = "db.redis"
	ResourceTypeDbPostgres          ResourceType = "db.postgres"
	ResourceTypeDbMySQL             ResourceType = "db.mysql"
	ResourceTypeDbMongo             ResourceType = "db.mongo"
	ResourceTypeDbCosmos            ResourceType = "db.cosmos"
	ResourceTypeHostContainerApp    ResourceType = "host.containerapp"
	ResourceTypeOpenAiModel         ResourceType = "ai.openai.model"
	ResourceTypeMessagingServiceBus ResourceType = "messaging.servicebus"
	ResourceTypeMessagingEventHubs  ResourceType = "messaging.eventhubs"
	ResourceTypeMessagingKafka      ResourceType = "messaging.kafka"
	ResourceTypeStorage             ResourceType = "storage"
)

type ConnectionInfoType string

const (
	Host             ConnectionInfoType = "host"
	Port             ConnectionInfoType = "port"
	Endpoint         ConnectionInfoType = "endpoint"
	DatabaseName     ConnectionInfoType = "databaseName"
	Namespace        ConnectionInfoType = "namespace"
	AccountName      ConnectionInfoType = "accountName"
	Username         ConnectionInfoType = "username"
	Password         ConnectionInfoType = "password"
	Url              ConnectionInfoType = "url"
	JdbcUrl          ConnectionInfoType = "jdbcUrl"
	ConnectionString ConnectionInfoType = "connectionString"
	IdentityClientId ConnectionInfoType = "identityClientId"
)

type EnvironmentVariableInformation struct {
	StringEnvironmentVariables    []StringEnvironmentVariable
	SecretRefEnvironmentVariables []SecretRefEnvironmentVariable
	ValueSecretDefinitions        []ValueSecretDefinition
	KeyVaultSecretDefinitions     []KeyVaultSecretDefinition
}

// StringEnvironmentVariable In generated bicep file, the Value will be quoted in bicep file.
// Example in bicep value:
//
//	value: 'jdbc:postgresql://${postgreServer.outputs.fqdn}:5432/${postgreSqlDatabaseName}'
type StringEnvironmentVariable struct {
	Name  string
	Value string
}

// SecretRefEnvironmentVariable In generated bicep file, the SecretRef will be quoted in bicep file.
// Example in bicep value:
//
//	secretRef: 'postgresql-password'
type SecretRefEnvironmentVariable struct {
	Name      string
	SecretRef string
}

// ValueSecretDefinition In generated bicep file, the SecretValue will be quoted in bicep file.
// Example in bicep value:
//
//	value: 'postgresql://${postgreSqlDatabaseUser}:${postgreSqlDatabasePassword}@${postgreServer.outputs.fqdn}:5432/${postgreSqlDatabaseName}'
type ValueSecretDefinition struct {
	SecretName  string
	SecretValue string
}

// KeyVaultSecretDefinition In generated bicep file, the KeyVaultUrl will be quoted in bicep file.
// Example in bicep value:
//
//	value: '${keyVault.outputs.uri}secrets/REDIS-URL'
type KeyVaultSecretDefinition struct {
	SecretName  string
	KeyVaultUrl string
}

type Frontend struct {
	Backends []ServiceReference
}

type Backend struct {
	Frontends []ServiceReference
}

type ServiceReference struct {
	Name string
}

type AIModelReference struct {
	Name string
}

func containerAppExistsParameter(serviceName string) Parameter {
	return Parameter{
		Name: BicepName(serviceName) + "Exists",
		Value: fmt.Sprintf("${SERVICE_%s_RESOURCE_EXISTS=false}",
			strings.ReplaceAll(strings.ToUpper(serviceName), "-", "_")),
		Type: "bool",
	}
}

type serviceDef struct {
	Settings []serviceDefSettings `json:"settings"`
}

type serviceDefSettings struct {
	Name         string `json:"name"`
	Value        string `json:"value"`
	Secret       bool   `json:"secret,omitempty"`
	SecretRef    string `json:"secretRef,omitempty"`
	CommentName  string `json:"_comment_name,omitempty"`
	CommentValue string `json:"_comment_value,omitempty"`
}

func serviceDefPlaceholder(serviceName string) Parameter {
	return Parameter{
		Name: BicepName(serviceName) + "Definition",
		Value: serviceDef{
			Settings: []serviceDefSettings{
				{
					Name:        "",
					Value:       "${VAR}",
					CommentName: "The name of the environment variable when running in Azure. If empty, ignored.",
					//nolint:lll
					CommentValue: "The value to provide. This can be a fixed literal, or an expression like ${VAR} to use the value of 'VAR' from the current environment.",
				},
				{
					Name:        "",
					Value:       "${VAR_S}",
					Secret:      true,
					CommentName: "The name of the environment variable when running in Azure. If empty, ignored.",
					//nolint:lll
					CommentValue: "The value to provide. This can be a fixed literal, or an expression like ${VAR_S} to use the value of 'VAR_S' from the current environment.",
				},
			},
		},
		Type:   "object",
		Secret: true,
	}
}
