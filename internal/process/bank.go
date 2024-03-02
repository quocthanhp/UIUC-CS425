package process

import (
	"fmt"
)

// Bank represents a bank with accounts and their balances
type Bank struct {
	balances map[string]int
}

// NewBank creates a new instance of a Bank
func NewBank() *Bank {
	return &Bank{balances: make(map[string]int)}
}

// Deposit adds amount to the account's balance
func (b *Bank) Deposit(account string, amount int) {
	b.balances[account] += amount
	b.printBalances()
}

// Transfer transfers amount from one account to another
func (b *Bank) Transfer(from, to string, amount int) bool {
	if balance, ok := b.balances[from]; !ok || balance < amount {
		return false
	}
	b.balances[from] -= amount
	b.balances[to] += amount
	b.printBalances()
	return true
}

// printBalances prints all account balances
func (b *Bank) printBalances() {
	fmt.Print("BALANCES")
	for account, balance := range b.balances {
		fmt.Printf(" %s:%d", account, balance)
	}
	fmt.Println()
}
