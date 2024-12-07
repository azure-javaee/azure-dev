package project

import (
	"fmt"

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
