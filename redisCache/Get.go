package rediscache

import (
	"context"
	"encoding/json"
	"fmt"
	"netto/models"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func GetAccountsOnCache(ctx context.Context, totalLen int) ([]models.GoogleAdwAccount, bool, error) {
	if !VerifyLen(totalLen, ctx) {
		fmt.Println("Precisamos gerar uma nova")
		return nil, false, nil
	}
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

func VerifyLen(totalLen int, ctx context.Context) bool {
	redisConn := RedisConnection()

	total, err := redisConn.Get(ctx, "total_accounts").Result()

	if err == redis.Nil {
		_ = redisConn.Set(ctx, "total_accounts", totalLen, 0).Err()
		return false
	}

	if err != nil {
		return false
	}

	savedValue, err := strconv.Atoi(total)
	if err != nil {
		return false
	}

	if savedValue == totalLen {
		return true
	}

	_ = redisConn.Set(ctx, "total_accounts", totalLen, 0).Err()
	return false
}
