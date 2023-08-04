package command_server

import (
	"WebHome/src/crontab_job"
	"WebHome/src/database/dao"
	"WebHome/src/database/model"
	"WebHome/src/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type CurrencyConvertCommand struct {
	BaseCommand
}

var foreignAliasIndex = map[string]string{
	"人民币":         "CNY",
	"美元":          "USD",
	"欧元":          "EUR",
	"日元":          "100JPY",
	"英镑":          "GBP",
	"港币":          "HKD",
	"澳元":          "AUD",
	"新西兰元":        "NZD",
	"新加坡元":        "SGD",
	"新币":          "SGD",
	"瑞士法郎":        "CHF",
	"加元":          "CAD",
	"马来西亚令吉":      "MYR",
	"令吉":          "MYR",
	"俄罗斯卢布":       "RUB",
	"卢布":          "RUB",
	"韩元":          "KRW",
	"阿拉伯联合酋长国迪拉姆": "AED",
	"迪拉姆":         "AED",
	"土耳其里拉":       "TRY",
	"里拉":          "TRY",
	"墨西哥比索":       "MXN",
	"泰铢":          "THB",
	"JPY":         "100JPY",
}

var supportForeignISO = []string{
	"CNY",
	"USD",
	"EUR",
	"100JPY",
	"HKD",
	"GBP",
	"AUD",
	"NZD",
	"SGD",
	"CHF",
	"CAD",
	"MYR",
	"RUB",
	"ZAR",
	"KRW",
	"AED",
	"SAR",
	"HUF",
	"PLN",
	"DKK",
	"SEK",
	"NOK",
	"TRY",
	"MXN",
	"THB",
}

func (ccc *CurrencyConvertCommand) ParseCommand(stdin string) {
	rawParts := strings.Split(stdin, " ")
	parts := utils.RemoveElements(rawParts, "").([]string)
	ccc.Required = make(map[string]string)
	var (
		sourceCurrency string
		targetCurrency string
	)
	sourceValid := true
	targetValid := true
	for i := 0; i < len(parts); i++ {
		arg := parts[i]
		switch {
		case strings.ToUpper(arg) == "-T":
			if i+1 < len(parts) {
				i++
				targetCurrency, targetValid = validationSupportCurrency(strings.ToUpper(parts[i]))
				ccc.Required["targetCurrency"] = targetCurrency
			} else {
				targetValid = false
			}
			continue
		case strings.ToUpper(arg) == "-S":
			if i+1 < len(parts) {
				i++
				sourceCurrency, sourceValid = validationSupportCurrency(strings.ToUpper(parts[i]))
				ccc.Required["sourceCurrency"] = sourceCurrency
			} else {
				sourceValid = false
			}
			continue
		}
		if sourceValid == false || targetValid == false {
			return
		}
		numStr := strings.ReplaceAll(arg, ",", "")
		_, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			currency, valid := validationSupportCurrency(arg)
			if valid {
				if sourceCurrency == "" {
					sourceCurrency = currency
					ccc.Required["sourceCurrency"] = sourceCurrency
					continue
				} else if targetCurrency == "" {
					targetCurrency = currency
					ccc.Required["targetCurrency"] = targetCurrency
					continue
				}
			} else {
				return
			}
		}
		ccc.Required["amount"] += numStr
	}
	_, err := strconv.ParseFloat(ccc.Required["amount"], 64)
	if err != nil {
		ccc.Required["amount"] = ""
	}
	if sourceCurrency == "" {
		ccc.Required["sourceCurrency"] = "CNY"
	}
	if targetCurrency == "" {
		ccc.Required["targetCurrency"] = "USD"
	}
}

func (ccc *CurrencyConvertCommand) ExecuteCommand(c *gin.Context) {
	sourceCurrency := ccc.Required["sourceCurrency"]
	if sourceCurrency == "" {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Unsupported source currency exchange rate conversions."})
		return
	}

	targetCurrency := ccc.Required["targetCurrency"]
	if targetCurrency == "" {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Unsupported target currency exchange rate conversions."})
		return
	}

	amountStr := ccc.Required["amount"]
	if amountStr == "" {
		c.JSON(http.StatusOK, gin.H{"response": "ERROR: Amount that cannot be parsed"})
		return
	}
	rate := currencyConvert(sourceCurrency, targetCurrency)
	amount, _ := strconv.ParseFloat(amountStr, 64)
	if sourceCurrency == "100JPY" {
		sourceCurrency = "JPY"
		amount /= 100
	} else if targetCurrency == "100JPY" {
		targetCurrency = "JPY"
		amount *= 100
	}
	res := amountConversion(amount, rate)
	c.JSON(http.StatusOK, gin.H{"response": fmt.Sprintf("%s %s = %f %s", amountStr, sourceCurrency, res, targetCurrency)})
}

func validationSupportCurrency(currency string) (string, bool) {
	sort.Strings(supportForeignISO)
	currencyIndex := sort.SearchStrings(supportForeignISO, currency)
	aliasIndex := sort.SearchStrings(supportForeignISO, foreignAliasIndex[currency])
	if currencyIndex < len(supportForeignISO) && supportForeignISO[currencyIndex] == currency {
		return currency, true
	} else if aliasIndex < len(supportForeignISO) && supportForeignISO[aliasIndex] == foreignAliasIndex[currency] {
		return supportForeignISO[aliasIndex], true
	} else {
		return "", false
	}
}

func currencyConvert(sourceCurrency, targetCurrency string) float64 {
	exchangeRateDao := dao.NewExchangeRateDao()
	exchangeRateInfo := getExchangeRateInfo(sourceCurrency, targetCurrency, exchangeRateDao)
	var sourceRate, targetRate *float64
	var sourceIsDirectQuotation, targetIsDirectQuotation bool
	if exchangeRateInfo["sourceCurrency"] != nil {
		sourceRate = exchangeRateInfo["sourceCurrency"].Value
		sourceIsDirectQuotation = exchangeRateInfo["sourceCurrency"].IsDirectQuotation
	} else if sourceCurrency == "CNY" {
		var CNYRate = 1.0
		sourceRate = &CNYRate
		sourceIsDirectQuotation = true
	}
	if exchangeRateInfo["targetCurrency"] != nil {
		targetRate = exchangeRateInfo["targetCurrency"].Value
		targetIsDirectQuotation = exchangeRateInfo["targetCurrency"].IsDirectQuotation
	} else if targetCurrency == "CNY" {
		var CNYRate = 1.0
		targetRate = &CNYRate
		targetIsDirectQuotation = true
	}
	if sourceRate != nil && targetRate != nil {
		return calculateExchangeRate(sourceRate, targetRate, sourceIsDirectQuotation, targetIsDirectQuotation)
	} else if sourceRate == nil {
		return -1.0
	} else {
		return -2.0
	}
}

func getExchangeRateInfo(sourceCurrency, targetCurrency string, dao *dao.ExchangeRateDao) map[string]*model.ExchangeRate {
	updateExchangeRateData(dao)
	lastTradeDay := dao.GetLastTradingDay()
	if lastTradeDay != "" {
		var exchangeRateInfo []model.ExchangeRate
		if sourceCurrency == "CNY" {
			_ = dao.
				Where("trade_date=? AND foreign_currency=?",
					lastTradeDay, targetCurrency).
				Find(&exchangeRateInfo)
		} else if targetCurrency == "CNY" {
			_ = dao.
				Where("trade_date=? AND foreign_currency=?",
					lastTradeDay, sourceCurrency).
				Find(&exchangeRateInfo)
		} else {
			_ = dao.
				Where("trade_date=? AND (foreign_currency=? OR foreign_currency=?)",
					lastTradeDay, sourceCurrency, targetCurrency).
				Find(&exchangeRateInfo)
		}
		rateInfo := map[string]*model.ExchangeRate{
			"sourceCurrency": nil,
			"targetCurrency": nil,
		}
		for _, rate := range exchangeRateInfo {
			if rate.ForeignCurrency == sourceCurrency {
				rateCopy := rate
				rateInfo["sourceCurrency"] = &rateCopy
			} else {
				rateCopy := rate
				rateInfo["targetCurrency"] = &rateCopy
			}
		}
		return rateInfo
	} else {
		return nil
	}
}

func updateExchangeRateData(dao *dao.ExchangeRateDao) {
	latestUpdateTime := dao.GetTheLatestUpdateTime()
	location, _ := time.LoadLocation("Asia/Shanghai")
	updatedTime := time.UnixMilli(latestUpdateTime).In(location)
	todayTime := time.Now().In(location)
	if updatedTime.Before(time.Date(todayTime.Year(), todayTime.Month(), todayTime.Day(), 9, 15, 59, 0, location)) &&
		todayTime.After(time.Date(todayTime.Year(), todayTime.Month(), todayTime.Day(), 9, 15, 59, 0, location)) {
		crontab_job.ExchangeRateProcess()
		dao.FlushUpdateTime(todayTime.UnixMilli())
	}
}

func calculateExchangeRate(sourceRate, targetRate *float64, sourceIsDirectQuotation, targetIsDirectQuotation bool) float64 {
	var res float64
	if sourceIsDirectQuotation == true && targetIsDirectQuotation == true {
		res = *sourceRate / *targetRate
	} else if sourceIsDirectQuotation == false && targetIsDirectQuotation == false {
		res = 1 / *sourceRate * *targetRate
	} else if sourceIsDirectQuotation == true && targetIsDirectQuotation == false {
		res = *sourceRate * *targetRate
	} else {
		res = 1 / (*sourceRate * *targetRate)
	}
	return res
}

func amountConversion(amount float64, rate float64) float64 {
	return utils.LimitDecimalPlaces(amount*rate, 5)
}
