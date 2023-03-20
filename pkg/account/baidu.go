package account

import (
	ia "accountbalance/internal/account"
	"fmt"

	bb "github.com/baidubce/bce-sdk-go/services/billing"
	bba "github.com/baidubce/bce-sdk-go/services/billing/api"
	"golang.org/x/sync/errgroup"
)

type BaiduClient struct {
	*bb.Client
}

func NewBaiduClient(ak string, sk string, region string) (ia.BalanceClient, error) {
	client, err := bb.NewClient(ak, sk, region)
	if err != nil {
		return nil, err
	}
	return &BaiduClient{client}, nil
}

func (bc *BaiduClient) QueryBalance() (string, error) {
	response, err := bc.Client.GetBalance()
	if err != nil {
		return "", err
	}
	balanceStr := fmt.Sprintf("%.2f", response.CashBalance)
	return balanceStr, nil
}

func (bc *BaiduClient) QueryExpense(billCycle string, billDate ...string) (float32, error) {
	var startTime, endTime string
	req := bba.BillingParams{}
	if len(billDate) != 0 {
		startTime = billDate[0]
		endTime = billDate[0]
		billCycle = billDate[0][:7]
	}
	req.Month = billCycle
	req.BeginTime = startTime
	req.EndTime = endTime
	orderExpenseChan := make(chan float64, 30)

	productTypes := []string{"prepay", "postpay"}
	for _, pt := range productTypes {
		req.ProductType = pt
		err := bc.queryOneTypeExpense(req, orderExpenseChan)
		if err != nil {
			return -1, err
		}
	}
	close(orderExpenseChan)
	sum := 0.0
	for v := range orderExpenseChan {
		// fmt.Println("channel: ", v)
		sum += v
	}
	// fmt.Println("channel sum: ", sum)

	return float32(sum), nil
}

func (bc *BaiduClient) queryOneTypeExpense(query bba.BillingParams, inChan chan<- float64) error {
	var eGroup = new(errgroup.Group)
	response, err := bc.Client.GetBilling(&query)
	if err != nil {
		return err
	}
	recordNums := response.TotalCount
	pageSize := 50
	pageNums := recordNums/pageSize + 1
	query.PageSize = pageSize

	for i := 1; i <= pageNums; i++ {
		m := i
		eGroup.Go(func() error {
			newBillingParams := query
			newBillingParams.PageNo = m
			response, err := bc.Client.GetBilling(&newBillingParams)
			if err != nil {
				return err
			}
			for _, v := range response.Bills {
				inChan <- v.Cash
			}
			return nil
		})
	}
	err = eGroup.Wait()

	if err != nil {
		return err
	} else {
		return nil
	}
}
