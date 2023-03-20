package internal

import (
	ia "accountbalance/internal/account"
	"accountbalance/internal/render"
	"accountbalance/internal/sender"
	pa "accountbalance/pkg/account"
	"accountbalance/pkg/setting"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SpecificInfo struct {
	account           setting.AccountName
	project           string
	provider          string
	balance           string
	expenseLastCycle  string
	expenseLastDay    string
	expenseLastTwoDay string
	expenseDayRate    string
	balanceKeepMonth  string
	wLock             sync.RWMutex
}

var wg sync.WaitGroup
var clientSet = map[string]ia.CreateBalanceClient{
	"aliyun":  pa.NewAliyunClient,
	"tencent": pa.NewTencentClient,
	"baidu":   pa.NewBaiduClient,
}

var dataChan = make(chan render.InMsg, 10)

func Balance(s *setting.Config) {
	lastMonth := time.Now().AddDate(0, -1, 0).Format("2006-01")
	lastDay := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	lastTwoDay := time.Now().AddDate(0, 0, -2).Format("2006-01-02")
	for account, info := range s.Accounts {
		wg.Add(1)
		go func(info setting.AccountInfo, account setting.AccountName) {
			defer wg.Done()

			var singleWg sync.WaitGroup
			newInfo := &SpecificInfo{
				account:  account,
				project:  info.Project,
				provider: info.Provider,
			}

			cli, err := clientSet[info.Provider](info.AK, info.SK, info.Region)
			if err != nil {
				return
			}

			// query balance
			singleWg.Add(1)
			go func(cli ia.BalanceClient, sInfo *SpecificInfo) {
				defer singleWg.Done()
				balanceStr, err := cli.QueryBalance()
				if err != nil {
					log.Println(err)
					return
				}
				sInfo.wLock.Lock()
				defer sInfo.wLock.Unlock()
				sInfo.balance = balanceStr
			}(cli, newInfo)

			// query expense of last month
			singleWg.Add(1)
			go func(cli ia.BalanceClient, sInfo *SpecificInfo) {
				defer singleWg.Done()
				expenseLastMonth, err := cli.QueryExpense(lastMonth)
				if err != nil {
					log.Println(err)
					return
				}
				sInfo.wLock.Lock()
				defer sInfo.wLock.Unlock()
				sInfo.expenseLastCycle = fmt.Sprintf("%.2f", expenseLastMonth)
			}(cli, newInfo)

			//query expense of yesterday
			singleWg.Add(1)
			go func(cli ia.BalanceClient, sInfo *SpecificInfo) {
				defer singleWg.Done()
				expenseLastDay, err := cli.QueryExpense(lastDay, lastDay)
				if err != nil {
					log.Println(err)
					return
				}
				sInfo.wLock.Lock()
				defer sInfo.wLock.Unlock()
				sInfo.expenseLastDay = fmt.Sprintf("%.2f", expenseLastDay)
			}(cli, newInfo)

			// query expense of the day before yesterday
			singleWg.Add(1)
			go func(cli ia.BalanceClient, sInfo *SpecificInfo) {
				defer singleWg.Done()
				expenseLastTwoDay, err := cli.QueryExpense(lastTwoDay, lastTwoDay)
				if err != nil {
					log.Println(err)
					return
				}
				sInfo.wLock.Lock()
				defer sInfo.wLock.Unlock()
				sInfo.expenseLastTwoDay = fmt.Sprintf("%.2f", expenseLastTwoDay)
			}(cli, newInfo)

			singleWg.Wait()
			newInfo.SetBalanceKeepMonth()
			newInfo.SetExpenseDayRate()

			log.Println("sending")
			dataChan <- newInfo
		}(info, account)
	}
	wg.Wait()
	close(dataChan)
}

func (s *SpecificInfo) SetExpenseDayRate() error {
	if len(strings.TrimSpace(s.expenseLastDay)) == 0 || len(strings.TrimSpace(s.expenseLastTwoDay)) == 0 {
		return errors.New("Calculating expense day rate: invalid lastDayExpense or lastTwoDayExpense last cycle")
	}

	lastDayFloat, err := strconv.ParseFloat(s.expenseLastDay, 32)
	if err != nil {
		log.Println(err)
		return err
	}

	lastTwoDayFloat, err := strconv.ParseFloat(s.expenseLastTwoDay, 32)
	if err != nil {
		log.Println(err)
		return err
	}
	if lastTwoDayFloat <= 0.0 || lastDayFloat <= 0.0 {
		s.expenseDayRate = "0.00%"
		return nil
	}
	s.expenseDayRate = fmt.Sprintf("%.2f%%", (float32(lastDayFloat-lastTwoDayFloat))/float32(lastTwoDayFloat)*100)
	return nil

}
func (s *SpecificInfo) SetBalanceKeepMonth() error {
	ParseComma := func(in string) string {
		var out []rune
		for _, char := range in {
			// 是否为","
			if char == 0x2C {
				continue
			}
			out = append(out, char)
		}
		return string(out)
	}

	if len(strings.TrimSpace(s.balance)) == 0 || len(strings.TrimSpace(s.expenseLastCycle)) == 0 {
		return errors.New("Calculating balance keep month: invalid balance or expense last cycle")
	}

	balanceFloat, err := strconv.ParseFloat(ParseComma(s.balance), 32)
	if err != nil {
		log.Println(err)
		return err
	}

	expenseLastCycleFloat, err := strconv.ParseFloat(ParseComma(s.expenseLastCycle), 32)
	if err != nil {
		log.Println(err)
		return err
	}
	if expenseLastCycleFloat == 0 {
		return nil
	}

	s.balanceKeepMonth = fmt.Sprintf("%.2f", float32(balanceFloat)/float32(expenseLastCycleFloat))
	return nil
}

func (s *SpecificInfo) GetInfo() []string {
	return []string{string(s.account), s.project, s.provider, s.balance, s.expenseLastCycle, s.balanceKeepMonth, s.expenseLastDay, s.expenseLastTwoDay, s.expenseDayRate}
}

func Send(s *setting.Config) {
	sender.SendMail(s, render.RenderHtml(render.FormatRows(dataChan)))
}
