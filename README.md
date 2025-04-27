# Fetch-Receipt-Processor
This is the code for Fetch Receipt Processor Challenge, which is written in Go. It was pretty fun to do and a good learning experience I'll say for sure.

## To Run the Code (Locally)
In your terminal:
1. Clone the repository:
   ```bash
   git clone git@github.com:asder8215/Fetch-Receipt-Processor.git
   ```
2. Switch to the receipt-processor directory via
   ```bash
   cd Fetch-Receipt-Processor/receipt-processor
   ```
3. Ensure that you have all the go dependencies using:
   ```bash
   go get .
   ```
4. Run the code
   ```bash
   go run .
   ```
5. From a new command line window (command + T on Mac), you can use curl to make a request to the running local service.
   - From the `examples` directory, you can perform a POST request to `/receipts/process` as follows:
     ```bash
     curl http://localhost:8080/receipts/process \
      --include \
      --header "Content-Type: application/json" \
      --request "POST" \
      --data @[insert_path_to_json_file]
   ```
   Where [insert_path_to_json_file] could be a path to a .json file containing a receipt information or you could directly create the receipt json instead of doing @[insert_path_to_json_file]
   - Once you have performed a POST request, the response of the POST request should give you a `<unique ID>`. You can use that `<unique ID>` to get how much points you are
   awarded for a specific receipt through the `/receipts/{id}/points` endpoint. You can perform the curl request as follows:
   ```bash
   curl http://localhost:8080/receipts/<unique_id>/points
   ```
That's pretty much it!

## Receipt Structure in Code
The following just shows how the `receipt` structure looks like in the code:
```go
type receipt struct {
    ID           string `json:"id"`
    Retailer     string `json:"retailer"`
    PurchaseDate string `json:"purchaseDate"`
    PurchaseTime string `json:"purchaseTime"`
    Items        []item `json:"items"`
    Total        string `json:"total"`
}
```
The `item` struct is defined as follows:
```go
type item struct {
    ShortDescription string `json:"shortDescription"`
    Price            string `json:"price"`
}
```
The `points` struct is returned within the GET request on the `/receipts/{id}/points` endpoint:
```go
type points struct {
    Points int64 `json:"points"`
}
```

An example of what a JSON file that a Receipt structure will take is as follows:
```json
{
    "retailer": "GameStop",
    "purchaseDate": "2021-01-28",
    "purchaseTime": "12:15",
    "total": "146.95",
    "items": [
        {"shortDescription": "Pokemon Sword", "price": "65.31"},
        {"shortDescription": "Legend of Zelda: Breath of the Wild", "price": "65.31"},
        {"shortDescription": "Hollow Knight", "price": "16.33"}
    ]
}
```

## Features Completed
- Implemented `/receipts/process` and `/receipts/{id}/points` as per instructed in [Fetch Receipt Processor Challenge](https://github.com/fetch-rewards/receipt-processor-challenge/tree/main)
  - Verified that in the POST method, that empty retailer name and item name are not possible, total and item prices are valid formats and are equal to each other, purchase date and time are valid formats
  - Points are awarded as specified in [Fetch Receipt Processor Challenge](https://github.com/fetch-rewards/receipt-processor-challenge/tree/main)

## References
- https://go.dev/doc/tutorial/web-service-gin ~ How to write RESTful API in Go (since I never done this in Go before; mostly Distributed System stuff)
- https://pkg.go.dev/time ~ Documentation on Time object in Go, which made it easier for me to parse
- https://blog.stackademic.com/unique-identifier-id-and-uuid-in-go-lang-99e6cc1b73b5 ~ How to generate a unique ID for the receipts
- Most other stuff like syntax, other packages used like (strconv, fmt, strings, unicode, etc.) I looked back on old code I've written in Go from past courses for reference
