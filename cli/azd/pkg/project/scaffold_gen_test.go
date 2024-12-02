// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package project

import (
	"reflect"
	"testing"

	"github.com/azure/azure-dev/cli/azd/internal/scaffold"
)

func Test_genBicepParamsFromEnvSubst(t *testing.T) {
	tests := []struct {
		// input
		value         string
		valueIsSecret bool
		// output
		want       string
		wantParams []scaffold.Parameter
	}{
		{"foo", false, "foo", nil},
		{"${MY_VAR}", false, "${myVar}", []scaffold.Parameter{{Name: "myVar", Value: "${MY_VAR}", Type: "string"}}},

		{"${MY_SECRET}", true, "${mySecret}",
			[]scaffold.Parameter{
				{Name: "mySecret", Value: "${MY_SECRET}", Type: "string", Secret: true}}},

		{"Hello, ${world:=okay}!", false, "${world}",
			[]scaffold.Parameter{
				{Name: "world", Value: "${world:=okay}", Type: "string"}}},

		{"${CAT} and ${DOG}", false, "${cat} and ${dog}",
			[]scaffold.Parameter{
				{Name: "cat", Value: "${CAT}", Type: "string"},
				{Name: "dog", Value: "${DOG}", Type: "string"}}},

		{"${DB_HOST:='local'}:${DB_USERNAME:='okay'}", true, "${dbHost}:${dbUsername}",
			[]scaffold.Parameter{
				{Name: "dbHost", Value: "${DB_HOST:='local'}", Type: "string", Secret: true},
				{Name: "dbUsername", Value: "${DB_USERNAME:='okay'}", Type: "string", Secret: true}}},
	}
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			spec := &scaffold.InfraSpec{}
			evaluated := genBicepParamsFromEnvSubst(tt.value, tt.valueIsSecret, spec)
			if tt.want != evaluated {
				t.Errorf("evalEnvValue() evaluatedValue = %v, want %v", evaluated, tt.want)
			}

			for i, param := range tt.wantParams {
				found := false
				for _, generated := range spec.Parameters {
					if generated.Name == param.Name {
						if generated.Secret != param.Secret {
							t.Errorf("evalEnvValue() secret = %v, want %v", generated.Secret, param.Secret)
						}

						if generated.Value != param.Value {
							t.Errorf("evalEnvValue() value = %v, want %v", generated.Value, param.Value)
						}

						if generated.Type != param.Type {
							t.Errorf("evalEnvValue() type = %v, want %v", generated.Type, param.Type)
						}
						found = true
						break
					}
				}

				if !found {
					t.Errorf("evalEnvValue() parameter = %v not found", spec.Parameters[i].Name)
				}
			}
		})
	}
}

func TestServiceConfigEnvToScaffoldEnv(t *testing.T) {
	tests := []struct {
		name            string
		serviceConfig   *ServiceConfig
		wantScaffoldEnv []scaffold.Env
	}{
		{
			name: "eureka-and-config-server",
			serviceConfig: &ServiceConfig{
				Env: map[string]string{
					"eureka.client.serviceUrl.defaultZone": "${DISCOVERY-SERVER_BASE_URL}/eureka",
					"eureka.instance.prefer-ip-address":    "true",
					"spring.config.import":                 "optional:configserver:${CONFIG-SERVER_BASE_URL}",
				},
			},
			wantScaffoldEnv: []scaffold.Env{
				{Name: "eureka.client.serviceUrl.defaultZone", Value: "\\${DISCOVERY-SERVER_BASE_URL}/eureka"},
				{Name: "eureka.instance.prefer-ip-address", Value: "true"},
				{Name: "spring.config.import", Value: "optional:configserver:\\${CONFIG-SERVER_BASE_URL}"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := serviceConfigEnv(tt.serviceConfig)
			if !reflect.DeepEqual(actual, tt.wantScaffoldEnv) {
				t.Errorf("serviceConfigEnv() got = %v, want %v", actual, tt.wantScaffoldEnv)
			}
		})
	}
}
