package handler

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/samipkalyani/posts_receipts_challenge/models"

	"github.com/gin-gonic/gin"
)

type PostReceiptResponse struct {
	Id string `json:"id"`
}

var regAlphaNum = regexp.MustCompile("[[:alnum:]]+")

func PostReceipts(ctx *gin.Context) {
	var newReceipt models.Receipt
	if err := ctx.BindJSON(&newReceipt); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var points int

	//1.	Add length of retailer
	matches := regAlphaNum.FindAllString(newReceipt.Retailer, -1)
	alphaNumRetailer := ""
	for _, match := range matches {
		alphaNumRetailer += match
	}
	points += len(alphaNumRetailer)

	//2.	50 points if price like .00
	total, parseTotalError := strconv.ParseFloat(newReceipt.Total, 64)
	if parseTotalError != nil {
		fmt.Println(parseTotalError)
		return
	}
	if total == float64(int64(total)) {
		points += 50
	}

	//3.	25 points if price multiple of 0.25
	if math.Mod(total, 0.25) == 0 {
		points += 25
	}

	//4.	5 points on every 2 items on the receipt.
	points += (len(newReceipt.Items) / 2) * 5

	//5.	If length of the item is multiple of 3.
	for _, item := range newReceipt.Items {
		if itemLength := len(strings.TrimSpace(item.ShortDescription)); itemLength%3 == 0 {
			itemPrice, parseItemPriceError := strconv.ParseFloat(item.Price, 64)
			if parseItemPriceError != nil {
				fmt.Println(parseItemPriceError)
				return
			}
			points += int(math.Ceil(itemPrice * 0.2))
		}
	}

	//6.	Odd day
	purchaseDate, parseError := time.Parse("2006-01-02", newReceipt.PurchaseDate)
	if parseError != nil {
		fmt.Println(parseError)
		return
	}
	if purchaseDate.Day()%2 != 0 {
		points += 6
	}

	//7.	Purchase B/W 2:00pm and 4:00pm
	parsedPurchaseTime, purchaseTimeParseError := time.Parse("15:04", newReceipt.PurchaseTime)
	if purchaseTimeParseError != nil {
		fmt.Println(purchaseTimeParseError)
		return
	}
	now := time.Now()
	purchaseTime := time.Date(now.Year(), now.Month(), now.Day(), parsedPurchaseTime.Hour(), parsedPurchaseTime.Minute(), 0, 0, time.Local)
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 14, 0, 0, 0, time.Local)
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 16, 0, 0, 0, time.Local)
	if purchaseTime.After(startTime) && purchaseTime.Before(endTime) {
		points += 10
	}

	// Set the random seed
	rand.Seed(time.Now().UnixNano())

	key := generateKey()

	PointsMap[key] = points

	ctx.JSON(http.StatusOK, PostReceiptResponse{
		Id: key,
	})
}

func generateKey() string {
	charPool := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-"
	pattern := "^\\S+$"
	regex := regexp.MustCompile(pattern)

	for {
		randomBytes := make([]byte, 36)

		// Making a string of 36 length
		for i := 0; i < 36; i++ {
			randomBytes[i] = charPool[rand.Intn(len(charPool))]
		}
		randomKey := string(randomBytes)

		if regex.MatchString(randomKey) {
			return randomKey
		}
	}
}
