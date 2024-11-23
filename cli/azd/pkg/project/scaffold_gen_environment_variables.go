package project

import (
	"fmt"
	"github.com/azure/azure-dev/cli/azd/internal"
	"github.com/azure/azure-dev/cli/azd/internal/scaffold"
)

var environmentVariableInformation = map[ResourceType]map[internal.AuthType]scaffold.EnvironmentVariableInformation{
	ResourceTypeDbPostgres: {
		internal.AuthTypePassword: scaffold.EnvironmentVariableInformation{
			StringEnvironmentVariables: []scaffold.StringEnvironmentVariable{
				{
					Name:  "POSTGRES_USERNAME",
					Value: "${postgreSqlDatabaseUser}", // todo manage all variables names
				},
				{
					Name:  "POSTGRES_HOST",
					Value: "${postgreServer.outputs.fqdn}", // todo manage variables like postgreServer
				},
				{
					Name:  "POSTGRES_DATABASE",
					Value: "${postgreSqlDatabaseName}",
				},
				{
					Name:  "POSTGRES_PORT",
					Value: "5432",
				},
				{
					Name:  "spring.datasource.url",
					Value: "jdbc:postgresql://${postgreServer.outputs.fqdn}:5432/${postgreSqlDatabaseName}",
				},
				{
					Name:  "spring.datasource.username",
					Value: "${postgreSqlDatabaseUser}",
				},
			},
			SecretRefEnvironmentVariables: []scaffold.SecretRefEnvironmentVariable{
				{
					Name:      "POSTGRES_URL",
					SecretRef: "postgresql-db-url",
				},
				{
					Name:      "POSTGRES_PASSWORD",
					SecretRef: "postgresql-password",
				},
				{
					Name:      "spring.datasource.password",
					SecretRef: "postgresql-password",
				},
			},
			ValueSecretDefinitions: []scaffold.ValueSecretDefinition{
				{
					SecretName:  "postgresql-db-url",
					SecretValue: "postgresql://${postgreSqlDatabaseUser}:${postgreSqlDatabasePassword}@${postgreServer.outputs.fqdn}:5432/${postgreSqlDatabaseName}",
				},
				{
					SecretName:  "postgresql-password",
					SecretValue: "${postgreSqlDatabasePassword}",
				},
			},
		},
		internal.AuthTypeUserAssignedManagedIdentity: scaffold.EnvironmentVariableInformation{
			StringEnvironmentVariables: []scaffold.StringEnvironmentVariable{
				// Some other environment variables are added by service connector,
				// should not add to bicep generation context
				{
					Name:  "POSTGRES_USERNAME",
					Value: "${postgreSqlDatabaseUser}", // todo manage all variables names
				},
				{
					Name:  "POSTGRES_HOST",
					Value: "${postgreServer.outputs.fqdn}", // todo manage variables like postgreServer
				},
				{
					Name:  "POSTGRES_DATABASE",
					Value: "${postgreSqlDatabaseName}",
				},
				{
					Name:  "POSTGRES_PORT",
					Value: "5432",
				},
			},
		},
	},
	ResourceTypeDbMySQL: {
		internal.AuthTypePassword: scaffold.EnvironmentVariableInformation{
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
		},
		internal.AuthTypeUserAssignedManagedIdentity: scaffold.EnvironmentVariableInformation{
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
		},
	},
	ResourceTypeDbRedis: {
		internal.AuthTypePassword: scaffold.EnvironmentVariableInformation{
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
		},
	},
	ResourceTypeDbMongo: {
		internal.AuthTypeUserAssignedManagedIdentity: scaffold.EnvironmentVariableInformation{
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
			ValueSecretDefinitions: []scaffold.ValueSecretDefinition{},
			KeyVaultSecretDefinitions: []scaffold.KeyVaultSecretDefinition{
				{
					SecretName:  "mongodb-url",
					KeyVaultUrl: "${cosmos.outputs.exportedSecrets['MONGODB-URL'].secretUri}",
				},
			},
		},
	},
	ResourceTypeHostContainerApp: {
		internal.AuthTypeUserAssignedManagedIdentity: scaffold.EnvironmentVariableInformation{},
	},
}

func getAllEnvironmentVariablesForPrint(resourceType ResourceType,
	authType internal.AuthType) (scaffold.EnvironmentVariableInformation, error) {
	information, ok := environmentVariableInformation[resourceType][authType]
	if !ok {
		return scaffold.EnvironmentVariableInformation{},
			fmt.Errorf("cannot get environment variable information, resourceType = %s, authType = %s",
				resourceType, authType)
	}
	additional, err := getAdditionalEnvironmentVariablesForPrint(resourceType, authType)
	if err != nil {
		return scaffold.EnvironmentVariableInformation{}, err
	}
	result, err := mergeWithDuplicationCheck(information, additional)
	if err != nil {
		return scaffold.EnvironmentVariableInformation{}, err
	}
	return result, nil
}

// Return environment variables added by service connector, they do not need to add to scaffold.ServiceSpec
// todo: Now only support springBoot application type. Need to support other types
func getAdditionalEnvironmentVariablesForPrint(resourceType ResourceType,
	authType internal.AuthType) (scaffold.EnvironmentVariableInformation, error) {
	switch resourceType {
	case ResourceTypeDbPostgres:
		switch authType {
		case internal.AuthTypePassword:
			return scaffold.EnvironmentVariableInformation{}, nil
		case internal.AuthTypeUserAssignedManagedIdentity:
			return scaffold.EnvironmentVariableInformation{
				StringEnvironmentVariables: []scaffold.StringEnvironmentVariable{
					{
						Name: "spring.datasource.url",
					},
					{
						Name: "spring.datasource.username",
					},
					{
						Name: "spring.datasource.azure.passwordless-enabled",
					},
				},
			}, nil
		default:
			// return error to make sure every case has been considered.
			return scaffold.EnvironmentVariableInformation{}, fmt.Errorf("unsupported auth type: %s", authType)
		}
	case ResourceTypeDbMySQL:
		switch authType {
		case internal.AuthTypePassword:
			return scaffold.EnvironmentVariableInformation{}, nil
		case internal.AuthTypeUserAssignedManagedIdentity:
			return scaffold.EnvironmentVariableInformation{
				StringEnvironmentVariables: []scaffold.StringEnvironmentVariable{
					{
						Name: "spring.datasource.url",
					},
					{
						Name: "spring.datasource.username",
					},
					{
						Name: "spring.datasource.azure.passwordless-enabled",
					},
				},
			}, nil
		default:
			// return error to make sure every case has been considered.
			return scaffold.EnvironmentVariableInformation{}, fmt.Errorf("unsupported auth type: %s", authType)
		}
	case ResourceTypeDbRedis:
		switch authType {
		case internal.AuthTypePassword:
			return scaffold.EnvironmentVariableInformation{}, nil
		default:
			// return error to make sure every case has been considered.
			return scaffold.EnvironmentVariableInformation{}, fmt.Errorf("unsupported auth type: %s", authType)
		}
	case ResourceTypeDbMongo:
		switch authType {
		case internal.AuthTypeUserAssignedManagedIdentity:
			return scaffold.EnvironmentVariableInformation{}, nil
		default:
			// return error to make sure every case has been considered.
			return scaffold.EnvironmentVariableInformation{}, fmt.Errorf("unsupported auth type: %s", authType)
		}
	case ResourceTypeHostContainerApp:
		return scaffold.EnvironmentVariableInformation{}, nil
	default:
		// return error to make sure every case has been considered.
		return scaffold.EnvironmentVariableInformation{}, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
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
				return scaffold.EnvironmentVariableInformation{}, fmt.Errorf(
					"duplicated environment variable. name = %s, value1 = %s, value2 = %s",
					v.Name, v.Value, existingValue)
			}
		} else {
			seen[v.Name] = existingValue
		}
	}
	for _, v := range result.SecretRefEnvironmentVariables {
		if existingRef, exist := seen[v.Name]; exist {
			if v.SecretRef != existingRef {
				return scaffold.EnvironmentVariableInformation{}, fmt.Errorf(
					"duplicated environment variable. Name = %s, value1 = %s, value2 = %s",
					v.Name, v.SecretRef, existingRef)
			}
		} else {
			seen[v.Name] = existingRef
		}
	}
	for _, v := range result.ValueSecretDefinitions {
		if existingRef, exist := seen[v.SecretName]; exist {
			if v.SecretValue != existingRef {
				return scaffold.EnvironmentVariableInformation{}, fmt.Errorf(
					"duplicated secret definition. Name = %s, value1 = %s, value2 = %s",
					v.SecretName, v.SecretValue, existingRef)
			}
		} else {
			seen[v.SecretName] = existingRef
		}
	}
	for _, v := range result.KeyVaultSecretDefinitions {
		if existingRef, exist := seen[v.SecretName]; exist {
			if v.SecretName != existingRef {
				return scaffold.EnvironmentVariableInformation{}, fmt.Errorf(
					"duplicated secret definition. Name = %s, value1 = %s, value2 = %s",
					v.SecretName, v.KeyVaultUrl, existingRef)
			}
		} else {
			seen[v.SecretName] = existingRef
		}
	}
	return result, nil
}
