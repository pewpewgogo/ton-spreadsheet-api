package handlers

import (
	"github.com/xssnick/tonutils-go/address"
	"log"
	"math/big"
	"net/http"
	"ton-balance-api/token"
	"ton-balance-api/ton"

	"github.com/gin-gonic/gin"
)

type BalanceRequest struct {
	Address string `json:"address"`
	Ticker  string `json:"ticker"`
}

type BalanceResponse struct {
	Address string   `json:"address"`
	Ticker  string   `json:"ticker"`
	Balance *big.Int `json:"balance"`
}

var (
	addressContractNOT  = "EQAvlWFDxGF2lXm67y4yzC17wYKD9A0guwPkMs1gOsM__NOT"
	addressContractDOGS = "EQCvxJy4eG8hyHBFsZ7eePxrRsUQSFE_jpptRAYBmcG_DOGS"
	addressContractUSDT = "EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"
)

var TicketConvert = map[string]*string{
	"NOT":  &addressContractNOT,
	"DOGS": &addressContractDOGS,
	"USDT": &addressContractUSDT,
	"TON":  nil,
}

func GetBalance(c *gin.Context) {
	var req BalanceRequest
	var ls, _ = ton.NewLsConnection(c)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var tokenContract = TicketConvert[req.Ticker]
	var userAddress, _ = address.ParseAddr(req.Address)

	var balance *big.Int
	var err error

	if tokenContract == nil {
		balance, err = token.GetTonBalance(c, ls, userAddress)
	} else {
		parsedTokenContract, _ := address.ParseAddr(*tokenContract)
		balance, err = token.GetJettonBalance(c, ls, parsedTokenContract, userAddress)
	}

	if err != nil {
		log.Println("Failed to fetch balance:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch balance"})
		return
	}

	c.JSON(http.StatusOK, BalanceResponse{
		Address: req.Address,
		Ticker:  req.Ticker,
		Balance: balance,
	})
}
