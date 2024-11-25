package scaffold

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetEnvBicepInfo(t *testing.T) {
	tests := []struct {
		name string
		in   ResourceConnectionEnv
		want BicepEnv
	}{
		{
			name: "Service connector created",
			in: ResourceConnectionEnv{
				ResourceConnectionEnvType: ResourceConnectionEnvTypeServiceConnectorCreated,
				Name:                      "spring.datasource.url",
			},
			want: BicepEnv{
				BicepEnvType: BicepEnvTypeServiceConnectorCreated,
				Name:         "spring.datasource.url",
			},
		},
		{
			name: "Plain text created by ResourceConnectionEnvTypePlainText",
			in: ResourceConnectionEnv{
				ResourceConnectionEnvType: ResourceConnectionEnvTypePlainText,
				Name:                      "spring.jms.servicebus.pricing-tier",
				PlainTextValue:            "premium",
			},
			want: BicepEnv{
				BicepEnvType: BicepEnvTypePlainText,
				Name:         "spring.jms.servicebus.pricing-tier",
				Value:        "premium",
			},
		},
		{
			name: "Plain text created by ResourceConnectionEnvTypeResourceSpecific",
			in: ResourceConnectionEnv{
				ResourceConnectionEnvType: ResourceConnectionEnvTypeResourceSpecific,
				Name:                      "POSTGRES_PORT",
				ResourceType:              ResourceTypeDbPostgres,
				ResourceInfoType:          ResourceInfoTypePort,
			},
			want: BicepEnv{
				BicepEnvType: BicepEnvTypePlainText,
				Name:         "POSTGRES_PORT",
				Value:        "5432",
				SecretName:   "",
				SecretValue:  "",
			},
		},
		{
			name: "Secret",
			in: ResourceConnectionEnv{
				ResourceConnectionEnvType: ResourceConnectionEnvTypeResourceSpecific,
				Name:                      "POSTGRES_PASSWORD",
				ResourceType:              ResourceTypeDbPostgres,
				ResourceInfoType:          ResourceInfoTypePassword,
			},
			want: BicepEnv{
				BicepEnvType: BicepEnvTypeSecret,
				Name:         "POSTGRES_PASSWORD",
				SecretName:   "db-postgres-password",
				SecretValue:  "${postgreSqlDatabasePassword}",
			},
		},
		{
			name: "KeuVault Secret",
			in: ResourceConnectionEnv{
				ResourceConnectionEnvType: ResourceConnectionEnvTypeResourceSpecific,
				Name:                      "REDIS_PASSWORD",
				ResourceType:              ResourceTypeDbRedis,
				ResourceInfoType:          ResourceInfoTypePassword,
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
