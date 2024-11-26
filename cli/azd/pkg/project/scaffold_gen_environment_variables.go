package project

import (
	"fmt"
	"github.com/azure/azure-dev/cli/azd/internal"
	"github.com/azure/azure-dev/cli/azd/internal/scaffold"
	"strings"
)

func getResourceConnectionEnvs(usedResource *ResourceConfig,
	infraSpec *scaffold.InfraSpec) ([]scaffold.Env, error) {
	resourceType := usedResource.Type
	authType, err := getAuthType(infraSpec, usedResource.Type)
	if err != nil {
		return []scaffold.Env{}, err
	}
	switch resourceType {
	case ResourceTypeDbPostgres:
		switch authType {
		case internal.AuthTypePassword:
			return []scaffold.Env{
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "POSTGRES_USERNAME",
					ResourceType:     scaffold.ResourceTypeDbPostgres,
					ResourceInfoType: scaffold.ResourceInfoTypeUsername,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "POSTGRES_PASSWORD",
					ResourceType:     scaffold.ResourceTypeDbPostgres,
					ResourceInfoType: scaffold.ResourceInfoTypePassword,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "POSTGRES_HOST",
					ResourceType:     scaffold.ResourceTypeDbPostgres,
					ResourceInfoType: scaffold.ResourceInfoTypeHost,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "POSTGRES_DATABASE",
					ResourceType:     scaffold.ResourceTypeDbPostgres,
					ResourceInfoType: scaffold.ResourceInfoTypeDatabaseName,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "POSTGRES_PORT",
					ResourceType:     scaffold.ResourceTypeDbPostgres,
					ResourceInfoType: scaffold.ResourceInfoTypePort,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "POSTGRES_URL",
					ResourceType:     scaffold.ResourceTypeDbPostgres,
					ResourceInfoType: scaffold.ResourceInfoTypeUrl,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "spring.datasource.url",
					ResourceType:     scaffold.ResourceTypeDbPostgres,
					ResourceInfoType: scaffold.ResourceInfoTypeJdbcUrl,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "spring.datasource.username",
					ResourceType:     scaffold.ResourceTypeDbPostgres,
					ResourceInfoType: scaffold.ResourceInfoTypeUsername,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "spring.datasource.password",
					ResourceType:     scaffold.ResourceTypeDbPostgres,
					ResourceInfoType: scaffold.ResourceInfoTypePassword,
				},
			}, nil
		case internal.AuthTypeUserAssignedManagedIdentity:
			return []scaffold.Env{
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "POSTGRES_USERNAME",
					ResourceType:     scaffold.ResourceTypeDbPostgres,
					ResourceInfoType: scaffold.ResourceInfoTypeUsername,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "POSTGRES_HOST",
					ResourceType:     scaffold.ResourceTypeDbPostgres,
					ResourceInfoType: scaffold.ResourceInfoTypeHost,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "POSTGRES_DATABASE",
					ResourceType:     scaffold.ResourceTypeDbPostgres,
					ResourceInfoType: scaffold.ResourceInfoTypeDatabaseName,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "POSTGRES_PORT",
					ResourceType:     scaffold.ResourceTypeDbPostgres,
					ResourceInfoType: scaffold.ResourceInfoTypePort,
				},
				{
					EnvType: scaffold.EnvTypeResourceConnectionServiceConnectorCreated,
					Name:    "spring.datasource.url",
				},
				{
					EnvType: scaffold.EnvTypeResourceConnectionServiceConnectorCreated,
					Name:    "spring.datasource.username",
				},
				{
					EnvType: scaffold.EnvTypeResourceConnectionServiceConnectorCreated,
					Name:    "spring.datasource.azure.passwordless-enabled",
				},
			}, nil
		default:
			return []scaffold.Env{}, unsupportedAuthTypeError(resourceType, authType)
		}
	case ResourceTypeDbMySQL:
		switch authType {
		case internal.AuthTypePassword:
			return []scaffold.Env{
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "MYSQL_USERNAME",
					ResourceType:     scaffold.ResourceTypeDbMySQL,
					ResourceInfoType: scaffold.ResourceInfoTypeUsername,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "MYSQL_PASSWORD",
					ResourceType:     scaffold.ResourceTypeDbMySQL,
					ResourceInfoType: scaffold.ResourceInfoTypePassword,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "MYSQL_HOST",
					ResourceType:     scaffold.ResourceTypeDbMySQL,
					ResourceInfoType: scaffold.ResourceInfoTypeHost,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "MYSQL_DATABASE",
					ResourceType:     scaffold.ResourceTypeDbMySQL,
					ResourceInfoType: scaffold.ResourceInfoTypeDatabaseName,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "MYSQL_PORT",
					ResourceType:     scaffold.ResourceTypeDbMySQL,
					ResourceInfoType: scaffold.ResourceInfoTypePort,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "MYSQL_URL",
					ResourceType:     scaffold.ResourceTypeDbMySQL,
					ResourceInfoType: scaffold.ResourceInfoTypeUrl,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "spring.datasource.url",
					ResourceType:     scaffold.ResourceTypeDbMySQL,
					ResourceInfoType: scaffold.ResourceInfoTypeJdbcUrl,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "spring.datasource.username",
					ResourceType:     scaffold.ResourceTypeDbMySQL,
					ResourceInfoType: scaffold.ResourceInfoTypeUsername,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "spring.datasource.password",
					ResourceType:     scaffold.ResourceTypeDbMySQL,
					ResourceInfoType: scaffold.ResourceInfoTypePassword,
				},
			}, nil
		case internal.AuthTypeUserAssignedManagedIdentity:
			return []scaffold.Env{
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "MYSQL_USERNAME",
					ResourceType:     scaffold.ResourceTypeDbMySQL,
					ResourceInfoType: scaffold.ResourceInfoTypeUsername,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "MYSQL_HOST",
					ResourceType:     scaffold.ResourceTypeDbMySQL,
					ResourceInfoType: scaffold.ResourceInfoTypeHost,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "MYSQL_PORT",
					ResourceType:     scaffold.ResourceTypeDbMySQL,
					ResourceInfoType: scaffold.ResourceInfoTypePort,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "MYSQL_DATABASE",
					ResourceType:     scaffold.ResourceTypeDbMySQL,
					ResourceInfoType: scaffold.ResourceInfoTypeDatabaseName,
				},
				{
					EnvType: scaffold.EnvTypeResourceConnectionServiceConnectorCreated,
					Name:    "spring.datasource.url",
				},
				{
					EnvType: scaffold.EnvTypeResourceConnectionServiceConnectorCreated,
					Name:    "spring.datasource.username",
				},
				{
					EnvType: scaffold.EnvTypeResourceConnectionServiceConnectorCreated,
					Name:    "spring.datasource.azure.passwordless-enabled",
				},
			}, nil
		default:
			return []scaffold.Env{}, unsupportedAuthTypeError(resourceType, authType)
		}
	case ResourceTypeDbRedis:
		switch authType {
		case internal.AuthTypePassword:
			return []scaffold.Env{
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "REDIS_HOST",
					ResourceType:     scaffold.ResourceTypeDbRedis,
					ResourceInfoType: scaffold.ResourceInfoTypeHost,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "REDIS_PORT",
					ResourceType:     scaffold.ResourceTypeDbRedis,
					ResourceInfoType: scaffold.ResourceInfoTypePort,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "REDIS_ENDPOINT",
					ResourceType:     scaffold.ResourceTypeDbRedis,
					ResourceInfoType: scaffold.ResourceInfoTypeEndpoint,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "REDIS_URL",
					ResourceType:     scaffold.ResourceTypeDbRedis,
					ResourceInfoType: scaffold.ResourceInfoTypeUrl,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "REDIS_PASSWORD",
					ResourceType:     scaffold.ResourceTypeDbRedis,
					ResourceInfoType: scaffold.ResourceInfoTypePassword,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "spring.data.redis.url",
					ResourceType:     scaffold.ResourceTypeDbRedis,
					ResourceInfoType: scaffold.ResourceInfoTypeUrl,
				},
			}, nil
		default:
			return []scaffold.Env{}, unsupportedAuthTypeError(resourceType, authType)
		}
	case ResourceTypeDbMongo:
		switch authType {
		case internal.AuthTypeUserAssignedManagedIdentity:
			return []scaffold.Env{
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "MONGODB_URL",
					ResourceType:     scaffold.ResourceTypeDbMongo,
					ResourceInfoType: scaffold.ResourceInfoTypeUrl,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "spring.data.mongodb.uri",
					ResourceType:     scaffold.ResourceTypeDbMongo,
					ResourceInfoType: scaffold.ResourceInfoTypeUrl,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "spring.data.mongodb.database",
					ResourceType:     scaffold.ResourceTypeDbMongo,
					ResourceInfoType: scaffold.ResourceInfoTypeDatabaseName,
				},
			}, nil
		default:
			return []scaffold.Env{}, unsupportedAuthTypeError(resourceType, authType)
		}
	case ResourceTypeDbCosmos:
		switch authType {
		case internal.AuthTypeUserAssignedManagedIdentity:
			return []scaffold.Env{
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "spring.cloud.azure.cosmos.endpoint",
					ResourceType:     scaffold.ResourceTypeDbCosmos,
					ResourceInfoType: scaffold.ResourceInfoTypeEndpoint,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "spring.cloud.azure.cosmos.database",
					ResourceType:     scaffold.ResourceTypeDbCosmos,
					ResourceInfoType: scaffold.ResourceInfoTypeDatabaseName,
				},
			}, nil
		default:
			return []scaffold.Env{}, unsupportedAuthTypeError(resourceType, authType)
		}
	case ResourceTypeMessagingServiceBus:
		if infraSpec.AzureServiceBus.IsJms {
			switch authType {
			case internal.AuthTypeUserAssignedManagedIdentity:
				return []scaffold.Env{
					{
						EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
						Name:           "spring.jms.servicebus.pricing-tier",
						PlainTextValue: "premium",
					},
					{
						EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
						Name:           "spring.jms.servicebus.passwordless-enabled",
						PlainTextValue: "true",
					},
					{
						EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
						Name:           "spring.jms.servicebus.credential.managed-identity-enabled",
						PlainTextValue: "true",
					},
					{
						EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
						Name:           "spring.jms.servicebus.credential.client-id",
						PlainTextValue: scaffold.PlaceHolderForServiceIdentityClientId(),
					},
					{
						EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
						Name:             "spring.jms.servicebus.namespace",
						ResourceType:     scaffold.ResourceTypeMessagingServiceBus,
						ResourceInfoType: scaffold.ResourceInfoTypeNamespace,
					},
					{
						EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
						Name:           "spring.jms.servicebus.connection-string",
						PlainTextValue: "",
					},
				}, nil
			case internal.AuthTypeConnectionString:
				return []scaffold.Env{
					{
						EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
						Name:           "spring.jms.servicebus.pricing-tier",
						PlainTextValue: "premium",
					},
					{
						EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
						Name:             "spring.jms.servicebus.connection-string",
						ResourceType:     scaffold.ResourceTypeMessagingServiceBus,
						ResourceInfoType: scaffold.ResourceInfoTypeConnectionString,
					},
					{
						EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
						Name:           "spring.jms.servicebus.passwordless-enabled",
						PlainTextValue: "false",
					},
					{
						EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
						Name:           "spring.jms.servicebus.credential.managed-identity-enabled",
						PlainTextValue: "false",
					},
					{
						EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
						Name:           "spring.jms.servicebus.credential.client-id",
						PlainTextValue: "",
					},
					{
						EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
						Name:           "spring.jms.servicebus.namespace",
						PlainTextValue: "",
					},
				}, nil
			default:
				return []scaffold.Env{}, unsupportedResourceTypeError(resourceType)
			}
		} else {
			// service bus, not jms
			switch authType {
			case internal.AuthTypeUserAssignedManagedIdentity:
				return []scaffold.Env{
					// Not add this: spring.cloud.azure.servicebus.connection-string = ""
					// because of this: https://github.com/Azure/azure-sdk-for-java/issues/42880
					{
						EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
						Name:           "spring.cloud.azure.servicebus.credential.managed-identity-enabled",
						PlainTextValue: "true",
					},
					{
						EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
						Name:           "spring.cloud.azure.servicebus.credential.client-id",
						PlainTextValue: scaffold.PlaceHolderForServiceIdentityClientId(),
					},
					{
						EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
						Name:             "spring.cloud.azure.servicebus.namespace",
						ResourceType:     scaffold.ResourceTypeMessagingServiceBus,
						ResourceInfoType: scaffold.ResourceInfoTypeNamespace,
					},
				}, nil
			case internal.AuthTypeConnectionString:
				return []scaffold.Env{
					{
						EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
						Name:             "spring.cloud.azure.servicebus.namespace",
						ResourceType:     scaffold.ResourceTypeMessagingServiceBus,
						ResourceInfoType: scaffold.ResourceInfoTypeNamespace,
					},
					{
						EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
						Name:             "spring.cloud.azure.servicebus.connection-string",
						ResourceType:     scaffold.ResourceTypeMessagingServiceBus,
						ResourceInfoType: scaffold.ResourceInfoTypeConnectionString,
					},
					{
						EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
						Name:           "spring.cloud.azure.servicebus.credential.managed-identity-enabled",
						PlainTextValue: "false",
					},
					{
						EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
						Name:           "spring.cloud.azure.servicebus.credential.client-id",
						PlainTextValue: "",
					},
				}, nil
			default:
				return []scaffold.Env{}, unsupportedResourceTypeError(resourceType)
			}
		}
	case ResourceTypeMessagingKafka:
		// event hubs for kafka
		var springBootVersionDecidedInformation []scaffold.Env
		if strings.HasPrefix(infraSpec.AzureEventHubs.SpringBootVersion, "2.") {
			springBootVersionDecidedInformation = []scaffold.Env{
				{
					EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
					Name:           "spring.cloud.stream.binders.kafka.environment.spring.main.sources",
					PlainTextValue: "com.azure.spring.cloud.autoconfigure.eventhubs.kafka.AzureEventHubsKafkaAutoConfiguration",
				},
			}
		} else {
			springBootVersionDecidedInformation = []scaffold.Env{
				{
					EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
					Name:           "spring.cloud.stream.binders.kafka.environment.spring.main.sources",
					PlainTextValue: "com.azure.spring.cloud.autoconfigure.implementation.eventhubs.kafka.AzureEventHubsKafkaAutoConfiguration",
				},
			}
		}
		var commonInformation []scaffold.Env
		switch authType {
		case internal.AuthTypeUserAssignedManagedIdentity:
			commonInformation = []scaffold.Env{
				// Not add this: spring.cloud.azure.eventhubs.connection-string = ""
				// because of this: https://github.com/Azure/azure-sdk-for-java/issues/42880
				{
					EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
					Name:           "spring.cloud.azure.eventhubs.credential.managed-identity-enabled",
					PlainTextValue: "true",
				},
				{
					EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
					Name:           "spring.cloud.azure.eventhubs.credential.client-id",
					PlainTextValue: scaffold.PlaceHolderForServiceIdentityClientId(),
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "spring.cloud.stream.kafka.binder.brokers",
					ResourceType:     scaffold.ResourceTypeMessagingKafka,
					ResourceInfoType: scaffold.ResourceInfoTypeEndpoint,
				},
			}
		case internal.AuthTypeConnectionString:
			commonInformation = []scaffold.Env{
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "spring.cloud.stream.kafka.binder.brokers",
					ResourceType:     scaffold.ResourceTypeMessagingKafka,
					ResourceInfoType: scaffold.ResourceInfoTypeEndpoint,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "spring.cloud.azure.eventhubs.connection-string",
					ResourceType:     scaffold.ResourceTypeMessagingKafka,
					ResourceInfoType: scaffold.ResourceInfoTypeConnectionString,
				},
				{
					EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
					Name:           "spring.cloud.azure.eventhubs.credential.managed-identity-enabled",
					PlainTextValue: "false",
				},
				{
					EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
					Name:           "spring.cloud.azure.eventhubs.credential.client-id",
					PlainTextValue: "",
				},
			}
		default:
			return []scaffold.Env{}, unsupportedAuthTypeError(resourceType, authType)
		}
		return mergeEnvWithDuplicationCheck(springBootVersionDecidedInformation, commonInformation)
	case ResourceTypeMessagingEventHubs:
		switch authType {
		case internal.AuthTypeUserAssignedManagedIdentity:
			return []scaffold.Env{
				// Not add this: spring.cloud.azure.eventhubs.connection-string = ""
				// because of this: https://github.com/Azure/azure-sdk-for-java/issues/42880
				{
					EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
					Name:           "spring.cloud.azure.eventhubs.credential.managed-identity-enabled",
					PlainTextValue: "true",
				},
				{
					EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
					Name:           "spring.cloud.azure.eventhubs.credential.client-id",
					PlainTextValue: scaffold.PlaceHolderForServiceIdentityClientId(),
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "spring.cloud.azure.eventhubs.namespace",
					ResourceType:     scaffold.ResourceTypeMessagingEventHubs,
					ResourceInfoType: scaffold.ResourceInfoTypeNamespace,
				},
			}, nil
		case internal.AuthTypeConnectionString:
			return []scaffold.Env{
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "spring.cloud.azure.eventhubs.namespace",
					ResourceType:     scaffold.ResourceTypeMessagingEventHubs,
					ResourceInfoType: scaffold.ResourceInfoTypeNamespace,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "spring.cloud.azure.eventhubs.connection-string",
					ResourceType:     scaffold.ResourceTypeMessagingEventHubs,
					ResourceInfoType: scaffold.ResourceInfoTypeConnectionString,
				},
				{
					EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
					Name:           "spring.cloud.azure.eventhubs.credential.managed-identity-enabled",
					PlainTextValue: "false",
				},
				{
					EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
					Name:           "spring.cloud.azure.eventhubs.credential.client-id",
					PlainTextValue: "",
				},
			}, nil
		default:
			return []scaffold.Env{}, unsupportedResourceTypeError(resourceType)
		}
	case ResourceTypeStorage:
		switch authType {
		case internal.AuthTypeUserAssignedManagedIdentity:
			return []scaffold.Env{
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "spring.cloud.azure.eventhubs.processor.checkpoint-store.account-name",
					ResourceType:     scaffold.ResourceTypeStorage,
					ResourceInfoType: scaffold.ResourceInfoTypeAccountName,
				},
				{
					EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
					Name:           "spring.cloud.azure.eventhubs.processor.checkpoint-store.credential.managed-identity-enabled",
					PlainTextValue: "true",
				},
				{
					EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
					Name:           "spring.cloud.azure.eventhubs.processor.checkpoint-store.credential.client-id",
					PlainTextValue: scaffold.PlaceHolderForServiceIdentityClientId(),
				},
				{
					EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
					Name:           "spring.cloud.azure.eventhubs.processor.checkpoint-store.connection-string",
					PlainTextValue: "",
				},
			}, nil
		case internal.AuthTypeConnectionString:
			return []scaffold.Env{
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "spring.cloud.azure.eventhubs.processor.checkpoint-store.account-name",
					ResourceType:     scaffold.ResourceTypeStorage,
					ResourceInfoType: scaffold.ResourceInfoTypeAccountName,
				},
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "spring.cloud.azure.eventhubs.processor.checkpoint-store.connection-string",
					ResourceType:     scaffold.ResourceTypeStorage,
					ResourceInfoType: scaffold.ResourceInfoTypeConnectionString,
				},
				{
					EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
					Name:           "spring.cloud.azure.eventhubs.processor.checkpoint-store.credential.managed-identity-enabled",
					PlainTextValue: "false",
				},
				{
					EnvType:        scaffold.EnvTypeResourceConnectionPlainText,
					Name:           "spring.cloud.azure.eventhubs.processor.checkpoint-store.credential.client-id",
					PlainTextValue: "",
				},
			}, nil
		default:
			return []scaffold.Env{}, unsupportedResourceTypeError(resourceType)
		}
	case ResourceTypeOpenAiModel:
		switch authType {
		case internal.AuthTypeUserAssignedManagedIdentity:
			return []scaffold.Env{
				{
					EnvType:          scaffold.EnvTypeResourceConnectionResourceInfo,
					Name:             "AZURE_OPENAI_ENDPOINT",
					ResourceType:     scaffold.ResourceTypeOpenAiModel,
					ResourceInfoType: scaffold.ResourceInfoTypeEndpoint,
				},
			}, nil
		default:
			return []scaffold.Env{}, unsupportedResourceTypeError(resourceType)
		}
	case ResourceTypeHostContainerApp: // todo improve this and delete Frontend and Backend in scaffold.ServiceSpec
		switch authType {
		case internal.AuthTypeUserAssignedManagedIdentity:
			return []scaffold.Env{}, nil
		default:
			return []scaffold.Env{}, unsupportedAuthTypeError(resourceType, authType)
		}
	default:
		return []scaffold.Env{}, unsupportedResourceTypeError(resourceType)
	}
}

func unsupportedResourceTypeError(resourceType ResourceType) error {
	return fmt.Errorf("unsupported resource type, resourceType = %s", resourceType)
}

func unsupportedAuthTypeError(resourceType ResourceType, authType internal.AuthType) error {
	return fmt.Errorf("unsupported auth type, resourceType = %s, authType = %s", resourceType, authType)
}

func mergeEnvWithDuplicationCheck(a []scaffold.Env,
	b []scaffold.Env) ([]scaffold.Env, error) {
	result := append(a, b...)
	seenName := make(map[string]scaffold.Env)
	for _, value := range result {
		if existingValue, exist := seenName[value.Name]; exist {
			if value != existingValue {
				return []scaffold.Env{}, duplicatedEnvError(existingValue, value)
			}
		} else {
			seenName[value.Name] = existingValue
		}
	}
	return result, nil
}

func addNewEnvironmentVariable(serviceSpec *scaffold.ServiceSpec, name string, value string) error {
	merged, err := mergeEnvWithDuplicationCheck(serviceSpec.Envs,
		[]scaffold.Env{
			{
				EnvType:        scaffold.EnvTypePlainText,
				Name:           name,
				PlainTextValue: value,
			},
		},
	)
	if err != nil {
		return err
	}
	serviceSpec.Envs = merged
	return nil
}

func duplicatedEnvError(existingValue scaffold.Env, value scaffold.Env) error {
	return fmt.Errorf(
		"duplicated environment variable. existingValue = %s, value = %s",
		existingValue.ToString(), value.ToString())
}
