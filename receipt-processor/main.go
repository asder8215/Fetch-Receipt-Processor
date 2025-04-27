package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type receipt struct {
    ID           string `json:"id"`
    Retailer     string `json:"retailer"`
    PurchaseDate string `json:"purchaseDate"`
    PurchaseTime string `json:"purchaseTime"`
    Items        []item `json:"items"`
    Total        string `json:"total"`
}

type item struct {
    ShortDescription string `json:"shortDescription"`
    Price            string `json:"price"`
}

func main() {
    router := gin.Default()
    router.GET("/receipts/:id/points", getReceiptPointsByID)
    router.POST("/receipts/process", postReceipt)

    router.Run("localhost:8080")
}


func getReceiptPointsByID(c *gin.Context) {
    id := c.Param("id")
}
