package scaffold

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetEnvBicepInfo(t *testing.T) {
	tests := []struct {
		name string
		in   ResourceConnectionEnv
		want EnvBicepInfo
	}{
		{
			name: "Not secret",
			in: ResourceConnectionEnv{
				Name:               "POSTGRES_PORT",
				ResourceType:       ResourceTypeDbPostgres,
				ConnectionInfoType: Port,
			},
			want: EnvBicepInfo{
				IsSecret:         false,
				IsKeyVaultSecret: false,
				Value:            "5432",
				SecretName:       "",
				SecretValue:      "",
			},
		},
		{
			name: "Is secret, not key vault secret",
			in: ResourceConnectionEnv{
				Name:               "POSTGRES_PASSWORD",
				ResourceType:       ResourceTypeDbPostgres,
				ConnectionInfoType: Password,
			},
			want: EnvBicepInfo{
				IsSecret:         true,
				IsKeyVaultSecret: false,
				Value:            "",
				SecretName:       "db-postgres-password",
				SecretValue:      "${postgreSqlDatabasePassword}",
			},
		},
		{
			name: "Is secret, is vault secret",
			in: ResourceConnectionEnv{
				Name:               "REDIS_PASSWORD",
				ResourceType:       ResourceTypeDbRedis,
				ConnectionInfoType: Password,
			},
			want: EnvBicepInfo{
				IsSecret:         true,
				IsKeyVaultSecret: true,
				Value:            "",
				SecretName:       "db-redis-password",
				SecretValue:      "${keyVault.outputs.uri}secrets/REDIS-PASSWORD",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := GetEnvBicepInfo(tt.in)
			assert.Equal(t, tt.want, actual)
		})
	}
}
