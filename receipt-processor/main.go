package main

import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
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

// Submits a receipt for processing.
// @return: Returns the ID assigned to the receipt.
func postReceipt() {
    var newReceipt receipt

    // Bind the received JSON to newReceipt
    if err := c.BindJSON(&newReceipt); err != nil {
        c.IndentedJSON(http.BadRequest,
                       gin.H{"description": "The receipt is invalid."})
        return
    }
    
    // Verifying if the value given to newReceipt is valid
    
    // cannot have a empty string for a retail place
    trimmed_retailer = string.TrimSpace(newReceipt.Retailer)
    if trimmed_retailer == "" {
        c.IndentedJSON(http.BadRequest,
                       gin.H{"description": "The receipt is invalid."})
        return
    }

    // https://pkg.go.dev/time in func Parse portion of the documentation
    // rely on time for validity of PurchaseDate and PurchaseTime
    parsedDate, date_err := time.Parse("2006-01-02", newReceipt.PurchaseDate)
    if date_err != nil { 
        c.IndentedJSON(http.BadRequest,
                       gin.H{"description": "The receipt is invalid."})
        return
    }

    parsedTime, time_err := time.Parse("15:04", newReceipt.PurchaseTime)
    if time_err != nil { 
        c.IndentedJSON(http.BadRequest,
                       gin.H{"description": "The receipt is invalid."})
        return
    }
    
    // verifying if total is formatted correctly
    if !verifyCost(newReceipt.Total) {
        c.IndentedJSON(http.BadRequest,
                       gin.H{"description": "The receipt is invalid."})
        return 
    }

    // need to verify if the total matches with the item cumulative prices
    // verify cost already checked if the total  string is good, so no 
    // need to double check
    float_total, _ := strconv.ParseFloat(newReceipt.Total, 64)
    
    var float_prices float64 := 0.00
    
    // verifying for each item in the receipt
    for _, i := range newReceipt.Items { 
        trimmed_str := strings.TrimSpace(i.ShortDescription)
        if trimmed_str == "" {
            c.IndentedJSON(http.BadRequest,
                           gin.H{"description": "The receipt is invalid."})
            return
        }

        if !verifyCost(i.Price) {
            c.IndentedJSON(http.BadRequest,
                           gin.H{"description": "The receipt is invalid."})
            return 
        } else {
            float_price, _ := strconv.ParseFloat(i.Price, 64) 
            float_prices += float_price
        }
    }
    
    if float_total != float_prices { 
        c.IndentedJSON(http.BadRequest,
                       gin.H{"description": "The receipt is invalid."})
        return
    }

    // assign a unique ID to the new receipt
    newReceipt.ID = uuid.New().String()
    
    // Add the new receipt to the receipts slice
    receipts = append(receipts, newReceipt)
    
    // Returns the ID as the JSON response
    c.IndentedJSON(http.StatusCreated, newReceipt.ID)
}

// The total and price of an object has to follow 
// at the basic level: "0.00" (at least 4 characters
// where the 3rd to last character must be a "."
// followed by two numeric value)
func verifyCost(p string) bool { 
    var check_digits := false
    var digits_after := 0
    for _, c := range p {
        if c == "." {
            check_digits = true
        } else if unicode.IsDigit(c) {
            if check_digits {
                digits_after += 1
            }
        } else {
            return false
        }
    }

    if digits_after != 2 {
        return false
    }

    return true
}

// Returns the points awarded for the receipt.
// @return: The number of points awarded.
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
    for _, i := range r.Items {
        trimmed_str := strings.TrimSpace(i.ShortDescription)
        if len(trimmed_str) % 3 == 0 {
            float_price, _ := strconv.ParseFloat(i.Price, 64)
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

// Quick method for checking if a char is within alphanumeric
// values ("a"-"z" or "A" - "Z" and "0" - "9")
func isAlphanumeric(c rune) bool {
    return unicode.IsLetter(c) || unicode.IsDigit(c)
}
