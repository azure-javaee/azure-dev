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
			SecretDefinitions: []scaffold.SecretDefinition{
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

// This is added by service connector, not need to add to scaffold.ServiceSpec
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
		SecretDefinitions:             append(a.SecretDefinitions, b.SecretDefinitions...),
	}
	seen := make(map[string]string)
	for _, v := range result.StringEnvironmentVariables {
		if existingValue, exist := seen[v.Name]; exist {
			if v.Value != existingValue {
				return scaffold.EnvironmentVariableInformation{}, fmt.Errorf(
					"duplicated environment name. name = %s, value1 = %s, value2 = %s",
					v.Name, v.Value, existingValue)
			}
		} else {
			seen[v.Name] = existingValue
		}
	}
	return result, nil
}
