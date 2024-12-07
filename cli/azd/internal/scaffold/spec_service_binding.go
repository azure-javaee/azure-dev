package scaffold

import (
	"fmt"
	"strings"

	"github.com/azure/azure-dev/cli/azd/internal"
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

func toServiceTypeAndServiceBindingInfoType(resourceConnectionEnv string) (
	serviceType ServiceType, infoType ServiceBindingInfoType) {
	if !isServiceBindingEnvValue(resourceConnectionEnv) {
		return "", ""
	}
	a := strings.Split(resourceConnectionEnv, ":")
	return ServiceType(a[1]), ServiceBindingInfoType(a[2])
}

func BindToPostgres(serviceSpec *ServiceSpec, postgres *DatabasePostgres) error {
	serviceSpec.DbPostgres = postgres
	envs, err := GetServiceBindingEnvsForPostgres(*postgres)
	if err != nil {
		return err
	}
	serviceSpec.Envs, err = mergeEnvWithDuplicationCheck(serviceSpec.Envs, envs)
	if err != nil {
		return err
	}
	return nil
}

func BindToMySql(serviceSpec *ServiceSpec, mysql *DatabaseMySql) error {
	serviceSpec.DbMySql = mysql
	envs, err := GetServiceBindingEnvsForMysql(*mysql)
	if err != nil {
		return err
	}
	serviceSpec.Envs, err = mergeEnvWithDuplicationCheck(serviceSpec.Envs, envs)
	if err != nil {
		return err
	}
	return nil
}

func BindToMongoDb(serviceSpec *ServiceSpec, mongo *DatabaseCosmosMongo) error {
	serviceSpec.DbCosmosMongo = mongo
	envs, err := GetServiceBindingEnvsForMongo()
	if err != nil {
		return err
	}
	serviceSpec.Envs, err = mergeEnvWithDuplicationCheck(serviceSpec.Envs, envs)
	if err != nil {
		return err
	}
	return nil
}

func BindToCosmosDb(serviceSpec *ServiceSpec, cosmos *DatabaseCosmosAccount) error {
	serviceSpec.DbCosmos = cosmos
	envs, err := GetServiceBindingEnvsForCosmos()
	if err != nil {
		return err
	}
	serviceSpec.Envs, err = mergeEnvWithDuplicationCheck(serviceSpec.Envs, envs)
	if err != nil {
		return err
	}
	return nil
}

func GetServiceBindingEnvsForPostgres(postgres DatabasePostgres) ([]Env, error) {
	switch postgres.AuthType {
	case internal.AuthTypePassword:
		return []Env{
			{
				Name:  "POSTGRES_USERNAME",
				Value: ToServiceBindingEnvValue(ServiceTypeDbPostgres, ServiceBindingInfoTypeUsername),
			},
			{
				Name:  "POSTGRES_PASSWORD",
				Value: ToServiceBindingEnvValue(ServiceTypeDbPostgres, ServiceBindingInfoTypePassword),
			},
			{
				Name:  "POSTGRES_HOST",
				Value: ToServiceBindingEnvValue(ServiceTypeDbPostgres, ServiceBindingInfoTypeHost),
			},
			{
				Name:  "POSTGRES_DATABASE",
				Value: ToServiceBindingEnvValue(ServiceTypeDbPostgres, ServiceBindingInfoTypeDatabaseName),
			},
			{
				Name:  "POSTGRES_PORT",
				Value: ToServiceBindingEnvValue(ServiceTypeDbPostgres, ServiceBindingInfoTypePort),
			},
			{
				Name:  "POSTGRES_URL",
				Value: ToServiceBindingEnvValue(ServiceTypeDbPostgres, ServiceBindingInfoTypeUrl),
			},
			{
				Name:  "spring.datasource.url",
				Value: ToServiceBindingEnvValue(ServiceTypeDbPostgres, ServiceBindingInfoTypeJdbcUrl),
			},
			{
				Name:  "spring.datasource.username",
				Value: ToServiceBindingEnvValue(ServiceTypeDbPostgres, ServiceBindingInfoTypeUsername),
			},
			{
				Name:  "spring.datasource.password",
				Value: ToServiceBindingEnvValue(ServiceTypeDbPostgres, ServiceBindingInfoTypePassword),
			},
		}, nil
	case internal.AuthTypeUserAssignedManagedIdentity:
		return []Env{
			{
				Name:  "POSTGRES_USERNAME",
				Value: ToServiceBindingEnvValue(ServiceTypeDbPostgres, ServiceBindingInfoTypeUsername),
			},
			{
				Name:  "POSTGRES_HOST",
				Value: ToServiceBindingEnvValue(ServiceTypeDbPostgres, ServiceBindingInfoTypeHost),
			},
			{
				Name:  "POSTGRES_DATABASE",
				Value: ToServiceBindingEnvValue(ServiceTypeDbPostgres, ServiceBindingInfoTypeDatabaseName),
			},
			{
				Name:  "POSTGRES_PORT",
				Value: ToServiceBindingEnvValue(ServiceTypeDbPostgres, ServiceBindingInfoTypePort),
			},
			{
				Name:  "spring.datasource.url",
				Value: ToServiceBindingEnvValue(ServiceTypeDbPostgres, ServiceBindingInfoTypeJdbcUrl),
			},
			{
				Name:  "spring.datasource.username",
				Value: ToServiceBindingEnvValue(ServiceTypeDbPostgres, ServiceBindingInfoTypeUsername),
			},
			{
				Name:  "spring.datasource.azure.passwordless-enabled",
				Value: "true",
			},
		}, nil
	default:
		return []Env{}, unsupportedAuthTypeError(ServiceTypeDbPostgres, postgres.AuthType)
	}
}

func GetServiceBindingEnvsForMysql(mysql DatabaseMySql) ([]Env, error) {
	switch mysql.AuthType {
	case internal.AuthTypePassword:
		return []Env{
			{
				Name:  "MYSQL_USERNAME",
				Value: ToServiceBindingEnvValue(ServiceTypeDbMySQL, ServiceBindingInfoTypeUsername),
			},
			{
				Name:  "MYSQL_PASSWORD",
				Value: ToServiceBindingEnvValue(ServiceTypeDbMySQL, ServiceBindingInfoTypePassword),
			},
			{
				Name:  "MYSQL_HOST",
				Value: ToServiceBindingEnvValue(ServiceTypeDbMySQL, ServiceBindingInfoTypeHost),
			},
			{
				Name:  "MYSQL_DATABASE",
				Value: ToServiceBindingEnvValue(ServiceTypeDbMySQL, ServiceBindingInfoTypeDatabaseName),
			},
			{
				Name:  "MYSQL_PORT",
				Value: ToServiceBindingEnvValue(ServiceTypeDbMySQL, ServiceBindingInfoTypePort),
			},
			{
				Name:  "MYSQL_URL",
				Value: ToServiceBindingEnvValue(ServiceTypeDbMySQL, ServiceBindingInfoTypeUrl),
			},
			{
				Name:  "spring.datasource.url",
				Value: ToServiceBindingEnvValue(ServiceTypeDbMySQL, ServiceBindingInfoTypeJdbcUrl),
			},
			{
				Name:  "spring.datasource.username",
				Value: ToServiceBindingEnvValue(ServiceTypeDbMySQL, ServiceBindingInfoTypeUsername),
			},
			{
				Name:  "spring.datasource.password",
				Value: ToServiceBindingEnvValue(ServiceTypeDbMySQL, ServiceBindingInfoTypePassword),
			},
		}, nil
	case internal.AuthTypeUserAssignedManagedIdentity:
		return []Env{
			{
				Name:  "MYSQL_USERNAME",
				Value: ToServiceBindingEnvValue(ServiceTypeDbMySQL, ServiceBindingInfoTypeUsername),
			},
			{
				Name:  "MYSQL_HOST",
				Value: ToServiceBindingEnvValue(ServiceTypeDbMySQL, ServiceBindingInfoTypeHost),
			},
			{
				Name:  "MYSQL_PORT",
				Value: ToServiceBindingEnvValue(ServiceTypeDbMySQL, ServiceBindingInfoTypePort),
			},
			{
				Name:  "MYSQL_DATABASE",
				Value: ToServiceBindingEnvValue(ServiceTypeDbMySQL, ServiceBindingInfoTypeDatabaseName),
			},
			{
				Name:  "spring.datasource.url",
				Value: ToServiceBindingEnvValue(ServiceTypeDbMySQL, ServiceBindingInfoTypeJdbcUrl),
			},
			{
				Name:  "spring.datasource.username",
				Value: ToServiceBindingEnvValue(ServiceTypeDbMySQL, ServiceBindingInfoTypeUsername),
			},
			{
				Name:  "spring.datasource.azure.passwordless-enabled",
				Value: "true",
			},
		}, nil
	default:
		return []Env{}, unsupportedAuthTypeError(ServiceTypeDbMySQL, mysql.AuthType)
	}
}

func GetServiceBindingEnvsForMongo() ([]Env, error) {
	return []Env{
		{
			Name:  "MONGODB_URL",
			Value: ToServiceBindingEnvValue(ServiceTypeDbMongo, ServiceBindingInfoTypeUrl),
		},
		{
			Name:  "spring.data.mongodb.uri",
			Value: ToServiceBindingEnvValue(ServiceTypeDbMongo, ServiceBindingInfoTypeUrl),
		},
		{
			Name:  "spring.data.mongodb.database",
			Value: ToServiceBindingEnvValue(ServiceTypeDbMongo, ServiceBindingInfoTypeDatabaseName),
		},
	}, nil
}

func GetServiceBindingEnvsForCosmos() ([]Env, error) {
	return []Env{
		{
			Name: "spring.cloud.azure.cosmos.endpoint",
			Value: ToServiceBindingEnvValue(
				ServiceTypeDbCosmos, ServiceBindingInfoTypeEndpoint),
		},
		{
			Name: "spring.cloud.azure.cosmos.database",
			Value: ToServiceBindingEnvValue(
				ServiceTypeDbCosmos, ServiceBindingInfoTypeDatabaseName),
		},
	}, nil
}

func unsupportedAuthTypeError(serviceType ServiceType, authType internal.AuthType) error {
	return fmt.Errorf("unsupported auth type, serviceType = %s, authType = %s", serviceType, authType)
}

func mergeEnvWithDuplicationCheck(a []Env, b []Env) ([]Env, error) {
	ab := append(a, b...)
	var result []Env
	seenName := make(map[string]Env)
	for _, value := range ab {
		if existingValue, exist := seenName[value.Name]; exist {
			if value != existingValue {
				return []Env{}, duplicatedEnvError(existingValue, value)
			}
		} else {
			seenName[value.Name] = value
			result = append(result, value)
		}
	}
	return result, nil
}

func duplicatedEnvError(existingValue Env, newValue Env) error {
	return fmt.Errorf(
		"duplicated environment variable. existingValue = %s, newValue = %s",
		existingValue, newValue,
	)
}
