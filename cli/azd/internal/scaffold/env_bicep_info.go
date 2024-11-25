package scaffold

import (
	"fmt"
	"strings"
)

func GetEnvBicepInfo(env ResourceConnectionEnv) EnvBicepInfo {
	value, ok := envValue[env.ResourceType][env.ConnectionInfoType]
	if !ok {
		panic(unsupportedType(env))
	}
	if isSecret(env.ConnectionInfoType) {
		if isKeyVaultSecret(value) {
			return EnvBicepInfo{
				IsSecret:         true,
				IsKeyVaultSecret: true,
				Value:            "",
				SecretName:       secretName(env),
				SecretValue:      unwrapKeyVaultSecretValue(value),
			}
		} else {
			return EnvBicepInfo{
				IsSecret:         true,
				IsKeyVaultSecret: false,
				Value:            "",
				SecretName:       secretName(env),
				SecretValue:      value,
			}
		}
	} else {
		return EnvBicepInfo{
			IsSecret:         false,
			IsKeyVaultSecret: false,
			Value:            value,
			SecretName:       "",
			SecretValue:      "",
		}
	}
}

type EnvBicepInfo struct {
	IsSecret         bool
	IsKeyVaultSecret bool
	Value            string
	SecretName       string
	SecretValue      string
}

var envValue = map[ResourceType]map[ConnectionInfoType]string{
	ResourceTypeDbPostgres: {
		Host:             "${postgreServer.outputs.fqdn}",
		Port:             "5432",
		Endpoint:         "",
		DatabaseName:     "${postgreSqlDatabaseName}",
		Namespace:        "",
		AccountName:      "",
		Username:         "${postgreSqlDatabaseUser}",
		Password:         "${postgreSqlDatabasePassword}",
		Url:              "postgresql://${postgreSqlDatabaseUser}:${postgreSqlDatabasePassword}@${postgreServer.outputs.fqdn}:5432/${postgreSqlDatabaseName}",
		JdbcUrl:          "jdbc:postgresql://${postgreServer.outputs.fqdn}:5432/${postgreSqlDatabaseName}",
		ConnectionString: "",
		IdentityClientId: identityClientId(),
	},
	ResourceTypeDbRedis: {
		Host:             "${redis.outputs.hostName}",
		Port:             "${redis.outputs.sslPort}",
		Endpoint:         "${redis.outputs.hostName}:${redis.outputs.sslPort}",
		DatabaseName:     "",
		Namespace:        "",
		AccountName:      "",
		Username:         "",
		Password:         wrapToKeyVaultSecretValue("${keyVault.outputs.uri}secrets/REDIS-PASSWORD"),
		Url:              wrapToKeyVaultSecretValue("${keyVault.outputs.uri}secrets/REDIS-URL"),
		JdbcUrl:          "",
		ConnectionString: "",
		IdentityClientId: identityClientId(),
	},
}

func unsupportedType(env ResourceConnectionEnv) string {
	return fmt.Sprintf("unsupported connection info type for resource type. "+
		"resourceType = %s, connectionInfoType = %s", env.ResourceType, env.ConnectionInfoType)
}

func identityClientId() string {
	return "__PlaceHolderForServiceIdentityClientId"
}

func isSecret(info ConnectionInfoType) bool {
	return info == Password || info == Url || info == ConnectionString
}

func secretName(env ResourceConnectionEnv) string {
	name := fmt.Sprintf("%s-%s", env.ResourceType, env.ConnectionInfoType)
	lowerCaseName := strings.ToLower(name)
	noDotName := strings.Replace(lowerCaseName, ".", "-", -1)
	noUnderscoreName := strings.Replace(noDotName, "_", "-", -1)
	return noUnderscoreName
}

var keyVaultSecretPrefix = "keyvault:"

func isKeyVaultSecret(value string) bool {
	return strings.HasPrefix(value, keyVaultSecretPrefix)
}

func wrapToKeyVaultSecretValue(value string) string {
	return fmt.Sprintf("%s%s", keyVaultSecretPrefix, value)
}

func unwrapKeyVaultSecretValue(value string) string {
	return strings.TrimPrefix(value, keyVaultSecretPrefix)
}
