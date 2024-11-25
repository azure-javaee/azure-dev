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
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "POSTGRES_USERNAME",
					ResourceType:              scaffold.ResourceTypeDbPostgres,
					ResourceInfoType:          scaffold.ResourceInfoTypeUsername,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "POSTGRES_PASSWORD",
					ResourceType:              scaffold.ResourceTypeDbPostgres,
					ResourceInfoType:          scaffold.ResourceInfoTypePassword,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "POSTGRES_HOST",
					ResourceType:              scaffold.ResourceTypeDbPostgres,
					ResourceInfoType:          scaffold.ResourceInfoTypeHost,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "POSTGRES_DATABASE",
					ResourceType:              scaffold.ResourceTypeDbPostgres,
					ResourceInfoType:          scaffold.ResourceInfoTypeDatabaseName,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "POSTGRES_PORT",
					ResourceType:              scaffold.ResourceTypeDbPostgres,
					ResourceInfoType:          scaffold.ResourceInfoTypePort,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "POSTGRES_URL",
					ResourceType:              scaffold.ResourceTypeDbPostgres,
					ResourceInfoType:          scaffold.ResourceInfoTypeUrl,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "spring.datasource.url",
					ResourceType:              scaffold.ResourceTypeDbPostgres,
					ResourceInfoType:          scaffold.ResourceInfoTypeJdbcUrl,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "spring.datasource.username",
					ResourceType:              scaffold.ResourceTypeDbPostgres,
					ResourceInfoType:          scaffold.ResourceInfoTypeUsername,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "spring.datasource.password",
					ResourceType:              scaffold.ResourceTypeDbPostgres,
					ResourceInfoType:          scaffold.ResourceInfoTypePassword,
				},
			}, nil
		case internal.AuthTypeUserAssignedManagedIdentity:
			return []scaffold.ResourceConnectionEnv{
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "POSTGRES_USERNAME",
					ResourceType:              scaffold.ResourceTypeDbPostgres,
					ResourceInfoType:          scaffold.ResourceInfoTypeUsername,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "POSTGRES_HOST",
					ResourceType:              scaffold.ResourceTypeDbPostgres,
					ResourceInfoType:          scaffold.ResourceInfoTypeHost,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "POSTGRES_DATABASE",
					ResourceType:              scaffold.ResourceTypeDbPostgres,
					ResourceInfoType:          scaffold.ResourceInfoTypeDatabaseName,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "POSTGRES_PORT",
					ResourceType:              scaffold.ResourceTypeDbPostgres,
					ResourceInfoType:          scaffold.ResourceInfoTypePort,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeServiceConnectorCreated,
					Name:                      "spring.datasource.url",
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeServiceConnectorCreated,
					Name:                      "spring.datasource.username",
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeServiceConnectorCreated,
					Name:                      "spring.datasource.azure.passwordless-enabled",
				},
			}, nil
		default:
			return []scaffold.ResourceConnectionEnv{}, unsupportedAuthTypeError(resourceType, authType)
		}
	case ResourceTypeDbMySQL:
		switch authType {
		case internal.AuthTypePassword:
			return []scaffold.ResourceConnectionEnv{
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "MYSQL_USERNAME",
					ResourceType:              scaffold.ResourceTypeDbMySQL,
					ResourceInfoType:          scaffold.ResourceInfoTypeUsername,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "MYSQL_PASSWORD",
					ResourceType:              scaffold.ResourceTypeDbMySQL,
					ResourceInfoType:          scaffold.ResourceInfoTypePassword,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "MYSQL_HOST",
					ResourceType:              scaffold.ResourceTypeDbMySQL,
					ResourceInfoType:          scaffold.ResourceInfoTypeHost,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "MYSQL_DATABASE",
					ResourceType:              scaffold.ResourceTypeDbMySQL,
					ResourceInfoType:          scaffold.ResourceInfoTypeDatabaseName,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "MYSQL_PORT",
					ResourceType:              scaffold.ResourceTypeDbMySQL,
					ResourceInfoType:          scaffold.ResourceInfoTypePort,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "MYSQL_URL",
					ResourceType:              scaffold.ResourceTypeDbMySQL,
					ResourceInfoType:          scaffold.ResourceInfoTypeUrl,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "spring.datasource.url",
					ResourceType:              scaffold.ResourceTypeDbMySQL,
					ResourceInfoType:          scaffold.ResourceInfoTypeJdbcUrl,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "spring.datasource.username",
					ResourceType:              scaffold.ResourceTypeDbMySQL,
					ResourceInfoType:          scaffold.ResourceInfoTypeUsername,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "spring.datasource.password",
					ResourceType:              scaffold.ResourceTypeDbMySQL,
					ResourceInfoType:          scaffold.ResourceInfoTypePassword,
				},
			}, nil
		case internal.AuthTypeUserAssignedManagedIdentity:
			return []scaffold.ResourceConnectionEnv{
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "MYSQL_USERNAME",
					ResourceType:              scaffold.ResourceTypeDbMySQL,
					ResourceInfoType:          scaffold.ResourceInfoTypeUsername,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "MYSQL_HOST",
					ResourceType:              scaffold.ResourceTypeDbMySQL,
					ResourceInfoType:          scaffold.ResourceInfoTypeHost,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "MYSQL_PORT",
					ResourceType:              scaffold.ResourceTypeDbMySQL,
					ResourceInfoType:          scaffold.ResourceInfoTypePort,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "MYSQL_DATABASE",
					ResourceType:              scaffold.ResourceTypeDbMySQL,
					ResourceInfoType:          scaffold.ResourceInfoTypeDatabaseName,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeServiceConnectorCreated,
					Name:                      "spring.datasource.url",
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeServiceConnectorCreated,
					Name:                      "spring.datasource.username",
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeServiceConnectorCreated,
					Name:                      "spring.datasource.azure.passwordless-enabled",
				},
			}, nil
		default:
			return []scaffold.ResourceConnectionEnv{}, unsupportedAuthTypeError(resourceType, authType)
		}
	case ResourceTypeDbRedis:
		switch authType {
		case internal.AuthTypePassword:
			return []scaffold.ResourceConnectionEnv{
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "REDIS_HOST",
					ResourceType:              scaffold.ResourceTypeDbRedis,
					ResourceInfoType:          scaffold.ResourceInfoTypeHost,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "REDIS_PORT",
					ResourceType:              scaffold.ResourceTypeDbRedis,
					ResourceInfoType:          scaffold.ResourceInfoTypePort,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "REDIS_ENDPOINT",
					ResourceType:              scaffold.ResourceTypeDbRedis,
					ResourceInfoType:          scaffold.ResourceInfoTypeEndpoint,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "REDIS_URL",
					ResourceType:              scaffold.ResourceTypeDbRedis,
					ResourceInfoType:          scaffold.ResourceInfoTypeUrl,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "REDIS_PASSWORD",
					ResourceType:              scaffold.ResourceTypeDbRedis,
					ResourceInfoType:          scaffold.ResourceInfoTypePassword,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "spring.data.redis.url",
					ResourceType:              scaffold.ResourceTypeDbRedis,
					ResourceInfoType:          scaffold.ResourceInfoTypeUrl,
				},
			}, nil
		default:
			return []scaffold.ResourceConnectionEnv{}, unsupportedAuthTypeError(resourceType, authType)
		}
	case ResourceTypeDbMongo:
		switch authType {
		case internal.AuthTypeUserAssignedManagedIdentity:
			return []scaffold.ResourceConnectionEnv{
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "MONGODB_URL",
					ResourceType:              scaffold.ResourceTypeDbMongo,
					ResourceInfoType:          scaffold.ResourceInfoTypeUrl,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "spring.data.mongodb.uri",
					ResourceType:              scaffold.ResourceTypeDbMongo,
					ResourceInfoType:          scaffold.ResourceInfoTypeUrl,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "spring.data.mongodb.database",
					ResourceType:              scaffold.ResourceTypeDbMongo,
					ResourceInfoType:          scaffold.ResourceInfoTypeDatabaseName,
				},
			}, nil
		default:
			return []scaffold.ResourceConnectionEnv{}, unsupportedAuthTypeError(resourceType, authType)
		}
	case ResourceTypeDbCosmos:
		switch authType {
		case internal.AuthTypeUserAssignedManagedIdentity:
			return []scaffold.ResourceConnectionEnv{
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "spring.cloud.azure.cosmos.endpoint",
					ResourceType:              scaffold.ResourceTypeDbCosmos,
					ResourceInfoType:          scaffold.ResourceInfoTypeEndpoint,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "spring.cloud.azure.cosmos.database",
					ResourceType:              scaffold.ResourceTypeDbCosmos,
					ResourceInfoType:          scaffold.ResourceInfoTypeDatabaseName,
				},
			}, nil
		default:
			return []scaffold.ResourceConnectionEnv{}, unsupportedAuthTypeError(resourceType, authType)
		}
	case ResourceTypeMessagingServiceBus:
		if infraSpec.AzureServiceBus.IsJms {
			switch authType {
			case internal.AuthTypeUserAssignedManagedIdentity:
				return []scaffold.ResourceConnectionEnv{
					{
						ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
						Name:                      "spring.jms.servicebus.pricing-tier",
						PlainTextValue:            "premium",
					},
					{
						ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
						Name:                      "spring.jms.servicebus.passwordless-enabled",
						PlainTextValue:            "true",
					},
					{
						ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
						Name:                      "spring.jms.servicebus.credential.managed-identity-enabled",
						PlainTextValue:            "true",
					},
					{
						ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
						Name:                      "spring.jms.servicebus.credential.client-id",
						PlainTextValue:            scaffold.PlaceHolderForServiceIdentityClientId(),
					},
					{
						ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
						Name:                      "spring.jms.servicebus.namespace",
						ResourceType:              scaffold.ResourceTypeMessagingServiceBus,
						ResourceInfoType:          scaffold.ResourceInfoTypeNamespace,
					},
					{
						ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
						Name:                      "spring.jms.servicebus.connection-string",
						PlainTextValue:            "",
					},
				}, nil
			case internal.AuthTypeConnectionString:
				return []scaffold.ResourceConnectionEnv{
					{
						ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
						Name:                      "spring.jms.servicebus.pricing-tier",
						PlainTextValue:            "premium",
					},
					{
						ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
						Name:                      "spring.jms.servicebus.connection-string",
						ResourceType:              scaffold.ResourceTypeMessagingServiceBus,
						ResourceInfoType:          scaffold.ResourceInfoTypeConnectionString,
					},
					{
						ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
						Name:                      "spring.jms.servicebus.passwordless-enabled",
						PlainTextValue:            "false",
					},
					{
						ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
						Name:                      "spring.jms.servicebus.credential.managed-identity-enabled",
						PlainTextValue:            "false",
					},
					{
						ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
						Name:                      "spring.jms.servicebus.credential.client-id",
						PlainTextValue:            "",
					},
					{
						ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
						Name:                      "spring.jms.servicebus.namespace",
						PlainTextValue:            "",
					},
				}, nil
			default:
				return []scaffold.ResourceConnectionEnv{}, unsupportedResourceTypeError(resourceType)
			}
		} else {
			// service bus, not jms
			switch authType {
			case internal.AuthTypeUserAssignedManagedIdentity:
				return []scaffold.ResourceConnectionEnv{
					// Not add this: spring.cloud.azure.servicebus.connection-string = ""
					// because of this: https://github.com/Azure/azure-sdk-for-java/issues/42880
					{
						ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
						Name:                      "spring.cloud.azure.servicebus.credential.managed-identity-enabled",
						PlainTextValue:            "true",
					},
					{
						ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
						Name:                      "spring.cloud.azure.servicebus.credential.client-id",
						PlainTextValue:            scaffold.PlaceHolderForServiceIdentityClientId(),
					},
					{
						ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
						Name:                      "spring.cloud.azure.servicebus.namespace",
						ResourceType:              scaffold.ResourceTypeMessagingServiceBus,
						ResourceInfoType:          scaffold.ResourceInfoTypeNamespace,
					},
				}, nil
			case internal.AuthTypeConnectionString:
				return []scaffold.ResourceConnectionEnv{
					{
						ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
						Name:                      "spring.cloud.azure.servicebus.namespace",
						ResourceType:              scaffold.ResourceTypeMessagingServiceBus,
						ResourceInfoType:          scaffold.ResourceInfoTypeNamespace,
					},
					{
						ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
						Name:                      "spring.cloud.azure.servicebus.connection-string",
						ResourceType:              scaffold.ResourceTypeMessagingServiceBus,
						ResourceInfoType:          scaffold.ResourceInfoTypeConnectionString,
					},
					{
						ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
						Name:                      "spring.cloud.azure.servicebus.credential.managed-identity-enabled",
						PlainTextValue:            "false",
					},
					{
						ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
						Name:                      "spring.cloud.azure.servicebus.credential.client-id",
						PlainTextValue:            "",
					},
				}, nil
			default:
				return []scaffold.ResourceConnectionEnv{}, unsupportedResourceTypeError(resourceType)
			}
		}
	case ResourceTypeMessagingKafka:
		// event hubs for kafka
		var springBootVersionDecidedInformation []scaffold.ResourceConnectionEnv
		if strings.HasPrefix(infraSpec.AzureEventHubs.SpringBootVersion, "2.") {
			springBootVersionDecidedInformation = []scaffold.ResourceConnectionEnv{
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
					Name:                      "spring.cloud.stream.binders.kafka.environment.spring.main.sources",
					PlainTextValue:            "com.azure.spring.cloud.autoconfigure.eventhubs.kafka.AzureEventHubsKafkaAutoConfiguration",
				},
			}
		} else {
			springBootVersionDecidedInformation = []scaffold.ResourceConnectionEnv{
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
					Name:                      "spring.cloud.stream.binders.kafka.environment.spring.main.sources",
					PlainTextValue:            "com.azure.spring.cloud.autoconfigure.implementation.eventhubs.kafka.AzureEventHubsKafkaAutoConfiguration",
				},
			}
		}
		var commonInformation []scaffold.ResourceConnectionEnv
		switch authType {
		case internal.AuthTypeUserAssignedManagedIdentity:
			commonInformation = []scaffold.ResourceConnectionEnv{
				// Not add this: spring.cloud.azure.eventhubs.connection-string = ""
				// because of this: https://github.com/Azure/azure-sdk-for-java/issues/42880
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
					Name:                      "spring.cloud.azure.eventhubs.credential.managed-identity-enabled",
					PlainTextValue:            "true",
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
					Name:                      "spring.cloud.azure.eventhubs.credential.client-id",
					PlainTextValue:            scaffold.PlaceHolderForServiceIdentityClientId(),
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "spring.cloud.stream.kafka.binder.brokers",
					ResourceType:              scaffold.ResourceTypeMessagingKafka,
					ResourceInfoType:          scaffold.ResourceInfoTypeEndpoint,
				},
			}
		case internal.AuthTypeConnectionString:
			commonInformation = []scaffold.ResourceConnectionEnv{
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "spring.cloud.stream.kafka.binder.brokers",
					ResourceType:              scaffold.ResourceTypeMessagingKafka,
					ResourceInfoType:          scaffold.ResourceInfoTypeEndpoint,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "spring.cloud.azure.eventhubs.connection-string",
					ResourceType:              scaffold.ResourceTypeMessagingKafka,
					ResourceInfoType:          scaffold.ResourceInfoTypeConnectionString,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
					Name:                      "spring.cloud.azure.eventhubs.credential.managed-identity-enabled",
					PlainTextValue:            "false",
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
					Name:                      "spring.cloud.azure.eventhubs.credential.client-id",
					PlainTextValue:            "",
				},
			}
		default:
			return []scaffold.ResourceConnectionEnv{}, unsupportedAuthTypeError(resourceType, authType)
		}
		return mergeResourceConnectionEnvWithDuplicationCheck(springBootVersionDecidedInformation, commonInformation)
	case ResourceTypeMessagingEventHubs:
		switch authType {
		case internal.AuthTypeUserAssignedManagedIdentity:
			return []scaffold.ResourceConnectionEnv{
				// Not add this: spring.cloud.azure.eventhubs.connection-string = ""
				// because of this: https://github.com/Azure/azure-sdk-for-java/issues/42880
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
					Name:                      "spring.cloud.azure.eventhubs.credential.managed-identity-enabled",
					PlainTextValue:            "true",
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
					Name:                      "spring.cloud.azure.eventhubs.credential.client-id",
					PlainTextValue:            scaffold.PlaceHolderForServiceIdentityClientId(),
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "spring.cloud.azure.eventhubs.namespace",
					ResourceType:              scaffold.ResourceTypeMessagingEventHubs,
					ResourceInfoType:          scaffold.ResourceInfoTypeNamespace,
				},
			}, nil
		case internal.AuthTypeConnectionString:
			return []scaffold.ResourceConnectionEnv{
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "spring.cloud.azure.eventhubs.namespace",
					ResourceType:              scaffold.ResourceTypeMessagingEventHubs,
					ResourceInfoType:          scaffold.ResourceInfoTypeNamespace,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "spring.cloud.azure.eventhubs.connection-string",
					ResourceType:              scaffold.ResourceTypeMessagingEventHubs,
					ResourceInfoType:          scaffold.ResourceInfoTypeConnectionString,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
					Name:                      "spring.cloud.azure.eventhubs.credential.managed-identity-enabled",
					PlainTextValue:            "false",
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
					Name:                      "spring.cloud.azure.eventhubs.credential.client-id",
					PlainTextValue:            "",
				},
			}, nil
		default:
			return []scaffold.ResourceConnectionEnv{}, unsupportedResourceTypeError(resourceType)
		}
	case ResourceTypeStorage:
		switch authType {
		case internal.AuthTypeUserAssignedManagedIdentity:
			return []scaffold.ResourceConnectionEnv{
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "spring.cloud.azure.eventhubs.processor.checkpoint-store.account-name",
					ResourceType:              scaffold.ResourceTypeStorage,
					ResourceInfoType:          scaffold.ResourceInfoTypeAccountName,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
					Name:                      "spring.cloud.azure.eventhubs.processor.checkpoint-store.credential.managed-identity-enabled",
					PlainTextValue:            "true",
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
					Name:                      "spring.cloud.azure.eventhubs.processor.checkpoint-store.credential.client-id",
					PlainTextValue:            scaffold.PlaceHolderForServiceIdentityClientId(),
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
					Name:                      "spring.cloud.azure.eventhubs.processor.checkpoint-store.connection-string",
					PlainTextValue:            "",
				},
			}, nil
		case internal.AuthTypeConnectionString:
			return []scaffold.ResourceConnectionEnv{
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "spring.cloud.azure.eventhubs.processor.checkpoint-store.account-name",
					ResourceType:              scaffold.ResourceTypeStorage,
					ResourceInfoType:          scaffold.ResourceInfoTypeAccountName,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "spring.cloud.azure.eventhubs.processor.checkpoint-store.connection-string",
					ResourceType:              scaffold.ResourceTypeStorage,
					ResourceInfoType:          scaffold.ResourceInfoTypeConnectionString,
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
					Name:                      "spring.cloud.azure.eventhubs.processor.checkpoint-store.credential.managed-identity-enabled",
					PlainTextValue:            "false",
				},
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypePlainText,
					Name:                      "spring.cloud.azure.eventhubs.processor.checkpoint-store.credential.client-id",
					PlainTextValue:            "",
				},
			}, nil
		default:
			return []scaffold.ResourceConnectionEnv{}, unsupportedResourceTypeError(resourceType)
		}
	case ResourceTypeOpenAiModel:
		switch authType {
		case internal.AuthTypeUserAssignedManagedIdentity:
			return []scaffold.ResourceConnectionEnv{
				{
					ResourceConnectionEnvType: scaffold.ResourceConnectionEnvTypeResourceSpecific,
					Name:                      "AZURE_OPENAI_ENDPOINT",
					ResourceType:              scaffold.ResourceTypeOpenAiModel,
					ResourceInfoType:          scaffold.ResourceInfoTypeEndpoint,
				},
			}, nil
		default:
			return []scaffold.ResourceConnectionEnv{}, unsupportedResourceTypeError(resourceType)
		}
	case ResourceTypeHostContainerApp: // todo improve this and delete Frontend and Backend in scaffold.ServiceSpec
		switch authType {
		case internal.AuthTypeUserAssignedManagedIdentity:
			return []scaffold.ResourceConnectionEnv{}, nil
		default:
			return []scaffold.ResourceConnectionEnv{}, unsupportedAuthTypeError(resourceType, authType)
		}
	default:
		return []scaffold.ResourceConnectionEnv{}, unsupportedResourceTypeError(resourceType)
	}
}

func unsupportedResourceTypeError(resourceType ResourceType) error {
	return fmt.Errorf("unsupported resource type, resourceType = %s", resourceType)
}

func unsupportedAuthTypeError(resourceType ResourceType, authType internal.AuthType) error {
	return fmt.Errorf("unsupported auth type, resourceType = %s, authType = %s", resourceType, authType)
}

func mergeResourceConnectionEnvWithDuplicationCheck(a []scaffold.ResourceConnectionEnv,
	b []scaffold.ResourceConnectionEnv) ([]scaffold.ResourceConnectionEnv, error) {
	result := append(a, b...)
	seenName := make(map[string]scaffold.ResourceConnectionEnv)
	for _, value := range result {
		if existingValue, exist := seenName[value.Name]; exist {
			if value != existingValue {
				return []scaffold.ResourceConnectionEnv{}, duplicatedResourceConnectionEnvError(existingValue, value)
			}
		} else {
			seenName[value.Name] = existingValue
		}
	}
	return result, nil
}

func mergeEnvWithDuplicationCheck(a []scaffold.Env, b []scaffold.Env) ([]scaffold.Env, error) {
	result := append(a, b...)
	seen := make(map[string]string)
	for _, v := range result {
		if existingValue, exist := seen[v.Name]; exist {
			if v.Value != existingValue {
				return []scaffold.Env{}, duplicatedEnvironmentError(v.Name, v.Value, existingValue)
			}
		} else {
			seen[v.Name] = existingValue
		}
	}
	return result, nil
}

func addNewEnvironmentVariable(serviceSpec *scaffold.ServiceSpec, name string, value string) error {
	merged, err := mergeEnvWithDuplicationCheck(serviceSpec.Envs,
		[]scaffold.Env{
			{
				Name:  name,
				Value: value,
			},
		},
	)
	if err != nil {
		return err
	}
	serviceSpec.Envs = merged
	return nil
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
