package gads

import (
	"fmt"
	"net/http"
	gads "netto/controllers/gads/AcessToken"
	gads2 "netto/controllers/gads/fetch"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var mu sync.Mutex
var wg sync.WaitGroup

var AccountList = []map[string]string{
	{"Nome_da_conta": "000-000-0000"}, //nome da conta, da sua preferencia | id da conta é o id da conta google ads, sem os traços(-)
}

func FetchGoogle(c *gin.Context) {
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
				opt := <-responseOptChanel
				for key, value := range opt {
					combinedData[key] = value
				}

				metrics := <-responseMetricsChannel
				for key, value := range metrics {
					combinedData[key] = value
				}

				mu.Lock()
				if data[v] == nil {
					data[v] = make(map[string]interface{})
				}

				data[v] = combinedData
				combinedData["Name"] = k

				mu.Unlock()
			}
		}(item)
	}

	wg.Wait()

	c.JSON(http.StatusOK, data)
}
