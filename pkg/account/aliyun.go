package account

import (
	"fmt"

	ia "accountbalance/internal/account"

	bssopenapi "github.com/alibabacloud-go/bssopenapi-20171214/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
)

var staticConfig = openapi.Config{
	// ms
	ReadTimeout:    tea.Int(20000),
	ConnectTimeout: tea.Int(10000),
}

type AliyunClient struct {
	*bssopenapi.Client
}

// func NewAliyunClient(ak string, sk string, region string) (*AliyunClient, error) {
func NewAliyunClient(ak string, sk string, region string) (ia.BalanceClient, error) {
	valueOfConfig := staticConfig
	config := &valueOfConfig
	config.AccessKeyId, config.AccessKeySecret, config.RegionId = &ak, &sk, &region
	result, err := bssopenapi.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &AliyunClient{result}, nil
}

func (ac *AliyunClient) QueryBalance() (string, error) {
	response, err := ac.Client.QueryAccountBalance()
	if err != nil {
		return "", err
	}

	if !*response.Body.Success {
		return "", fmt.Errorf("查询失败")
	}
	return *response.Body.Data.AvailableAmount, nil
}

func (ac *AliyunClient) QueryExpense(billCycle string, billDate ...string) (float32, error) {
	req := &bssopenapi.QueryAccountBillRequest{
		BillingCycle: &billCycle,
	}
	if len(billDate) != 0 {
		req.SetGranularity("DAILY").SetBillingDate(billDate[0])
	}

	response, err := ac.Client.QueryAccountBill(req)
	if err != nil {
		return -1, err
	}
	if !*response.Body.Success {
		return -1, fmt.Errorf("查询失败")
	}

	var expense float32 = 0
	for _, item := range response.Body.Data.Items.Item {
		expense += *item.PaymentAmount
	}
	return expense, nil
}
