package scaffold

import (
	"fmt"
	"strings"
)

func ToBicepEnv(env Env) BicepEnv {
	switch env.EnvType {
	case EnvTypeResourceConnectionServiceConnectorCreated:
		return BicepEnv{
			BicepEnvType: BicepEnvTypeOthers,
			Name:         env.Name,
		}
	case EnvTypePlainText, EnvTypeResourceConnectionPlainText:
		return BicepEnv{
			BicepEnvType:   BicepEnvTypePlainText,
			Name:           env.Name,
			PlainTextValue: toBicepEnvPlainTextValue(env.PlainTextValue),
		}
	case EnvTypeResourceConnectionResourceInfo:
		value, ok := resourceSpecificBicepEnvValue[env.ResourceType][env.ResourceInfoType]
		if !ok {
			panic(unsupportedType(env))
		}
		if isSecret(env.ResourceInfoType) {
			if isKeyVaultSecret(value) {
				return BicepEnv{
					BicepEnvType: BicepEnvTypeKeyVaultSecret,
					Name:         env.Name,
					SecretName:   secretName(env),
					SecretValue:  unwrapKeyVaultSecretValue(value),
				}
			} else {
				return BicepEnv{
					BicepEnvType: BicepEnvTypeSecret,
					Name:         env.Name,
					SecretName:   secretName(env),
					SecretValue:  value,
				}
			}
		} else {
			return BicepEnv{
				BicepEnvType:   BicepEnvTypePlainText,
				Name:           env.Name,
				PlainTextValue: value,
			}
		}
	default:
		panic(unsupportedType(env))
	}
}

// inputStringExample -> 'inputStringExample'
func addQuotation(input string) string {
	return fmt.Sprintf("'%s'", input)
}

// 'inputStringExample' -> 'inputStringExample'
// '${inputSingleVariableExample}' -> inputSingleVariableExample
// '${HOST}:${PORT}' -> '${HOST}:${PORT}'
func removeQuotationIfItIsASingleVariable(input string) string {
	prefix := "'${"
	suffix := "}'"
	if strings.HasPrefix(input, prefix) && strings.HasSuffix(input, suffix) {
		prefixTrimmed := strings.TrimPrefix(input, prefix)
		trimmed := strings.TrimSuffix(prefixTrimmed, suffix)
		if strings.IndexAny(trimmed, "}") == -1 {
			return trimmed
		} else {
			return input
		}
	} else {
		return input
	}
}

// The BicepEnv.PlainTextValue is handled as variable by default.
// If the value is string, it should contain (').
// Here are some examples of input and output:
// inputStringExample -> 'inputStringExample'
// ${inputSingleVariableExample} -> inputSingleVariableExample
// ${HOST}:${PORT} -> '${HOST}:${PORT}'
func toBicepEnvPlainTextValue(input string) string {
	return removeQuotationIfItIsASingleVariable(addQuotation(input))
}

// BicepEnv
//
// For Name and SecretName, they are handled as string by default.
// Which means quotation will be added before they are used in bicep file, because they are always string value.
//
// For PlainTextValue and SecretValue, they are handled as variable by default.
// When they are string value, quotation should be contained by themselves.
// Set variable as default is mainly to avoid this problem:
// https://learn.microsoft.com/en-us/azure/azure-resource-manager/bicep/linter-rule-simplify-interpolation
type BicepEnv struct {
	BicepEnvType   BicepEnvType
	Name           string
	PlainTextValue string
	SecretName     string
	SecretValue    string
}

type BicepEnvType string

const (
	BicepEnvTypePlainText      BicepEnvType = "plainText"
	BicepEnvTypeSecret         BicepEnvType = "secret"
	BicepEnvTypeKeyVaultSecret BicepEnvType = "keyVaultSecret"
	BicepEnvTypeOthers         BicepEnvType = "others" // This will not be added in to bicep file
)

// Note: If the value is string, it should contain quotation, otherwise it will be viewed as variable.
var resourceSpecificBicepEnvValue = map[ResourceType]map[ResourceInfoType]string{
	ResourceTypeDbPostgres: {
		ResourceInfoTypeHost:         "postgreServer.outputs.fqdn",
		ResourceInfoTypePort:         "'5432'",
		ResourceInfoTypeDatabaseName: "postgreSqlDatabaseName",
		ResourceInfoTypeUsername:     "postgreSqlDatabaseUser",
		ResourceInfoTypePassword:     "postgreSqlDatabasePassword",
		ResourceInfoTypeUrl:          "'postgresql://${postgreSqlDatabaseUser}:${postgreSqlDatabasePassword}@${postgreServer.outputs.fqdn}:5432/${postgreSqlDatabaseName}'",
		ResourceInfoTypeJdbcUrl:      "'jdbc:postgresql://${postgreServer.outputs.fqdn}:5432/${postgreSqlDatabaseName}'",
	},
	ResourceTypeDbMySQL: {
		ResourceInfoTypeHost:         "mysqlServer.outputs.fqdn",
		ResourceInfoTypePort:         "'3306'",
		ResourceInfoTypeDatabaseName: "mysqlDatabaseName",
		ResourceInfoTypeUsername:     "mysqlDatabaseUser",
		ResourceInfoTypePassword:     "mysqlDatabasePassword",
		ResourceInfoTypeUrl:          "'mysql://${mysqlDatabaseUser}:${mysqlDatabasePassword}@${mysqlServer.outputs.fqdn}:3306/${mysqlDatabaseName}'",
		ResourceInfoTypeJdbcUrl:      "'jdbc:mysql://${mysqlServer.outputs.fqdn}:3306/${mysqlDatabaseName}'",
	},
	ResourceTypeDbRedis: {
		ResourceInfoTypeHost:     "redis.outputs.hostName",
		ResourceInfoTypePort:     "string(redis.outputs.sslPort)",
		ResourceInfoTypeEndpoint: "'${redis.outputs.hostName}:${redis.outputs.sslPort}'",
		ResourceInfoTypePassword: wrapToKeyVaultSecretValue("'${keyVault.outputs.uri}secrets/REDIS-PASSWORD'"),
		ResourceInfoTypeUrl:      wrapToKeyVaultSecretValue("'${keyVault.outputs.uri}secrets/REDIS-URL'"),
	},
	ResourceTypeDbMongo: {
		ResourceInfoTypeDatabaseName: "mongoDatabaseName",
		ResourceInfoTypeUrl:          wrapToKeyVaultSecretValue("cosmos.outputs.exportedSecrets['MONGODB-URL'].secretUri"),
	},
	ResourceTypeDbCosmos: {
		ResourceInfoTypeEndpoint:     "cosmos.outputs.endpoint",
		ResourceInfoTypeDatabaseName: "cosmosDatabaseName",
	},
	ResourceTypeMessagingServiceBus: {
		ResourceInfoTypeNamespace:        "serviceBusNamespace.outputs.name",
		ResourceInfoTypeConnectionString: wrapToKeyVaultSecretValue("serviceBusConnectionString.outputs.keyVaultUrl"),
	},
	ResourceTypeMessagingEventHubs: {
		ResourceInfoTypeNamespace:        "eventHubNamespace.outputs.name",
		ResourceInfoTypeConnectionString: wrapToKeyVaultSecretValue("'${keyVault.outputs.uri}secrets/EVENT-HUBS-CONNECTION-STRING'"),
	},
	ResourceTypeMessagingKafka: {
		ResourceInfoTypeEndpoint:         "${eventHubNamespace.outputs.name}.servicebus.windows.net:909",
		ResourceInfoTypeConnectionString: wrapToKeyVaultSecretValue("'${keyVault.outputs.uri}secrets/EVENT-HUBS-CONNECTION-STRING'"),
	},
	ResourceTypeStorage: {
		ResourceInfoTypeAccountName:      "storageAccountName",
		ResourceInfoTypeConnectionString: wrapToKeyVaultSecretValue("'${keyVault.outputs.uri}secrets/STORAGE-ACCOUNT-CONNECTION-STRING'"),
	},
	ResourceTypeOpenAiModel: {
		ResourceInfoTypeEndpoint: "account.outputs.endpoint",
	},
	ResourceTypeHostContainerApp: {},
}

func unsupportedType(env Env) string {
	return fmt.Sprintf("unsupported connection info type for resource type. "+
		"resourceType = %s, connectionInfoType = %s", env.ResourceType, env.ResourceInfoType)
}

func PlaceHolderForServiceIdentityClientId() string {
	return "__PlaceHolderForServiceIdentityClientId"
}

func isSecret(info ResourceInfoType) bool {
	return info == ResourceInfoTypePassword || info == ResourceInfoTypeUrl || info == ResourceInfoTypeConnectionString
}

func secretName(env Env) string {
	name := fmt.Sprintf("%s-%s", env.ResourceType, env.ResourceInfoType)
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
