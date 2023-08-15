package crontab_job

import (
	"WebHome/src/database"
	"WebHome/src/database/dao"
	"WebHome/src/database/model"
	"WebHome/src/utils"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type DailyExchangeRate struct {
	tradeDate         string
	foreignCurrency   string
	value             *float64
	isDirectQuotation bool
}

type currencyBasicInfo struct {
	name              string
	isDirectQuotation bool
}

var (
	once         sync.Once
	db           *gorm.DB
	currencyList []currencyBasicInfo
	location     *time.Location
)

func init() {
	once.Do(func() {
		db, _ = database.ConnectDB()
	})
}

func ExchangeRateProcess() {
	initTime := determineTheInitialDate()
	currentTime := utils.GetCurrentTime()
	location, _ = time.LoadLocation("Asia/Shanghai")
	endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, location)
	for initTime.Before(endTime) && initTime.AddDate(1, 0, -1).Before(endTime) {
		startDate := initTime.Format("2006-01-02")
		initTime = initTime.AddDate(1, 0, -1)
		endDate := initTime.Format("2006-01-02")
		getAndSaveExchangeRateData(startDate, endDate)
		initTime = initTime.AddDate(0, 0, 1)
	}
	startDate := initTime.Format("2006-01-02")
	endDate := endTime.Format("2006-01-02")
	if startDate <= endDate {
		getAndSaveExchangeRateData(startDate, endDate)
	}
}

func determineTheInitialDate() (initTime time.Time) {
	exchangeRateDao := dao.NewExchangeRateDao()
	lastTradeDay := exchangeRateDao.GetLastTradingDay()
	if lastTradeDay == "" {
		location, _ = time.LoadLocation("Asia/Shanghai")
		initTime = time.Date(2006, 1, 4, 0, 0, 0, 0, location)
	} else {
		initTime, _ = time.Parse("2006-01-02", lastTradeDay)
		initTime = initTime.AddDate(0, 0, 1)
	}
	return initTime
}

func getAndSaveExchangeRateData(startDate, endDate string) {
	rawData := requestExchangeRateData(100, startDate, endDate)
	pageTotal := getPageTotal(rawData)
	resultCh := make(chan []model.ExchangeRate, pageTotal)
	semaphore := make(chan struct{}, 1)
	for pageNum := 1; pageNum <= pageTotal; pageNum++ {
		time.Sleep(345 * time.Millisecond)
		go func(pageNum int) {
			rawData := requestExchangeRateData(pageNum, startDate, endDate)
			resultCh <- parseExchangeRateData(rawData, semaphore)
		}(pageNum)
	}
	saveData(resultCh, pageTotal)
}

func requestExchangeRateData(pageNum int, startDate, endDate string) map[string]interface{} {
	url := fmt.Sprintf("https://www.chinamoney.com.cn/ags/ms/cm-u-bk-ccpr/CcprHisNew?startDate=%s&endDate=%s&pageNum=%d", startDate, endDate, pageNum)
	resp := utils.HttpGet(url)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)
	var rawData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&rawData); err != nil {
		log.Fatalf("Failed to decode JSON: %s", err)
	}
	return rawData
}

func getPageTotal(rawData map[string]interface{}) (pageTotal int) {
	if len(currencyList) == 0 {
		searchList := rawData["data"].(map[string]interface{})["searchlist"].([]interface{})
		parseCurrencyBasicInfo(searchList)
	}
	return int(rawData["data"].(map[string]interface{})["pageTotal"].(float64))
}

func parseExchangeRateData(rawData map[string]interface{}, semaphore chan struct{}) (exchangeRateData []model.ExchangeRate) {
	if len(currencyList) == 0 {
		searchList := rawData["data"].(map[string]interface{})["searchlist"].([]interface{})
		parseCurrencyBasicInfo(searchList)
	}
	for _, records := range rawData["records"].([]interface{}) {
		tradeDate := records.(map[string]interface{})["date"].(string)
		values := records.(map[string]interface{})["values"].([]interface{})
		for i := 0; i < len(values); i++ {
			semaphore <- struct{}{}
			singleData := model.ExchangeRate{BaseModel: *model.NewBaseModel()}
			singleData.TradeDate = tradeDate
			singleData.ForeignCurrency = currencyList[i].name
			value, err := strconv.ParseFloat(values[i].(string), 64)
			if err != nil {
				singleData.Value = nil
			} else {
				singleData.Value = &value
			}
			singleData.IsDirectQuotation = currencyList[i].isDirectQuotation
			exchangeRateData = append(exchangeRateData, singleData)
			<-semaphore
		}
	}
	return exchangeRateData
}

func parseCurrencyBasicInfo(searchList []interface{}) {
	for _, value := range searchList {
		var currency currencyBasicInfo
		arr := strings.Split(value.(string), "/")
		if arr[0] == "CNY" {
			currency.name = arr[1]
			currency.isDirectQuotation = false
		} else {
			currency.name = arr[0]
			currency.isDirectQuotation = true
		}
		currencyList = append(currencyList, currency)
	}
}

func saveData(resultCh chan []model.ExchangeRate, pageTotal int) {
	for i := 1; i <= pageTotal; i++ {
		res := <-resultCh
		db.Create(res)
	}
}
