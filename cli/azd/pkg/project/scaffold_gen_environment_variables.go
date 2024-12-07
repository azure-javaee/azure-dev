package project

import (
	"fmt"
	"strings"

	"github.com/azure/azure-dev/cli/azd/internal"
	"github.com/azure/azure-dev/cli/azd/internal/scaffold"
)

func GetResourceConnectionEnvs(usedResource *ResourceConfig,
	infraSpec *scaffold.InfraSpec) ([]scaffold.Env, error) {
	resourceType := usedResource.Type
	authType, err := getAuthType(infraSpec, usedResource.Type)
	if err != nil {
		return []scaffold.Env{}, err
	}
	switch resourceType {
	case ResourceTypeDbRedis:
		switch authType {
		case internal.AuthTypePassword:
			return []scaffold.Env{
				{
					Name: "REDIS_HOST",
					Value: scaffold.ToServiceBindingEnvValue(
						scaffold.ServiceTypeDbRedis, scaffold.ServiceBindingInfoTypeHost),
				},
				{
					Name: "REDIS_PORT",
					Value: scaffold.ToServiceBindingEnvValue(
						scaffold.ServiceTypeDbRedis, scaffold.ServiceBindingInfoTypePort),
				},
				{
					Name: "REDIS_ENDPOINT",
					Value: scaffold.ToServiceBindingEnvValue(
						scaffold.ServiceTypeDbRedis, scaffold.ServiceBindingInfoTypeEndpoint),
				},
				{
					Name: "REDIS_URL",
					Value: scaffold.ToServiceBindingEnvValue(
						scaffold.ServiceTypeDbRedis, scaffold.ServiceBindingInfoTypeUrl),
				},
				{
					Name: "REDIS_PASSWORD",
					Value: scaffold.ToServiceBindingEnvValue(
						scaffold.ServiceTypeDbRedis, scaffold.ServiceBindingInfoTypePassword),
				},
				{
					Name: "spring.data.redis.url",
					Value: scaffold.ToServiceBindingEnvValue(
						scaffold.ServiceTypeDbRedis, scaffold.ServiceBindingInfoTypeUrl),
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
						Name:  "spring.jms.servicebus.pricing-tier",
						Value: "premium",
					},
					{
						Name:  "spring.jms.servicebus.passwordless-enabled",
						Value: "true",
					},
					{
						Name:  "spring.jms.servicebus.credential.managed-identity-enabled",
						Value: "true",
					},
					{
						Name:  "spring.jms.servicebus.credential.client-id",
						Value: scaffold.PlaceHolderForServiceIdentityClientId(),
					},
					{
						Name: "spring.jms.servicebus.namespace",
						Value: scaffold.ToServiceBindingEnvValue(
							scaffold.ServiceTypeMessagingServiceBus, scaffold.ServiceBindingInfoTypeNamespace),
					},
					{
						Name:  "spring.jms.servicebus.connection-string",
						Value: "",
					},
				}, nil
			case internal.AuthTypeConnectionString:
				return []scaffold.Env{
					{
						Name:  "spring.jms.servicebus.pricing-tier",
						Value: "premium",
					},
					{
						Name: "spring.jms.servicebus.connection-string",
						Value: scaffold.ToServiceBindingEnvValue(
							scaffold.ServiceTypeMessagingServiceBus, scaffold.ServiceBindingInfoTypeConnectionString),
					},
					{
						Name:  "spring.jms.servicebus.passwordless-enabled",
						Value: "false",
					},
					{
						Name:  "spring.jms.servicebus.credential.managed-identity-enabled",
						Value: "false",
					},
					{
						Name:  "spring.jms.servicebus.credential.client-id",
						Value: "",
					},
					{
						Name:  "spring.jms.servicebus.namespace",
						Value: "",
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
						Name:  "spring.cloud.azure.servicebus.credential.managed-identity-enabled",
						Value: "true",
					},
					{
						Name:  "spring.cloud.azure.servicebus.credential.client-id",
						Value: scaffold.PlaceHolderForServiceIdentityClientId(),
					},
					{
						Name: "spring.cloud.azure.servicebus.namespace",
						Value: scaffold.ToServiceBindingEnvValue(
							scaffold.ServiceTypeMessagingServiceBus, scaffold.ServiceBindingInfoTypeNamespace),
					},
				}, nil
			case internal.AuthTypeConnectionString:
				return []scaffold.Env{
					{
						Name: "spring.cloud.azure.servicebus.namespace",
						Value: scaffold.ToServiceBindingEnvValue(
							scaffold.ServiceTypeMessagingServiceBus, scaffold.ServiceBindingInfoTypeNamespace),
					},
					{
						Name: "spring.cloud.azure.servicebus.connection-string",
						Value: scaffold.ToServiceBindingEnvValue(
							scaffold.ServiceTypeMessagingServiceBus, scaffold.ServiceBindingInfoTypeConnectionString),
					},
					{
						Name:  "spring.cloud.azure.servicebus.credential.managed-identity-enabled",
						Value: "false",
					},
					{
						Name:  "spring.cloud.azure.servicebus.credential.client-id",
						Value: "",
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
					Name:  "spring.cloud.stream.binders.kafka.environment.spring.main.sources",
					Value: "com.azure.spring.cloud.autoconfigure.eventhubs.kafka.AzureEventHubsKafkaAutoConfiguration",
				},
			}
		} else {
			springBootVersionDecidedInformation = []scaffold.Env{
				{
					Name: "spring.cloud.stream.binders.kafka.environment.spring.main.sources",
					Value: "com.azure.spring.cloud.autoconfigure.implementation.eventhubs.kafka" +
						".AzureEventHubsKafkaAutoConfiguration",
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
					Name: "spring.cloud.stream.kafka.binder.brokers",
					Value: scaffold.ToServiceBindingEnvValue(
						scaffold.ServiceTypeMessagingKafka, scaffold.ServiceBindingInfoTypeEndpoint),
				},
				{
					Name:  "spring.cloud.azure.eventhubs.credential.managed-identity-enabled",
					Value: "true",
				},
				{
					Name:  "spring.cloud.azure.eventhubs.credential.client-id",
					Value: scaffold.PlaceHolderForServiceIdentityClientId(),
				},
			}
		case internal.AuthTypeConnectionString:
			commonInformation = []scaffold.Env{
				{
					Name: "spring.cloud.stream.kafka.binder.brokers",
					Value: scaffold.ToServiceBindingEnvValue(
						scaffold.ServiceTypeMessagingKafka, scaffold.ServiceBindingInfoTypeEndpoint),
				},
				{
					Name: "spring.cloud.azure.eventhubs.connection-string",
					Value: scaffold.ToServiceBindingEnvValue(
						scaffold.ServiceTypeMessagingKafka, scaffold.ServiceBindingInfoTypeConnectionString),
				},
				{
					Name:  "spring.cloud.azure.eventhubs.credential.managed-identity-enabled",
					Value: "false",
				},
				{
					Name:  "spring.cloud.azure.eventhubs.credential.client-id",
					Value: "",
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
					Name:  "spring.cloud.azure.eventhubs.credential.managed-identity-enabled",
					Value: "true",
				},
				{
					Name:  "spring.cloud.azure.eventhubs.credential.client-id",
					Value: scaffold.PlaceHolderForServiceIdentityClientId(),
				},
				{
					Name: "spring.cloud.azure.eventhubs.namespace",
					Value: scaffold.ToServiceBindingEnvValue(
						scaffold.ServiceTypeMessagingEventHubs, scaffold.ServiceBindingInfoTypeNamespace),
				},
			}, nil
		case internal.AuthTypeConnectionString:
			return []scaffold.Env{
				{
					Name: "spring.cloud.azure.eventhubs.namespace",
					Value: scaffold.ToServiceBindingEnvValue(
						scaffold.ServiceTypeMessagingEventHubs, scaffold.ServiceBindingInfoTypeNamespace),
				},
				{
					Name: "spring.cloud.azure.eventhubs.connection-string",
					Value: scaffold.ToServiceBindingEnvValue(
						scaffold.ServiceTypeMessagingEventHubs, scaffold.ServiceBindingInfoTypeConnectionString),
				},
				{
					Name:  "spring.cloud.azure.eventhubs.credential.managed-identity-enabled",
					Value: "false",
				},
				{
					Name:  "spring.cloud.azure.eventhubs.credential.client-id",
					Value: "",
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
					Name: "spring.cloud.azure.eventhubs.processor.checkpoint-store.account-name",
					Value: scaffold.ToServiceBindingEnvValue(
						scaffold.ServiceTypeStorage, scaffold.ServiceBindingInfoTypeAccountName),
				},
				{
					Name:  "spring.cloud.azure.eventhubs.processor.checkpoint-store.credential.managed-identity-enabled",
					Value: "true",
				},
				{
					Name:  "spring.cloud.azure.eventhubs.processor.checkpoint-store.credential.client-id",
					Value: scaffold.PlaceHolderForServiceIdentityClientId(),
				},
				{
					Name:  "spring.cloud.azure.eventhubs.processor.checkpoint-store.connection-string",
					Value: "",
				},
			}, nil
		case internal.AuthTypeConnectionString:
			return []scaffold.Env{
				{
					Name: "spring.cloud.azure.eventhubs.processor.checkpoint-store.account-name",
					Value: scaffold.ToServiceBindingEnvValue(
						scaffold.ServiceTypeStorage, scaffold.ServiceBindingInfoTypeAccountName),
				},
				{
					Name: "spring.cloud.azure.eventhubs.processor.checkpoint-store.connection-string",
					Value: scaffold.ToServiceBindingEnvValue(
						scaffold.ServiceTypeStorage, scaffold.ServiceBindingInfoTypeConnectionString),
				},
				{
					Name:  "spring.cloud.azure.eventhubs.processor.checkpoint-store.credential.managed-identity-enabled",
					Value: "false",
				},
				{
					Name:  "spring.cloud.azure.eventhubs.processor.checkpoint-store.credential.client-id",
					Value: "",
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
					Name: "AZURE_OPENAI_ENDPOINT",
					Value: scaffold.ToServiceBindingEnvValue(
						scaffold.ServiceTypeOpenAiModel, scaffold.ServiceBindingInfoTypeEndpoint),
				},
			}, nil
		default:
			return []scaffold.Env{}, unsupportedResourceTypeError(resourceType)
		}
	case ResourceTypeHostContainerApp: // todo improve this and delete Frontend and Backend in scaffold.ServiceSpec
		switch authType {
		case internal.AuthTypeUnspecified:
			return []scaffold.Env{}, nil
		default:
			return []scaffold.Env{}, unsupportedAuthTypeError(resourceType, authType)
		}
	case ResourceTypeJavaEurekaServer:
		return []scaffold.Env{
			{
				Name:  "eureka.client.register-with-eureka",
				Value: "true",
			},
			{
				Name:  "eureka.client.fetch-registry",
				Value: "true",
			},
			{
				Name:  "eureka.instance.prefer-ip-address",
				Value: "true",
			},
			{
				Name:  "eureka.client.serviceUrl.defaultZone",
				Value: fmt.Sprintf("%s/eureka", scaffold.GetContainerAppHost(usedResource.Name)),
			},
		}, nil
	case ResourceTypeJavaConfigServer:
		return []scaffold.Env{
			{
				Name: "spring.config.import",
				Value: fmt.Sprintf("optional:configserver:%s?fail-fast=true",
					scaffold.GetContainerAppHost(usedResource.Name)),
			},
		}, nil
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
	ab := append(a, b...)
	var result []scaffold.Env
	seenName := make(map[string]scaffold.Env)
	for _, value := range ab {
		if existingValue, exist := seenName[value.Name]; exist {
			if value != existingValue {
				return []scaffold.Env{}, duplicatedEnvError(existingValue, value)
			}
		} else {
			seenName[value.Name] = value
			result = append(result, value)
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

func duplicatedEnvError(existingValue scaffold.Env, newValue scaffold.Env) error {
	return fmt.Errorf("duplicated environment variable. existingValue = %s, newValue = %s",
		existingValue, newValue)
}
