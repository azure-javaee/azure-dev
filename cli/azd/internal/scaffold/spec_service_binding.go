package scaffold

import (
	"github.com/azure/azure-dev/cli/azd/internal"
	"github.com/azure/azure-dev/cli/azd/internal/binding"
)

func BindToPostgres(sourceType binding.SourceType, serviceSpec *ServiceSpec, postgres *DatabasePostgres) error {
	serviceSpec.DbPostgres = postgres
	envs, err := binding.GetBindingEnvs(binding.Source{Type: sourceType},
		binding.Target{Type: binding.AzureDatabaseForPostgresql, AuthType: postgres.AuthType})
	if err != nil {
		return err
	}
	serviceSpec.Envs, err = binding.MergeMapWithDuplicationCheck(serviceSpec.Envs, envs)
	if err != nil {
		return err
	}
	return nil
}

func BindToMySql(sourceType binding.SourceType, serviceSpec *ServiceSpec, mysql *DatabaseMySql) error {
	serviceSpec.DbMySql = mysql
	envs, err := binding.GetBindingEnvs(binding.Source{Type: sourceType},
		binding.Target{Type: binding.AzureDatabaseForMysql, AuthType: mysql.AuthType})
	if err != nil {
		return err
	}
	serviceSpec.Envs, err = binding.MergeMapWithDuplicationCheck(serviceSpec.Envs, envs)
	if err != nil {
		return err
	}
	return nil
}

func BindToMongoDb(sourceType binding.SourceType, serviceSpec *ServiceSpec, mongo *DatabaseCosmosMongo) error {
	serviceSpec.DbCosmosMongo = mongo
	envs, err := binding.GetBindingEnvs(binding.Source{Type: sourceType},
		binding.Target{Type: binding.AzureCosmosDBForMongoDB, AuthType: internal.AuthTypeUserAssignedManagedIdentity})
	if err != nil {
		return err
	}
	serviceSpec.Envs, err = binding.MergeMapWithDuplicationCheck(serviceSpec.Envs, envs)
	if err != nil {
		return err
	}
	return nil
}

func BindToCosmosDb(sourceType binding.SourceType, serviceSpec *ServiceSpec, cosmos *DatabaseCosmosAccount) error {
	serviceSpec.DbCosmos = cosmos
	envs, err := binding.GetBindingEnvs(binding.Source{Type: sourceType},
		binding.Target{Type: binding.AzureCosmosDBForNoSQL, AuthType: internal.AuthTypeUserAssignedManagedIdentity})
	if err != nil {
		return err
	}
	serviceSpec.Envs, err = binding.MergeMapWithDuplicationCheck(serviceSpec.Envs, envs)
	if err != nil {
		return err
	}
	return nil
}

func BindToRedis(sourceType binding.SourceType, serviceSpec *ServiceSpec, redis *DatabaseRedis) error {
	serviceSpec.DbRedis = redis
	envs, err := binding.GetBindingEnvs(binding.Source{Type: sourceType},
		binding.Target{Type: binding.AzureCacheForRedis, AuthType: internal.AuthTypePassword})
	if err != nil {
		return err
	}
	serviceSpec.Envs, err = binding.MergeMapWithDuplicationCheck(serviceSpec.Envs, envs)
	if err != nil {
		return err
	}
	return nil
}

func BindToServiceBus(sourceType binding.SourceType, serviceSpec *ServiceSpec, serviceBus *AzureDepServiceBus) error {
	serviceSpec.AzureServiceBus = serviceBus
	envs, err := binding.GetBindingEnvs(binding.Source{Type: sourceType, IsSpringBootJms: serviceBus.IsJms},
		binding.Target{Type: binding.AzureServiceBus, AuthType: serviceBus.AuthType})
	if err != nil {
		return err
	}
	serviceSpec.Envs, err = binding.MergeMapWithDuplicationCheck(serviceSpec.Envs, envs)
	if err != nil {
		return err
	}
	return nil
}

func BindToEventHubs(sourceType binding.SourceType, serviceSpec *ServiceSpec, eventHubs *AzureDepEventHubs) error {
	serviceSpec.AzureEventHubs = eventHubs
	envs, err := binding.GetBindingEnvs(binding.Source{
		Type:              sourceType,
		IsSpringBootKafka: eventHubs.UseKafka,
		SpringBootVersion: eventHubs.SpringBootVersion,
	},
		binding.Target{Type: binding.AzureEventHubs, AuthType: eventHubs.AuthType})
	if err != nil {
		return err
	}
	serviceSpec.Envs, err = binding.MergeMapWithDuplicationCheck(serviceSpec.Envs, envs)
	if err != nil {
		return err
	}
	return nil
}

func BindToStorageAccount(sourceType binding.SourceType, serviceSpec *ServiceSpec,
	account *AzureDepStorageAccount) error {
	serviceSpec.AzureStorageAccount = account
	envs, err := binding.GetBindingEnvs(binding.Source{Type: sourceType},
		binding.Target{Type: binding.AzureStorageAccount, AuthType: account.AuthType})
	if err != nil {
		return err
	}
	serviceSpec.Envs, err = binding.MergeMapWithDuplicationCheck(serviceSpec.Envs, envs)
	if err != nil {
		return err
	}
	return nil
}

func BindToAIModels(sourceType binding.SourceType, serviceSpec *ServiceSpec, model string) error {
	serviceSpec.AIModels = append(serviceSpec.AIModels, AIModelReference{Name: model})
	envs, err := binding.GetBindingEnvs(binding.Source{Type: sourceType},
		binding.Target{Type: binding.AzureOpenAiModel, AuthType: internal.AuthTypeUnspecified})
	if err != nil {
		return err
	}
	serviceSpec.Envs, err = binding.MergeMapWithDuplicationCheck(serviceSpec.Envs, envs)
	if err != nil {
		return err
	}
	return nil
}

// BindToContainerApp a call b
// todo:
//  1. Add field in ServiceSpec to identify b's app type like Eureka server and Config server.
//  2. Create GetServiceBindingEnvsForContainerApp
//  3. Merge GetServiceBindingEnvsForEurekaServer and GetServiceBindingEnvsForConfigServer into
//     GetServiceBindingEnvsForContainerApp.
//  4. Delete printHintsAboutUseHostContainerApp use GetServiceBindingEnvsForContainerApp instead
func BindToContainerApp(a *ServiceSpec, b *ServiceSpec) {
	if a.Frontend == nil {
		a.Frontend = &Frontend{}
	}
	a.Frontend.Backends = append(a.Frontend.Backends, ServiceReference{Name: b.Name})
	if b.Backend == nil {
		b.Backend = &Backend{}
	}
	b.Backend.Frontends = append(b.Backend.Frontends, ServiceReference{Name: a.Name})
}
