package gads

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	gads "netto/controllers/gads/AcessToken"
	gads2 "netto/controllers/gads/fetch"
	"netto/models"
	rediscache "netto/redisCache"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var mu sync.Mutex
var wg sync.WaitGroup

var AccountList = []map[string]string{
	{"Nome da conta": "000-000-0000"}, //nome da conta, da sua preferencia | id da conta é o id da conta google ads

}

func FetchGoogle(c *gin.Context) {
	/* Vamos Utilizar o redis (https://github.com/redis/go-redis)
	Verificando antes se a busca ja existe no cache, se sim, não buscamos na API e pegamos direto do cache
	1) economizamos tempo na busca
	2) reduzimos o uso da maquina
	3) prevenção de request limit do google em cima do IP
	*/
	ctx := context.Background()
	accountFromCache, status, err := rediscache.GetAccountsOnCache(ctx, len(AccountList))
	if !status || err != nil {
		//Não encontramos a busca no cache,

		data := make(map[string]map[string]interface{})

		acessToken, err := gads.GetAcessToken()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"erro": err.Error(),
			})
			return
		}
		actualMonth := int(time.Now().Month())
		actualDay := int(time.Now().Day())
		actualYear := int(time.Now().Year())
		var cacheList []models.GoogleAdwAccount
		wg.Add(len(AccountList))

		for _, item := range AccountList {
			go func(item map[string]string) {
				defer wg.Done()

				responseOptChanel := make(chan map[string]interface{})
				responseMetricsChannel := make(chan map[string]interface{})

				go func() {
					defer close(responseOptChanel)
					for _, v := range item {
						responseOpt, err := gads2.FetchOptimization(acessToken, v)
						if err != nil {
							fmt.Printf("Erro ao pegar opt para %v\n", v)
						}
						responseOptChanel <- responseOpt
					}
				}()

				go func() {
					defer close(responseMetricsChannel)
					for _, v := range item {
						responseMetrics, err := gads2.FetchAccountDatas(acessToken, v, actualMonth, actualDay, actualYear)
						if err != nil {
							fmt.Printf("Erro ao pegar metrics para %v\n", v)
						}
						responseMetricsChannel <- responseMetrics
					}
				}()

				combinedData := make(map[string]interface{})

				for k, v := range item {
					var cacheItem models.GoogleAdwAccount
					cacheItem.AccountName = k
					cacheItem.AccountId = v

					opt := <-responseOptChanel
					for key, value := range opt {
						optValue, err := value.(float64)
						if !err {
							optValue = 0.0
						}
						cacheItem.AccountOptimazionScore = optValue
						combinedData[key] = value
					}

					metrics := <-responseMetricsChannel
					var ListCacheCampaign []models.Campaign
					for key, value := range metrics {

						if key != "cost" {
							campaigns, ok := value.([]map[string]interface{})
							if !ok {
								fmt.Println("value não é lista de campanhas")
								continue
							}
							for _, item := range campaigns {
								var CacheCampaign models.Campaign
								jsonBytes, err := json.Marshal(item)
								if err != nil {
									fmt.Println("error1:", err.Error())
								}
								err = json.Unmarshal(jsonBytes, &CacheCampaign)
								if err != nil {
									fmt.Println("value:", value)
									fmt.Println("key:", key)
									fmt.Println("error2:", err.Error())
								}
								ListCacheCampaign = append(ListCacheCampaign, CacheCampaign)
							}

						}
						acc_cost, errorCost := metrics["cost"].(float64)
						if !errorCost {
							acc_cost = 0.0
						}
						cacheItem.AccountCost = acc_cost
						combinedData[key] = value
					}
					cacheItem.Campaigns = ListCacheCampaign
					mu.Lock()
					if data[v] == nil {
						data[v] = make(map[string]interface{})
					}

					data[v] = combinedData
					combinedData["Name"] = k
					cacheList = append(cacheList, cacheItem)
					mu.Unlock()
				}
			}(item)
		}

		wg.Wait()
		errCache := rediscache.SaveAccountsOnCache(cacheList, ctx)
		if errCache != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"erro": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, cacheList)
		return
	} else {
		c.JSON(http.StatusOK, accountFromCache)
	}
}
