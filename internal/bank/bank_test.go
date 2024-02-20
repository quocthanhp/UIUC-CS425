package bank_test

import (
	"mp1_node/internal/bank"
	"os"
	"testing"
)

var myBank *bank.Bank

func TestMain(m *testing.M) {
	myBank = bank.NewBank()
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestAddAccounts(t *testing.T) {
	myBank.AddAccounts("wqkby", "yxpqg")
	result := myBank.GetSize()
	expected := 2

	if result != expected {
		t.Errorf("Expected: %d, but got: %d", expected, result)
	}
}

func TestDeposit(t *testing.T) {
	myBank.Deposit("wqkby", 10)
	result, _ := myBank.GetBalance("wqkby")
	expected := uint(10)

	if result != expected {
		t.Errorf("Expected: %d, but got: %d", expected, result)
	}
}

func TestTransfer(t *testing.T) {
	myBank.Deposit("yxpqg", 75)
	myBank.Transfer("yxpqg", "hreqp", 13)

	result1, _ := myBank.GetBalance("hreqp")
	expected1 := uint(13)

	if result1 != expected1 {
		t.Errorf("Expected: %d, but got: %d", expected1, result1)
	}

	result2, _ := myBank.GetBalance("yxpqg")
	expected2 := uint(62)

	if result2 != expected2 {
		t.Errorf("Expected: %d, but got: %d", expected2, result2)
	}
}

func TestInvalidTransfer(t *testing.T) {
	_, err := myBank.Transfer("yxpqg", "wqkby", 70)

	if err == nil {
		t.Error("Invalid transfer did not return an error as expected")
	}
}

func TestPrintBalance(t *testing.T) {
	myBank.PrintBalances()
}
