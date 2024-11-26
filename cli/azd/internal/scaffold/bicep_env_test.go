package scaffold

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetEnvBicepInfo(t *testing.T) {
	tests := []struct {
		name string
		in   Env
		want BicepEnv
	}{
		{
			name: "Service connector created",
			in: Env{
				EnvType: EnvTypeResourceConnectionServiceConnectorCreated,
				Name:    "spring.datasource.url",
			},
			want: BicepEnv{
				BicepEnvType: BicepEnvTypeOthers,
				Name:         "spring.datasource.url",
			},
		},
		{
			name: "Plain text",
			in: Env{
				EnvType:        EnvTypePlainText,
				Name:           "enable-customer-related-feature",
				PlainTextValue: "true",
			},
			want: BicepEnv{
				BicepEnvType:   BicepEnvTypePlainText,
				Name:           "enable-customer-related-feature",
				PlainTextValue: "'true'",
			},
		},
		{
			name: "Plain text from EnvTypeResourceConnectionPlainText",
			in: Env{
				EnvType:        EnvTypeResourceConnectionPlainText,
				Name:           "spring.jms.servicebus.pricing-tier",
				PlainTextValue: "premium",
			},
			want: BicepEnv{
				BicepEnvType:   BicepEnvTypePlainText,
				Name:           "spring.jms.servicebus.pricing-tier",
				PlainTextValue: "'premium'",
			},
		},
		{
			name: "Plain text from EnvTypeResourceConnectionResourceInfo",
			in: Env{
				EnvType:          EnvTypeResourceConnectionResourceInfo,
				Name:             "POSTGRES_PORT",
				ResourceType:     ResourceTypeDbPostgres,
				ResourceInfoType: ResourceInfoTypePort,
			},
			want: BicepEnv{
				BicepEnvType:   BicepEnvTypePlainText,
				Name:           "POSTGRES_PORT",
				PlainTextValue: "'5432'",
			},
		},
		{
			name: "Secret",
			in: Env{
				EnvType:          EnvTypeResourceConnectionResourceInfo,
				Name:             "POSTGRES_PASSWORD",
				ResourceType:     ResourceTypeDbPostgres,
				ResourceInfoType: ResourceInfoTypePassword,
			},
			want: BicepEnv{
				BicepEnvType: BicepEnvTypeSecret,
				Name:         "POSTGRES_PASSWORD",
				SecretName:   "db-postgres-password",
				SecretValue:  "postgreSqlDatabasePassword",
			},
		},
		{
			name: "KeuVault Secret",
			in: Env{
				EnvType:          EnvTypeResourceConnectionResourceInfo,
				Name:             "REDIS_PASSWORD",
				ResourceType:     ResourceTypeDbRedis,
				ResourceInfoType: ResourceInfoTypePassword,
			},
			want: BicepEnv{
				BicepEnvType: BicepEnvTypeKeyVaultSecret,
				Name:         "REDIS_PASSWORD",
				SecretName:   "db-redis-password",
				SecretValue:  "${keyVault.outputs.uri}secrets/REDIS-PASSWORD",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ToBicepEnv(tt.in)
			assert.Equal(t, tt.want, actual)
		})
	}
}
