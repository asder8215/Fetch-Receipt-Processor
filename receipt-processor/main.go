package main

import (
    "github.com/gin-gonic/gin"
    "math"
    "net/http"
    "strconv"
    "strings"
    "time"
    "unicode"
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

type points struct {
    Points int64 `json:"points"`
}

var receipts = []receipt{}


func main() {
    router := gin.Default()
    router.GET("/receipts/:id/points", getReceiptPointsByID)
    router.POST("/receipts/process", postReceipt)

    router.Run("localhost:8080")
}

func postReceipt() {
    
}

func getReceiptPointsByID(c *gin.Context) {
    id := c.Param("id")
    
    for _, r := range receipts {
        if r.ID == id {
            json_points = points { processPoints(r) }
            c.IndentedJSON(http.StatusOK, json_points)
            return
        }
    }

    c.IndentedJSON(http.StatusNotFound, 
    gin.H{"description": "No receipt found for that ID."})
}

func processPoints(r receipt) {
    var total_points int64 := 0
    
    // One point for every alphanumeric character in the retailer name
    for _, c := range r.Retailer {
        if isAlphanumeric(c) {
            total_points += 1
        }
    }
    
    // ignoring error returned from strconv because
    // POST method should make sure receipt is valid
    float_total, _ := strconv.ParseFloat(r.Total, 64)

    rounded_float_total := int64(float_total * 100)
    
    // 50 points if the total is a round dollar amount with no cents
    if rounded_float_total % 100 == 0 {
        total_points += 50
    }

    // 25 points if the total is a multiple of 0.25
    if rounded_float_total % 25 == 0 {
        total_points += 25 
    }

    // 5 points for every two items on the receipt 
    total_points += (len(r.Items) / 2) * 5
    
    // If the trimmed length of the item description is a multiple of 3, 
    // multiply the price by 0.2 and round up to the nearest integer. 
    // The result is the number of points earned.
    for _, item := range r.Items {
        trimmed_str := strings.TrimSpace(item.ShortDescription)
        if len(trimmed_str) % 3 == 0 {
            float_price, _ := strconv.ParseFloat(item.Price, 64)
            ceiled_price := math.Ceil(float_price * 0.2)
            total_price += int64(ceiled_price)
        }
    }

    // 6 points if the *day* in the purchase date is odd
    // Purchase Date will always appear in the format of 
    // yyyy-mm-dd
    
    // converts the string into a specific structure
    parsedDate, _ := time.Parse("2006-01-02", r.PurchaseDate)
    // gives day as an integer
    day := parsedDate.Day()
    
    if day % 2 == 1 {
        total_points += 6
    }


    // 10 points if the time of purchase is after 2:00pm and before 4:00pm
    parsedTime, _ := time.Parse("15:04", r.PurchaseTime)
    hour := parsedDate.Hour()

    if hour > 2 && hour < 4 {
        total_points += 10
    }

    return total_points
}
