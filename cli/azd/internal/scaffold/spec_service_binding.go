package scaffold

import (
	"fmt"
	"strings"
)

// todo merge ServiceType and project.ResourceType
// Not use project.ResourceType because it will cause cycle import.
// Not merge it in current PR to avoid conflict with upstream main branch.
// Solution proposal: define a ServiceType in lower level that can be used both in scaffold and project package.

type ServiceType string

const (
	ServiceTypeDbRedis             ServiceType = "db.redis"
	ServiceTypeDbPostgres          ServiceType = "db.postgres"
	ServiceTypeDbMySQL             ServiceType = "db.mysql"
	ServiceTypeDbMongo             ServiceType = "db.mongo"
	ServiceTypeDbCosmos            ServiceType = "db.cosmos"
	ServiceTypeHostContainerApp    ServiceType = "host.containerapp"
	ServiceTypeOpenAiModel         ServiceType = "ai.openai.model"
	ServiceTypeMessagingServiceBus ServiceType = "messaging.servicebus"
	ServiceTypeMessagingEventHubs  ServiceType = "messaging.eventhubs"
	ServiceTypeMessagingKafka      ServiceType = "messaging.kafka"
	ServiceTypeStorage             ServiceType = "storage"
)

type ServiceBindingInfoType string

const (
	ServiceBindingInfoTypeHost             ServiceBindingInfoType = "host"
	ServiceBindingInfoTypePort             ServiceBindingInfoType = "port"
	ServiceBindingInfoTypeEndpoint         ServiceBindingInfoType = "endpoint"
	ServiceBindingInfoTypeDatabaseName     ServiceBindingInfoType = "databaseName"
	ServiceBindingInfoTypeNamespace        ServiceBindingInfoType = "namespace"
	ServiceBindingInfoTypeAccountName      ServiceBindingInfoType = "accountName"
	ServiceBindingInfoTypeUsername         ServiceBindingInfoType = "username"
	ServiceBindingInfoTypePassword         ServiceBindingInfoType = "password"
	ServiceBindingInfoTypeUrl              ServiceBindingInfoType = "url"
	ServiceBindingInfoTypeJdbcUrl          ServiceBindingInfoType = "jdbcUrl"
	ServiceBindingInfoTypeConnectionString ServiceBindingInfoType = "connectionString"
)

var serviceBindingEnvValuePrefix = "$service.binding"

func isServiceBindingEnvValue(env string) bool {
	if !strings.HasPrefix(env, serviceBindingEnvValuePrefix) {
		return false
	}
	a := strings.Split(env, ":")
	if len(a) != 3 {
		return false
	}
	return a[0] != "" && a[1] != "" && a[2] != ""
}

func ToServiceBindingEnvValue(resourceType ServiceType, resourceInfoType ServiceBindingInfoType) string {
	return fmt.Sprintf("%s:%s:%s", serviceBindingEnvValuePrefix, resourceType, resourceInfoType)
}

func toServiceTypeAndServiceBindingInfoType(resourceConnectionEnv string) (serviceType ServiceType,
	infoType ServiceBindingInfoType) {
	if !isServiceBindingEnvValue(resourceConnectionEnv) {
		return "", ""
	}
	a := strings.Split(resourceConnectionEnv, ":")
	return ServiceType(a[1]), ServiceBindingInfoType(a[2])
}
