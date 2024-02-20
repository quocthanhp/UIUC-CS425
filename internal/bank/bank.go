package bank

import (
	"errors"
	"fmt"

	"github.com/emirpasic/gods/maps/treemap"
)

const (
	SUCCESS = 1
	FAILURE = 0
)

type Bank struct {
	Accounts treemap.Map
}

func NewBank() *Bank {
	bank := &Bank{}
	bank.Accounts = *treemap.NewWithStringComparator()
	return bank
}

func (bank *Bank) AddAccounts(accts ...string) {
	for _, acct := range accts {
		bank.Accounts.Put(acct, uint(0))
	}
}

func (bank *Bank) Transfer(payer string, payee string, amount uint) (int, error) {
	payerBal, err := bank.GetBalance(payer)
	if err != nil {
		return FAILURE, err
	}

	if payerBal < amount {
		return FAILURE, fmt.Errorf("insufficient balance. Available balance: %d", payerBal)
	}

	bank.Accounts.Put(payer, payerBal-amount)
	bank.Deposit(payee, amount)

	return SUCCESS, nil
}

func (bank *Bank) Deposit(acct string, amount uint) {
	bal, err := bank.GetBalance(acct)
	if err != nil {
		// account not exist, create new account
		bank.Accounts.Put(acct, amount)
		return
	}

	bank.Accounts.Put(acct, bal+amount)
}

func (bank *Bank) GetBalance(acct string) (uint, error) {
	bal, found := bank.Accounts.Get(acct)

	if !found {
		return FAILURE, errors.New("account does not exist")
	}

	balUint, ok := bal.(uint)
	if !ok {
		return FAILURE, errors.New("invalid balance type")
	}

	return balUint, nil
}

func (bank *Bank) GetSize() int {
	return bank.Accounts.Size()
}

func (bank *Bank) PrintBalances() {
	statement := "BALANCES"
	it := bank.Accounts.Iterator()

	for it.Next() {
		acct := it.Key().(string)
		bal := it.Value().(uint)

		if bal != 0 {
			statement += fmt.Sprintf(" %s:%d", acct, bal)
		}
	}

	fmt.Println(statement)
}
