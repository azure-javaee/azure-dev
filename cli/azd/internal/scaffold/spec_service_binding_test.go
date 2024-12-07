package scaffold

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToServiceBindingEnvName(t *testing.T) {
	tests := []struct {
		name                  string
		inputResourceType     ServiceType
		inputResourceInfoType ServiceBindingInfoType
		want                  string
	}{
		{
			name:                  "mysql username",
			inputResourceType:     ServiceTypeDbMySQL,
			inputResourceInfoType: ServiceBindingInfoTypeUsername,
			want:                  "$service.binding:db.mysql:username",
		},
		{
			name:                  "postgres password",
			inputResourceType:     ServiceTypeDbPostgres,
			inputResourceInfoType: ServiceBindingInfoTypePassword,
			want:                  "$service.binding:db.postgres:password",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ToServiceBindingEnvValue(tt.inputResourceType, tt.inputResourceInfoType)
			assert.Equal(t, tt.want, actual)
		})
	}
}

func TestIsServiceBindingEnvName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "valid",
			input: "$service.binding:db.postgres:password",
			want:  true,
		},
		{
			name:  "invalid",
			input: "$service.binding:db.postgres:",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isServiceBindingEnvValue(tt.input)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestToServiceTypeAndServiceBindingInfoType(t *testing.T) {
	tests := []struct {
		name                 string
		input                string
		wantResourceType     ServiceType
		wantResourceInfoType ServiceBindingInfoType
	}{
		{
			name:                 "invalid input",
			input:                "$service.binding:db.mysql::username",
			wantResourceType:     "",
			wantResourceInfoType: "",
		},
		{
			name:                 "mysql username",
			input:                "$service.binding:db.mysql:username",
			wantResourceType:     ServiceTypeDbMySQL,
			wantResourceInfoType: ServiceBindingInfoTypeUsername,
		},
		{
			name:                 "postgres password",
			input:                "$service.binding:db.postgres:password",
			wantResourceType:     ServiceTypeDbPostgres,
			wantResourceInfoType: ServiceBindingInfoTypePassword,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceType, resourceInfoType := toServiceTypeAndServiceBindingInfoType(tt.input)
			assert.Equal(t, tt.wantResourceType, resourceType)
			assert.Equal(t, tt.wantResourceInfoType, resourceInfoType)
		})
	}
}
