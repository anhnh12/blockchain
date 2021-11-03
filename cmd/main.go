package main

import (
	"fmt"
	"math/big"
	"net/http"
	"strings"

	Comm "github.com/anhnh12/market-event-logs/pkg/common"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
)

type MarketItemSold struct {
	Value  *big.Int
}

func main() {
	router := gin.Default()
	router.GET("total-volume", ReadTotalTradedVolume())
	router.Run(":1231")
}

func ReadTotalTradedVolume() gin.HandlerFunc {
	return func(c *gin.Context) {
		client, err := ethclient.Dial("https://rinkeby.infura.io/v3/1eaaa354da99418fbca47f657fbb7f88")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"message": err.Error(),
			})
			return
		}
		
		address := common.HexToAddress("0x0E2949198a0f3464b0dCA2c030BdB412252e4F0f")
		query := ethereum.FilterQuery{
			Addresses: []common.Address{address},
		}

		logs, err := client.FilterLogs(c, query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"message": err.Error(),
			})
			return
		}

		abi, err := abi.JSON(strings.NewReader(Comm.MARKET_ABI))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"message": err.Error(),
			})
			return
		}

		result := big.NewInt(0)
		for _, log := range logs {
			var item MarketItemSold
			err := abi.UnpackIntoInterface(&item, "MarketItemSold", log.Data)
			fmt.Println(item)
			if err != nil {
				continue
			}
			fmt.Println("not error: ", item)
			result = result.Add(result, item.Value)
		}

		c.JSON(http.StatusOK, gin.H{
			"result": result,
		})
	}
}