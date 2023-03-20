package account

import (
	"accountbalance/pkg/setting"
)

type BalanceClient interface {
	QueryBalance() (string, error)
	QueryExpense(string, ...string) (float32, error)
}
type CreateBalanceClient func(ak string, sk string, region string) (BalanceClient, error)
type baseAccount struct {
	*setting.AccountInfo
}

func (a *baseAccount) GetInfo() (*string, *string, *string) {
	return &a.AccountInfo.AK, &a.AccountInfo.SK, &a.AccountInfo.Region
}

func newBaseAccount(a *setting.AccountInfo) *baseAccount {
	return &baseAccount{AccountInfo: a}
}
