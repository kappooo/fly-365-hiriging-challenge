package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type FlyTypeA struct {
	Amount         int64  `json:"amount"`
	Currency       string `json:"currency"`
	StatusCode     int    `json:"statusCode"`
	OrderReference string `json:"orderReference"`
	TransactionId  string `json:"transactionId"`
	FlyType        string
	Url            string
	Status         string
}

type FlyTypeB struct {
	Amount         int64  `json:"value"`
	Currency       string `json:"transactionCurrency"`
	StatusCode     int    `json:"statusCode"`
	OrderReference string `json:"orderInfo"`
	TransactionId  string `json:"paymentId"`
	FlyType        string
	Url            string
	Status         string
}

type FlyType struct {
	Amount         int64
	Currency       string
	StatusCode     int
	OrderReference string
	TransactionId  string
	FlyType        string
	Url            string
	Status         string
}

type FlyAPaymentData struct {
	PaymentData []FlyTypeA `json:"transactions"`
}
type FlyBPaymentData struct {
	PaymentData []FlyTypeB `json:"transactions"`
}
type FlyTypeGetter interface {
	getFlyData() []FlyType
}

func (fly FlyTypeA) getFlyData() []FlyType {
	jsonData, err := ioutil.ReadFile(fly.Url)
	if err != nil {
		fmt.Println("Error when opeing file")
	}
	var paymentData FlyAPaymentData
	var retData []FlyType
	err2 := json.Unmarshal(jsonData, &paymentData)

	if err2 != nil {
		fmt.Println("Error happend when parsing json", err2)
	}
	for i := 0; i < len(paymentData.PaymentData); i++ {
		paymentData.PaymentData[i].FlyType = fly.FlyType
		switch paymentData.PaymentData[i].StatusCode {
		case 1:
			paymentData.PaymentData[i].Status = "authorised"
		case 2:
			paymentData.PaymentData[i].Status = "decline"
		case 3:
			paymentData.PaymentData[i].Status = "refunded"
		}
		retData = append(retData, FlyType{
			Amount:         paymentData.PaymentData[i].Amount,
			Currency:       paymentData.PaymentData[i].Currency,
			StatusCode:     paymentData.PaymentData[i].StatusCode,
			OrderReference: paymentData.PaymentData[i].OrderReference,
			TransactionId:  paymentData.PaymentData[i].TransactionId,
			Status:         paymentData.PaymentData[i].Status,
			FlyType:        paymentData.PaymentData[i].FlyType,
		})
	}
	return retData
}
func (fly FlyTypeB) getFlyData() []FlyType {
	jsonData, err := ioutil.ReadFile(fly.Url)
	if err != nil {
		fmt.Println("Error when opeing file")
	}
	var paymentData FlyBPaymentData
	var retData []FlyType
	err2 := json.Unmarshal(jsonData, &paymentData)

	if err2 != nil {
		fmt.Println("Error happend when parsing json", err2)
	}
	for i := 0; i < len(paymentData.PaymentData); i++ {
		paymentData.PaymentData[i].FlyType = fly.FlyType

		switch paymentData.PaymentData[i].StatusCode {
		case 100:
			paymentData.PaymentData[i].Status = "authorised"
		case 200:
			paymentData.PaymentData[i].Status = "decline"
		case 300:
			paymentData.PaymentData[i].Status = "refunded"
		}
		retData = append(retData, FlyType{
			Amount:         paymentData.PaymentData[i].Amount,
			Currency:       paymentData.PaymentData[i].Currency,
			StatusCode:     paymentData.PaymentData[i].StatusCode,
			OrderReference: paymentData.PaymentData[i].OrderReference,
			TransactionId:  paymentData.PaymentData[i].TransactionId,
			Status:         paymentData.PaymentData[i].Status,
			FlyType:        paymentData.PaymentData[i].FlyType,
		})
	}
	return retData
}

func main() {

	http.HandleFunc("/api/payment/transaction", sayHello)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
func sayHello(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	data := getData()
	u, _ := url.Parse(r.URL.String())

	filterData := applyFilter(data, u.Query())
	js, _ := json.Marshal(filterData)
	w.Write(js)

}

func getData() []FlyType {
	var flyTypeA FlyTypeGetter
	t, _ := getFlyTypeFactory("flytypeA").(FlyTypeGetter)
	flyTypeA = t
	var flyTypeB FlyTypeGetter
	t2, _ := getFlyTypeFactory("flytypeB").(FlyTypeGetter)
	flyTypeB = t2
	return mergeArrayData(flyTypeA.getFlyData(), flyTypeB.getFlyData())
}

func mergeArrayData(types ...[]FlyType) []FlyType {
	data := make([]FlyType, 0)
	for _, typeItem := range types {
		for _, value := range typeItem {
			data = append(data, value)
		}
	}
	return data
}

func applyFilter(data []FlyType, filters map[string][]string) []FlyType {
	validateFilter := false
	retData := make([]FlyType, 0)
	for _, value := range data {
		validateFilter = false
		for key, filter_value := range filters {
			if (key == "provide") && (filter_value[0] != value.FlyType) {
				validateFilter = true
			}
			if (key == "statusCode") && (filter_value[0] != value.Status) {
				validateFilter = true
			}
			if (key == "currency") && (filter_value[0] != value.Currency) {
				validateFilter = true
			}
			if key == "amountMin" {
				i2, _ := strconv.ParseInt(filter_value[0], 10, 64)
				if i2 > value.Amount {
					validateFilter = true
				}
			}
			if key == "amountMax" {
				i2, _ := strconv.ParseInt(filter_value[0], 10, 64)
				if i2 < value.Amount {
					validateFilter = true
				}
			}
		}
		if !validateFilter {
			retData = append(retData, value)
		}
	}
	return retData
}

func getFlyTypeFactory(TypeName string) interface{} {

	switch TypeName {
	case "flytypeA":
		return FlyTypeA{
			Url:     "FlyPayA.json",
			FlyType: "flypayA",
		}
	case "flytypeB":
		return FlyTypeB{
			Url:     "FlyPayB.json",
			FlyType: "flypayB",
		}
	}
	return FlyType{}

}
