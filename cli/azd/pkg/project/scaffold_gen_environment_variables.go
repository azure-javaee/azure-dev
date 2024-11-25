package project

import (
	"fmt"
	"github.com/azure/azure-dev/cli/azd/internal"
	"github.com/azure/azure-dev/cli/azd/internal/scaffold"
	"strings"
)

func getResourceConnectionEnvs(usedResource *ResourceConfig,
	infraSpec *scaffold.InfraSpec) ([]scaffold.ResourceConnectionEnv, error) {
	resourceType := usedResource.Type
	authType, err := getAuthType(infraSpec, usedResource.Type)
	if err != nil {
		return []scaffold.ResourceConnectionEnv{}, err
	}
	switch resourceType {
	case ResourceTypeDbPostgres:
		switch authType {
		case internal.AuthTypePassword:
			return []scaffold.ResourceConnectionEnv{
				{
					Name:               "POSTGRES_USERNAME",
					ResourceType:       scaffold.ResourceTypeDbPostgres,
					ConnectionInfoType: scaffold.Username,
				},
				{
					Name:               "POSTGRES_PASSWORD",
					ResourceType:       scaffold.ResourceTypeDbPostgres,
					ConnectionInfoType: scaffold.Password,
				},
				{
					Name:               "POSTGRES_HOST",
					ResourceType:       scaffold.ResourceTypeDbPostgres,
					ConnectionInfoType: scaffold.Host,
				},
				{
					Name:               "POSTGRES_DATABASE",
					ResourceType:       scaffold.ResourceTypeDbPostgres,
					ConnectionInfoType: scaffold.DatabaseName,
				},
				{
					Name:               "POSTGRES_PORT",
					ResourceType:       scaffold.ResourceTypeDbPostgres,
					ConnectionInfoType: scaffold.Port,
				},
				{
					Name:               "POSTGRES_URL",
					ResourceType:       scaffold.ResourceTypeDbPostgres,
					ConnectionInfoType: scaffold.Url,
				},
				{
					Name:               "spring.datasource.url",
					ResourceType:       scaffold.ResourceTypeDbPostgres,
					ConnectionInfoType: scaffold.JdbcUrl,
				},
				{
					Name:               "spring.datasource.username",
					ResourceType:       scaffold.ResourceTypeDbPostgres,
					ConnectionInfoType: scaffold.Username,
				},
				{
					Name:               "spring.datasource.password",
					ResourceType:       scaffold.ResourceTypeDbPostgres,
					ConnectionInfoType: scaffold.Password,
				},
			}, nil
		case internal.AuthTypeUserAssignedManagedIdentity:
			// Some other environment variables are added by service connector,
			// should not add to bicep generation context
			return []scaffold.ResourceConnectionEnv{
				{
					Name:               "POSTGRES_USERNAME",
					ResourceType:       scaffold.ResourceTypeDbPostgres,
					ConnectionInfoType: scaffold.Username,
				},
				{
					Name:               "POSTGRES_HOST",
					ResourceType:       scaffold.ResourceTypeDbPostgres,
					ConnectionInfoType: scaffold.Host,
				},
				{
					Name:               "POSTGRES_DATABASE",
					ResourceType:       scaffold.ResourceTypeDbPostgres,
					ConnectionInfoType: scaffold.DatabaseName,
				},
				{
					Name:               "POSTGRES_PORT",
					ResourceType:       scaffold.ResourceTypeDbPostgres,
					ConnectionInfoType: scaffold.Port,
				},
			}, nil
		default:
			return []scaffold.ResourceConnectionEnv{}, unsupportedAuthTypeWhenGetResourceConnectionEnvsError(resourceType, authType)
		}
	default:
		return []scaffold.ResourceConnectionEnv{}, unsupportedResourceTypeWhenGetResourceConnectionEnvsError(resourceType)
	}
}

func getEnvironmentVariableInformation(usedResource *ResourceConfig,
	infraSpec *scaffold.InfraSpec) (scaffold.EnvironmentVariableInformation, error) {
	resourceType := usedResource.Type
	authType, err := getAuthType(infraSpec, usedResource.Type)
	if err != nil {
		return scaffold.EnvironmentVariableInformation{}, err
	}
	switch resourceType {
	case ResourceTypeDbMySQL:
		switch authType {
		case internal.AuthTypePassword:
			return scaffold.EnvironmentVariableInformation{
				StringEnvironmentVariables: []scaffold.StringEnvironmentVariable{
					{
						Name:  "MYSQL_USERNAME",
						Value: "${mysqlDatabaseUser}",
					},
					{
						Name:  "MYSQL_HOST",
						Value: "${mysqlServer.outputs.fqdn}",
					},
					{
						Name:  "MYSQL_DATABASE",
						Value: "${mysqlDatabaseName}",
					},
					{
						Name:  "MYSQL_PORT",
						Value: "3306",
					},
					{
						Name:  "spring.datasource.url",
						Value: "jdbc:mysql://${mysqlServer.outputs.fqdn}:3306/${mysqlDatabaseName}",
					},
					{
						Name:  "spring.datasource.username",
						Value: "${mysqlDatabaseUser}",
					},
				},
				SecretRefEnvironmentVariables: []scaffold.SecretRefEnvironmentVariable{
					{
						Name:      "MYSQL_URL",
						SecretRef: "mysql-db-url",
					},
					{
						Name:      "MYSQL_PASSWORD",
						SecretRef: "mysql-password",
					},
					{
						Name:      "spring.datasource.password",
						SecretRef: "mysql-password",
					},
				},
				ValueSecretDefinitions: []scaffold.ValueSecretDefinition{
					{
						SecretName:  "mysql-db-url",
						SecretValue: "mysql://${mysqlDatabaseUser}:${mysqlDatabasePassword}@${mysqlServer.outputs.fqdn}:3306/${mysqlDatabaseName}",
					},
					{
						SecretName:  "mysql-password",
						SecretValue: "${mysqlDatabasePassword}",
					},
				},
			}, nil
		case internal.AuthTypeUserAssignedManagedIdentity:
			return scaffold.EnvironmentVariableInformation{
				StringEnvironmentVariables: []scaffold.StringEnvironmentVariable{
					// Some other environment variables are added by service connector,
					// should not add to bicep generation context
					{
						Name:  "MYSQL_USERNAME",
						Value: "${mysqlDatabaseUser}",
					},
					{
						Name:  "MYSQL_HOST",
						Value: "${mysqlServer.outputs.fqdn}",
					},
					{
						Name:  "MYSQL_DATABASE",
						Value: "${mysqlDatabaseName}",
					},
					{
						Name:  "MYSQL_PORT",
						Value: "3306",
					},
				},
			}, nil
		default:
			return scaffold.EnvironmentVariableInformation{},
				unsupportedAuthTypeWhenGetResourceConnectionEnvsError(resourceType, authType)
		}
	case ResourceTypeDbRedis:
		switch authType {
		case internal.AuthTypePassword:
			return scaffold.EnvironmentVariableInformation{
				StringEnvironmentVariables: []scaffold.StringEnvironmentVariable{
					{
						Name:  "REDIS_HOST",
						Value: "${redis.outputs.hostName}",
					},
					{
						Name:  "REDIS_PORT",
						Value: "${redis.outputs.sslPort}",
					},
					{
						Name:  "REDIS_ENDPOINT",
						Value: "${redis.outputs.hostName}:${redis.outputs.sslPort}",
					},
				},
				SecretRefEnvironmentVariables: []scaffold.SecretRefEnvironmentVariable{
					{
						Name:      "REDIS_URL",
						SecretRef: "redis-url",
					},
					{
						Name:      "REDIS_PASSWORD",
						SecretRef: "redis-pass",
					},
					{
						Name:      "spring.data.redis.url",
						SecretRef: "redis-url",
					},
				},
				KeyVaultSecretDefinitions: []scaffold.KeyVaultSecretDefinition{
					{
						SecretName:  "redis-pass",
						KeyVaultUrl: "${keyVault.outputs.uri}secrets/REDIS-PASSWORD",
					},
					{
						SecretName:  "redis-url",
						KeyVaultUrl: "${keyVault.outputs.uri}secrets/REDIS-URL",
					},
				},
			}, nil
		default:
			return scaffold.EnvironmentVariableInformation{},
				unsupportedAuthTypeWhenGetResourceConnectionEnvsError(resourceType, authType)
		}
	case ResourceTypeDbMongo:
		switch authType {
		case internal.AuthTypeUserAssignedManagedIdentity:
			return scaffold.EnvironmentVariableInformation{
				StringEnvironmentVariables: []scaffold.StringEnvironmentVariable{
					{
						Name:  "spring.data.mongodb.database",
						Value: "${mongoDatabaseName}",
					},
				},
				SecretRefEnvironmentVariables: []scaffold.SecretRefEnvironmentVariable{
					{
						Name:      "MONGODB_URL",
						SecretRef: "mongodb-url",
					},
					{
						Name:      "spring.data.mongodb.uri",
						SecretRef: "mongodb-url",
					},
				},
				KeyVaultSecretDefinitions: []scaffold.KeyVaultSecretDefinition{
					{
						SecretName:  "mongodb-url",
						KeyVaultUrl: "${cosmos.outputs.exportedSecrets['MONGODB-URL'].secretUri}",
					},
				},
			}, nil
		default:
			return scaffold.EnvironmentVariableInformation{},
				unsupportedAuthTypeWhenGetResourceConnectionEnvsError(resourceType, authType)
		}
	case ResourceTypeDbCosmos:
		switch authType {
		case internal.AuthTypeUserAssignedManagedIdentity:
			return scaffold.EnvironmentVariableInformation{
				StringEnvironmentVariables: []scaffold.StringEnvironmentVariable{
					{
						Name:  "spring.cloud.azure.cosmos.endpoint",
						Value: "${cosmos.outputs.endpoint}",
					},
					{
						Name:  "spring.cloud.azure.cosmos.database",
						Value: "${cosmosDatabaseName}",
					},
				},
			}, nil
		default:
			return scaffold.EnvironmentVariableInformation{},
				unsupportedAuthTypeWhenGetResourceConnectionEnvsError(resourceType, authType)
		}
	case ResourceTypeMessagingServiceBus:
		if infraSpec.AzureServiceBus.IsJms {
			switch authType {
			case internal.AuthTypeUserAssignedManagedIdentity:
				return scaffold.EnvironmentVariableInformation{
					StringEnvironmentVariables: []scaffold.StringEnvironmentVariable{
						{
							Name:  "spring.jms.servicebus.passwordless-enabled",
							Value: "true",
						},
						{
							Name:  "spring.jms.servicebus.namespace",
							Value: "${serviceBusNamespace.outputs.name}",
						},
						{
							Name:  "spring.jms.servicebus.credential.managed-identity-enabled",
							Value: "true",
						},
						{
							Name:  "spring.jms.servicebus.credential.client-id",
							Value: "__PlaceHolderForServiceIdentityClientId",
						},
						{
							Name:  "spring.jms.servicebus.pricing-tier",
							Value: "premium",
						},
						{
							Name:  "spring.jms.servicebus.connection-string",
							Value: "",
						},
					},
				}, nil
			case internal.AuthTypeConnectionString:
				return scaffold.EnvironmentVariableInformation{
					StringEnvironmentVariables: []scaffold.StringEnvironmentVariable{
						{
							Name:  "spring.jms.servicebus.pricing-tier",
							Value: "premium",
						},
						{
							Name:  "spring.jms.servicebus.passwordless-enabled",
							Value: "false",
						},
						{
							Name:  "spring.jms.servicebus.namespace",
							Value: "",
						},
						{
							Name:  "spring.jms.servicebus.credential.managed-identity-enabled",
							Value: "false",
						},
						{
							Name:  "spring.jms.servicebus.credential.client-id",
							Value: "",
						},
					},
					SecretRefEnvironmentVariables: []scaffold.SecretRefEnvironmentVariable{
						{
							Name:      "spring.jms.servicebus.connection-string",
							SecretRef: "servicebus-connection-string",
						},
					},
					KeyVaultSecretDefinitions: []scaffold.KeyVaultSecretDefinition{
						{
							SecretName:  "servicebus-connection-string",
							KeyVaultUrl: "${keyVault.outputs.uri}secrets/SERVICEBUS-CONNECTION-STRING",
						},
					},
				}, nil
			default:
				return scaffold.EnvironmentVariableInformation{}, unsupportedAuthTypeWhenGetResourceConnectionEnvsError(resourceType, authType)
			}
		} else {
			// service bus, not jms
			switch authType {
			case internal.AuthTypeUserAssignedManagedIdentity:
				return scaffold.EnvironmentVariableInformation{
					StringEnvironmentVariables: []scaffold.StringEnvironmentVariable{
						{
							Name:  "spring.cloud.azure.servicebus.namespace",
							Value: "${serviceBusNamespace.outputs.name}",
						},
						{
							Name:  "spring.cloud.azure.servicebus.credential.managed-identity-enabled",
							Value: "true",
						},
						{
							Name:  "spring.cloud.azure.servicebus.credential.client-id",
							Value: "__PlaceHolderForServiceIdentityClientId",
						},
						{
							Name:  "spring.cloud.azure.servicebus.connection-string",
							Value: "",
						},
						// Not add it because of this: https://github.com/Azure/azure-sdk-for-java/issues/42880
						// Not add it even though the issue fixed, because customer may not use the new version
						//{
						//	Name:  "spring.cloud.azure.servicebus.connection-string",
						//	Value: "",
						//},
					},
				}, nil
			case internal.AuthTypeConnectionString:
				return scaffold.EnvironmentVariableInformation{
					StringEnvironmentVariables: []scaffold.StringEnvironmentVariable{
						{
							Name:  "spring.cloud.azure.servicebus.namespace",
							Value: "${serviceBusNamespace.outputs.name}",
						},
						{
							Name:  "spring.cloud.azure.servicebus.credential.managed-identity-enabled",
							Value: "false",
						},
						{
							Name:  "spring.cloud.azure.servicebus.credential.client-id",
							Value: "",
						},
					},
					SecretRefEnvironmentVariables: []scaffold.SecretRefEnvironmentVariable{
						{
							Name:      "spring.cloud.azure.servicebus.connection-string",
							SecretRef: "servicebus-connection-string",
						},
					},
					KeyVaultSecretDefinitions: []scaffold.KeyVaultSecretDefinition{
						{
							SecretName:  "servicebus-connection-string",
							KeyVaultUrl: "${keyVault.outputs.uri}secrets/SERVICEBUS-CONNECTION-STRING",
						},
					},
				}, nil
			default:
				return scaffold.EnvironmentVariableInformation{}, unsupportedAuthTypeWhenGetResourceConnectionEnvsError(resourceType, authType)
			}
		}
	case ResourceTypeMessagingKafka:
		// event hubs for kafka
		springBootVersionDecidedInformation := scaffold.EnvironmentVariableInformation{}
		if strings.HasPrefix(infraSpec.AzureEventHubs.SpringBootVersion, "2.") {
			springBootVersionDecidedInformation = scaffold.EnvironmentVariableInformation{
				StringEnvironmentVariables: []scaffold.StringEnvironmentVariable{
					{
						Name:  "spring.cloud.stream.binders.kafka.environment.spring.main.sources",
						Value: "com.azure.spring.cloud.autoconfigure.eventhubs.kafka.AzureEventHubsKafkaAutoConfiguration",
					},
				},
			}
		} else {
			springBootVersionDecidedInformation = scaffold.EnvironmentVariableInformation{
				StringEnvironmentVariables: []scaffold.StringEnvironmentVariable{
					{
						Name:  "spring.cloud.stream.binders.kafka.environment.spring.main.sources",
						Value: "com.azure.spring.cloud.autoconfigure.implementation.eventhubs.kafka.AzureEventHubsKafkaAutoConfiguration",
					},
				},
			}
		}
		commonInformation := scaffold.EnvironmentVariableInformation{}
		switch authType {
		case internal.AuthTypeUserAssignedManagedIdentity:
			commonInformation = scaffold.EnvironmentVariableInformation{
				StringEnvironmentVariables: []scaffold.StringEnvironmentVariable{
					{
						Name:  "spring.cloud.stream.kafka.binder.brokers",
						Value: "${eventHubNamespace.outputs.name}.servicebus.windows.net:9093",
					},
					{
						Name:  "spring.cloud.azure.eventhubs.credential.managed-identity-enabled",
						Value: "true",
					},
					{
						Name:  "spring.cloud.azure.eventhubs.credential.client-id",
						Value: "__PlaceHolderForServiceIdentityClientId",
					},
					// Not add it because of this: https://github.com/Azure/azure-sdk-for-java/issues/42880
					// Not add it even though the issue fixed, because customer may not use the new version
					//{
					//	Name:  "spring.cloud.azure.eventhubs.connection-string",
					//	Value: "",
					//},
				},
			}
		case internal.AuthTypeConnectionString:
			commonInformation = scaffold.EnvironmentVariableInformation{
				StringEnvironmentVariables: []scaffold.StringEnvironmentVariable{
					{
						Name:  "spring.cloud.stream.kafka.binder.brokers",
						Value: "${eventHubNamespace.outputs.name}.servicebus.windows.net:9093",
					},
					{
						Name:  "spring.cloud.azure.eventhubs.credential.managed-identity-enabled",
						Value: "false",
					},
					{
						Name:  "spring.cloud.azure.eventhubs.credential.client-id",
						Value: "",
					},
				},
				SecretRefEnvironmentVariables: []scaffold.SecretRefEnvironmentVariable{
					{
						Name:      "spring.cloud.azure.eventhubs.connection-string",
						SecretRef: "event-hubs-connection-string",
					},
				},
				KeyVaultSecretDefinitions: []scaffold.KeyVaultSecretDefinition{
					{
						SecretName:  "event-hubs-connection-string",
						KeyVaultUrl: "${keyVault.outputs.uri}secrets/EVENT-HUBS-CONNECTION-STRING",
					},
				},
			}
		default:
			return scaffold.EnvironmentVariableInformation{}, unsupportedAuthTypeWhenGetResourceConnectionEnvsError(resourceType, authType)
		}
		return mergeWithDuplicationCheck(springBootVersionDecidedInformation, commonInformation)
	case ResourceTypeMessagingEventHubs:
		// event hubs, not kafka
		switch authType {
		case internal.AuthTypeUserAssignedManagedIdentity:
			return scaffold.EnvironmentVariableInformation{
				StringEnvironmentVariables: []scaffold.StringEnvironmentVariable{
					{
						Name:  "spring.cloud.azure.eventhubs.namespace",
						Value: "${eventHubNamespace.outputs.name}",
					},
					{
						Name:  "spring.cloud.azure.eventhubs.credential.managed-identity-enabled",
						Value: "true",
					},
					{
						Name:  "spring.cloud.azure.eventhubs.credential.client-id",
						Value: "__PlaceHolderForServiceIdentityClientId",
					},
					// Not add it because of this: https://github.com/Azure/azure-sdk-for-java/issues/42880
					// Not add it even though the issue fixed, because customer may not use the new version
					//{
					//	Name:  "spring.cloud.azure.eventhubs.connection-string",
					//	Value: "",
					//},
				},
			}, nil
		case internal.AuthTypeConnectionString:
			return scaffold.EnvironmentVariableInformation{
				StringEnvironmentVariables: []scaffold.StringEnvironmentVariable{
					{
						Name:  "spring.cloud.azure.eventhubs.namespace",
						Value: "${eventHubNamespace.outputs.name}",
					},
					{
						Name:  "spring.cloud.azure.eventhubs.credential.managed-identity-enabled",
						Value: "false",
					},
					{
						Name:  "spring.cloud.azure.eventhubs.credential.client-id",
						Value: "",
					},
				},
				SecretRefEnvironmentVariables: []scaffold.SecretRefEnvironmentVariable{
					{
						Name:      "spring.cloud.azure.eventhubs.connection-string",
						SecretRef: "event-hubs-connection-string",
					},
				},
				KeyVaultSecretDefinitions: []scaffold.KeyVaultSecretDefinition{
					{
						SecretName:  "event-hubs-connection-string",
						KeyVaultUrl: "${keyVault.outputs.uri}secrets/EVENT-HUBS-CONNECTION-STRING",
					},
				},
			}, nil
		default:
			return scaffold.EnvironmentVariableInformation{}, unsupportedAuthTypeWhenGetResourceConnectionEnvsError(resourceType, authType)
		}
	case ResourceTypeStorage:
		switch authType {
		case internal.AuthTypeUserAssignedManagedIdentity:
			return scaffold.EnvironmentVariableInformation{
				StringEnvironmentVariables: []scaffold.StringEnvironmentVariable{
					{
						Name:  "spring.cloud.azure.eventhubs.processor.checkpoint-store.account-name",
						Value: "${storageAccountName}",
					},
					{
						Name:  "spring.cloud.azure.eventhubs.processor.checkpoint-store.credential.managed-identity-enabled",
						Value: "true",
					},
					{
						Name:  "spring.cloud.azure.eventhubs.processor.checkpoint-store.credential.client-id",
						Value: "__PlaceHolderForServiceIdentityClientId",
					},
					{
						Name:  "spring.cloud.azure.eventhubs.processor.checkpoint-store.connection-string",
						Value: "",
					},
				},
			}, nil
		case internal.AuthTypeConnectionString:
			return scaffold.EnvironmentVariableInformation{
				StringEnvironmentVariables: []scaffold.StringEnvironmentVariable{
					{
						Name:  "spring.cloud.azure.eventhubs.processor.checkpoint-store.account-name",
						Value: "${storageAccountName}",
					},
					{
						Name:  "spring.cloud.azure.eventhubs.processor.checkpoint-store.credential.managed-identity-enabled",
						Value: "false",
					},
					{
						Name:  "spring.cloud.azure.eventhubs.processor.checkpoint-store.credential.client-id",
						Value: "",
					},
				},
				SecretRefEnvironmentVariables: []scaffold.SecretRefEnvironmentVariable{
					{
						Name:      "spring.cloud.azure.eventhubs.processor.checkpoint-store.connection-string",
						SecretRef: "storage-account-connection-string",
					},
				},
				KeyVaultSecretDefinitions: []scaffold.KeyVaultSecretDefinition{
					{
						SecretName:  "storage-account-connection-string",
						KeyVaultUrl: "${keyVault.outputs.uri}secrets/STORAGE-ACCOUNT-CONNECTION-STRING",
					},
				},
			}, nil
		default:
			return scaffold.EnvironmentVariableInformation{}, unsupportedAuthTypeWhenGetResourceConnectionEnvsError(resourceType, authType)
		}
	case ResourceTypeOpenAiModel: // Keep this as code template
		switch authType {
		case internal.AuthTypeUserAssignedManagedIdentity:
			return scaffold.EnvironmentVariableInformation{
				StringEnvironmentVariables: []scaffold.StringEnvironmentVariable{
					{
						Name:  "AZURE_OPENAI_ENDPOINT",
						Value: "${account.outputs.endpoint}",
					},
				},
			}, nil
		default:
			return scaffold.EnvironmentVariableInformation{}, unsupportedAuthTypeWhenGetResourceConnectionEnvsError(resourceType, authType)
		}
		//case OtherType: // Keep this as code template
		//	switch authType {
		//	default:
		//		return scaffold.EnvironmentVariableInformation{}, unsupportedAuthTypeWhenGetResourceConnectionEnvsError(resourceType, authType)
		//	}
	case ResourceTypeHostContainerApp: // todo improve this and delete Frontend and Backend in scaffold.ServiceSpec
		switch authType {
		case internal.AuthTypeUserAssignedManagedIdentity:
			return scaffold.EnvironmentVariableInformation{}, nil
		default:
			return scaffold.EnvironmentVariableInformation{}, unsupportedAuthTypeWhenGetResourceConnectionEnvsError(resourceType, authType)
		}
	default:
		return scaffold.EnvironmentVariableInformation{}, unsupportedResourceTypeWhenGetResourceConnectionEnvsError(resourceType)
	}
}

func unsupportedResourceTypeWhenGetResourceConnectionEnvsError(resourceType ResourceType) error {
	return fmt.Errorf("unsupported resource type when getResourceConnectionEnvs, resourceType = %s",
		resourceType)
}

func unsupportedAuthTypeWhenGetResourceConnectionEnvsError(resourceType ResourceType, authType internal.AuthType) error {
	return fmt.Errorf("unsupported auth type when getResourceConnectionEnvs, "+
		"resourceType = %s, authType = %s", resourceType, authType)
}

func getAllResourceConnectionEnvForPrint(usedResource *ResourceConfig,
	infraSpec *scaffold.InfraSpec) ([]scaffold.ResourceConnectionEnv, error) {
	envs, err := getResourceConnectionEnvs(usedResource, infraSpec)
	if err != nil {
		return []scaffold.ResourceConnectionEnv{}, err
	}
	resourceType := usedResource.Type
	authType, err := getAuthType(infraSpec, usedResource.Type)
	if err != nil {
		return []scaffold.ResourceConnectionEnv{}, err
	}
	additional, err := getResourceConnectionEnvNamesCreatedByServiceConnector(resourceType, authType)
	if err != nil {
		return []scaffold.ResourceConnectionEnv{}, err
	}
	result, err := mergeResourceConnectionEnvWithDuplicationCheck(envs, additional)
	if err != nil {
		return []scaffold.ResourceConnectionEnv{}, err
	}
	return result, nil
}

// Return environment variables added by service connector, they do not need to add to scaffold.ServiceSpec
// todo: Now only support springBoot application type. Need to support other types
func getResourceConnectionEnvNamesCreatedByServiceConnector(resourceType ResourceType,
	authType internal.AuthType) ([]scaffold.ResourceConnectionEnv, error) {
	switch resourceType {
	case ResourceTypeDbPostgres:
		switch authType {
		case internal.AuthTypePassword:
			return []scaffold.ResourceConnectionEnv{}, nil
		case internal.AuthTypeUserAssignedManagedIdentity:
			return []scaffold.ResourceConnectionEnv{
				{
					Name: "spring.datasource.url",
				},
				{
					Name: "spring.datasource.username",
				},
				{
					Name: "spring.datasource.azure.passwordless-enabled",
				},
			}, nil
		default:
			// return error to make sure every case has been considered.
			return []scaffold.ResourceConnectionEnv{}, unsupportedAuthTypeWhenGetResourceConnectionEnvsError(resourceType, authType)
		}
	case ResourceTypeDbMySQL:
		switch authType {
		case internal.AuthTypePassword:
			return []scaffold.ResourceConnectionEnv{}, nil
		case internal.AuthTypeUserAssignedManagedIdentity:
			return []scaffold.ResourceConnectionEnv{
				{
					Name: "spring.datasource.url",
				},
				{
					Name: "spring.datasource.username",
				},
				{
					Name: "spring.datasource.azure.passwordless-enabled",
				},
			}, nil
		default:
			// return error to make sure every case has been considered.
			return []scaffold.ResourceConnectionEnv{}, unsupportedAuthTypeWhenGetResourceConnectionEnvsError(resourceType, authType)
		}
	case ResourceTypeDbRedis:
		switch authType {
		case internal.AuthTypePassword:
			return []scaffold.ResourceConnectionEnv{}, nil
		default:
			// return error to make sure every case has been considered.
			return []scaffold.ResourceConnectionEnv{}, unsupportedAuthTypeWhenGetResourceConnectionEnvsError(resourceType, authType)
		}
	case ResourceTypeDbMongo,
		ResourceTypeDbCosmos,
		ResourceTypeOpenAiModel,
		ResourceTypeHostContainerApp:
		switch authType {
		case internal.AuthTypeUserAssignedManagedIdentity:
			return []scaffold.ResourceConnectionEnv{}, nil
		default:
			// return error to make sure every case has been considered.
			return []scaffold.ResourceConnectionEnv{}, unsupportedAuthTypeWhenGetResourceConnectionEnvsError(resourceType, authType)
		}
	case
		ResourceTypeMessagingServiceBus,
		ResourceTypeMessagingEventHubs,
		ResourceTypeMessagingKafka,
		ResourceTypeStorage:
		switch authType {
		case internal.AuthTypeUserAssignedManagedIdentity, internal.AuthTypeConnectionString:
			return []scaffold.ResourceConnectionEnv{}, nil
		default:
			// return error to make sure every case has been considered.
			return []scaffold.ResourceConnectionEnv{}, unsupportedAuthTypeWhenGetResourceConnectionEnvsError(resourceType, authType)
		}
	default:
		// return error to make sure every case has been considered.
		return []scaffold.ResourceConnectionEnv{}, unsupportedResourceTypeWhenGetResourceConnectionEnvsError(resourceType)
	}
}

func mergeResourceConnectionEnvWithDuplicationCheck(a []scaffold.ResourceConnectionEnv,
	b []scaffold.ResourceConnectionEnv) ([]scaffold.ResourceConnectionEnv, error) {
	result := append(a, b...)
	seen := make(map[string]scaffold.ResourceConnectionEnv)
	for _, value := range result {
		if existingValue, exist := seen[value.Name]; exist {
			if value != existingValue {
				return []scaffold.ResourceConnectionEnv{}, duplicatedResourceConnectionEnvError(existingValue, value)
			}
		} else {
			seen[value.Name] = existingValue
		}
	}
	return result, nil
}

func mergeWithDuplicationCheck(a scaffold.EnvironmentVariableInformation,
	b scaffold.EnvironmentVariableInformation) (scaffold.EnvironmentVariableInformation, error) {
	result := scaffold.EnvironmentVariableInformation{
		StringEnvironmentVariables:    append(a.StringEnvironmentVariables, b.StringEnvironmentVariables...),
		SecretRefEnvironmentVariables: append(a.SecretRefEnvironmentVariables, b.SecretRefEnvironmentVariables...),
		ValueSecretDefinitions:        append(a.ValueSecretDefinitions, b.ValueSecretDefinitions...),
		KeyVaultSecretDefinitions:     append(a.KeyVaultSecretDefinitions, b.KeyVaultSecretDefinitions...),
	}
	seen := make(map[string]string)
	for _, v := range result.StringEnvironmentVariables {
		if existingValue, exist := seen[v.Name]; exist {
			if v.Value != existingValue {
				return scaffold.EnvironmentVariableInformation{},
					duplicatedEnvironmentError(v.Name, v.Value, existingValue)
			}
		} else {
			seen[v.Name] = existingValue
		}
	}
	for _, v := range result.SecretRefEnvironmentVariables {
		if existingRef, exist := seen[v.Name]; exist {
			if v.SecretRef != existingRef {
				return scaffold.EnvironmentVariableInformation{},
					duplicatedEnvironmentError(v.Name, v.SecretRef, existingRef)
			}
		} else {
			seen[v.Name] = existingRef
		}
	}
	seen = make(map[string]string)
	for _, v := range result.ValueSecretDefinitions {
		if existingRef, exist := seen[v.SecretName]; exist {
			if v.SecretValue != existingRef {
				return scaffold.EnvironmentVariableInformation{},
					duplicatedSecretDefinitionError(v.SecretName, v.SecretValue, existingRef)
			}
		} else {
			seen[v.SecretName] = existingRef
		}
	}
	for _, v := range result.KeyVaultSecretDefinitions {
		if existingRef, exist := seen[v.SecretName]; exist {
			if v.SecretName != existingRef {
				return scaffold.EnvironmentVariableInformation{},
					duplicatedSecretDefinitionError(v.SecretName, v.KeyVaultUrl, existingRef)
			}
		} else {
			seen[v.SecretName] = existingRef
		}
	}
	return result, nil
}

func addNewEnvironmentVariable(serviceSpec *scaffold.ServiceSpec, name string, value string) error {
	merged, err := mergeWithDuplicationCheck(serviceSpec.EnvironmentVariableInformation,
		scaffold.EnvironmentVariableInformation{
			StringEnvironmentVariables: []scaffold.StringEnvironmentVariable{
				{
					Name:  name,
					Value: value,
				},
			},
		},
	)
	if err != nil {
		return err
	}
	serviceSpec.EnvironmentVariableInformation = merged
	return nil
}

func duplicatedSecretDefinitionError(name string, value1 string, value2 string) error {
	return duplicatedError("secret definition", name, value1, value2)
}

func duplicatedResourceConnectionEnvError(existingValue scaffold.ResourceConnectionEnv,
	value scaffold.ResourceConnectionEnv) error {
	return fmt.Errorf(
		"duplicated ResourceConnectionEnv. existingValue = %s, value = %s",
		existingValue.ToString(), value.ToString())
}

func duplicatedEnvironmentError(name string, value1 string, value2 string) error {
	return duplicatedError("environment variable", name, value1, value2)
}

func duplicatedError(description string, name string, value1 string, value2 string) error {
	return fmt.Errorf(
		"duplicated %s. name = %s, value1 = %s, value2 = %s", description, name, value1, value2)
}
