package rediscache

import (
	"context"
	"encoding/json"
	"netto/models"

	"github.com/redis/go-redis/v9"
)

func GetAccountsOnCache(ctx context.Context) ([]models.GoogleAdwAccount, bool, error) {
	redisConn := RedisConnection()
	var accounts []models.GoogleAdwAccount
	cachedData, err := redisConn.Get(ctx, "adw_accounts").Result()
	if err == redis.Nil {
		//nao tem esse valor salvo
		return nil, false, nil
	}
	if err != nil {
		//algum erro real com o redis
		return nil, false, err
	}
	if err := json.Unmarshal([]byte(cachedData), &accounts); err != nil {
		//erro ao deserializar os dados
		return nil, false, err
	}
	//retorno transformado da lista de contas
	return accounts, true, nil
}
