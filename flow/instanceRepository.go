package flow

import (
	"github.com/go-redis/redis/v8"
	"time"
)

type FlowInstanceRepository interface {
	Init(url string,db int,expire time.Duration)
	saveInstance(instance *flowInstance)(error)
	getInstance(instanceID string)(*flowInstance)
}

type DefaultFlowInstanceRepository struct {
	client *redis.Client
	expire time.Duration
}

func (repo *DefaultFlowInstanceRepository)Init(url string,db int,expire time.Duration){
	repo.client=redis.NewClient(&redis.Options{
        Addr:     url,
        Password: "", // no password set
        DB:       db,  // use default DB
    })
	repo.expire=expire
}

func (repo *DefaultFlowInstanceRepository)saveInstance(instance *flowInstance)(error){
	
	return repo.client.Set(repo.client.Context(), instance.InstanceID, *instance, repo.expire).Err()
}

func (repo *DefaultFlowInstanceRepository)getInstance(instanceID string)(flowInstance,error){
	return repo.client.Get(repo.client.Context(), instanceID).Result()
}

