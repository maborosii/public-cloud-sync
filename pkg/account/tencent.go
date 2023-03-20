package account

import (
	ia "accountbalance/internal/account"
	"fmt"
	"math"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	tb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/billing/v20180709"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

var timeOut = 20

type TencentClient struct {
	*tb.Client
}

func NewTencentClient(ak string, sk string, region string) (ia.BalanceClient, error) {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqTimeout = timeOut
	credential := common.NewCredential(ak, sk)
	client, err := tb.NewClient(credential, region, cpf)
	if err != nil {
		return nil, err
	}
	return &TencentClient{client}, nil
}

func (tc *TencentClient) QueryBalance() (string, error) {
	// var respStr = []byte{}
	req := tb.NewDescribeAccountBalanceRequest()

	response, err := tc.Client.DescribeAccountBalance(req)
	if err != nil {
		return "", err
	}
	// err = response.ParseErrorFromHTTPResponse(respStr)
	// if err != nil {
	// 	return "", fmt.Errorf("查询失败,%v", err)
	// }
	balanceStr := fmt.Sprintf("%.2f", float64(*response.Response.Balance)/100.0)

	return balanceStr, nil
}

func (tc *TencentClient) QueryExpense(billCycle string, billDate ...string) (float32, error) {
	var startTime, endTime string
	req := tb.NewDescribeBillListRequest()
	if len(billDate) == 0 {
		startTime = billCycle + "-01 00:00:00"
		endTime = utilFormatTimeForEndOfMonth(startTime)
	} else {
		startTime = billDate[0] + " 00:00:00"
		endTime = utilFormatTimeForEndOfDay(startTime)
	}

	req.StartTime = tea.String(startTime)
	req.EndTime = tea.String(endTime)
	req.Limit = tea.Uint64(1)
	req.Offset = tea.Uint64(0)
	response, err := tc.Client.DescribeBillList(req)
	if err != nil {
		return -1, err
	}
	expense := *response.Response.DeductAmount / 100

	return float32(math.Abs(expense)), nil
}

func utilFormatTimeForEndOfMonth(billCycle string) string {
	formatTime, _ := time.Parse("2006-01-02 15:04:05", billCycle)
	endTime := formatTime.AddDate(0, 1, 0).Format("2006-01-02 15:04:05")
	return endTime
}
func utilFormatTimeForEndOfDay(billDate string) string {
	formatTime, _ := time.Parse("2006-01-02 15:04:05", billDate)
	endTime := formatTime.AddDate(0, 0, 1).Format("2006-01-02 15:04:05")
	return endTime
}
