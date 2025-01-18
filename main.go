package main

import (
	"fmt"
)

type User struct {
	ID   string
	Name string
}

func NewUser(id, name string) *User {
	return &User{ID: id, Name: name}
}

type Split interface {
	GetAmount() float64
	SetAmount(amount float64)
	GetUser() *User
}

type EqualSplit struct {
	User   *User
	Amount float64
}

func (e *EqualSplit) GetAmount() float64 {
	return e.Amount
}

func (e *EqualSplit) SetAmount(amount float64) {
	e.Amount = amount
}

func (e *EqualSplit) GetUser() *User {
	return e.User
}

type ExactSplit struct {
	User   *User
	Amount float64
}

func (e *ExactSplit) GetAmount() float64 {
	return e.Amount
}

func (e *ExactSplit) SetAmount(amount float64) {
	e.Amount = amount
}

func (e *ExactSplit) GetUser() *User {
	return e.User
}

type PercentSplit struct {
	User    *User
	Percent float64
	Amount  float64
}

func (p *PercentSplit) GetAmount() float64 {
	return p.Amount
}

func (p *PercentSplit) SetAmount(amount float64) {
	p.Amount = amount
}

func (p *PercentSplit) GetUser() *User {
	return p.User
}

type ExpenseMetadata struct {
	Name string
}

type Expense interface {
	Validate() bool
	GetAmount() float64
	GetSplits() []Split
	GetPaidBy() *User
}

type EqualExpense struct {
	Amount   float64
	PaidBy   *User
	Splits   []Split
	Metadata *ExpenseMetadata
}

func (e *EqualExpense) Validate() bool {
	for _, split := range e.Splits {
		if _, ok := split.(*EqualSplit); !ok {
			return false
		}
	}
	return true
}

func (e *EqualExpense) GetAmount() float64 {
	return e.Amount
}

func (e *EqualExpense) GetSplits() []Split {
	return e.Splits
}

func (e *EqualExpense) GetPaidBy() *User {
	return e.PaidBy
}

type ExactExpense struct {
	Amount   float64
	PaidBy   *User
	Splits   []Split
	Metadata *ExpenseMetadata
}

func (e *ExactExpense) Validate() bool {
	totalAmount := e.Amount
	sumSplitAmount := 0.0
	for _, split := range e.Splits {
		if exactSplit, ok := split.(*ExactSplit); ok {
			sumSplitAmount += exactSplit.Amount
		}
	}

	if totalAmount != sumSplitAmount {
		return false
	}
	return true
}

func (e *ExactExpense) GetAmount() float64 {
	return e.Amount
}

func (e *ExactExpense) GetSplits() []Split {
	return e.Splits
}

func (e *ExactExpense) GetPaidBy() *User {
	return e.PaidBy
}

type PercentExpense struct {
	Amount   float64
	PaidBy   *User
	Splits   []Split
	Metadata *ExpenseMetadata
}

func (e *PercentExpense) Validate() bool {
	totalPercent := 100.0
	sumSplitPercent := 0.0
	for _, split := range e.Splits {
		if percentSplit, ok := split.(*PercentSplit); ok {
			sumSplitPercent += percentSplit.Percent
		}
	}

	if totalPercent != sumSplitPercent {
		return false
	}
	return true
}

func (e *PercentExpense) GetAmount() float64 {
	return e.Amount
}

func (e *PercentExpense) GetSplits() []Split {
	return e.Splits
}

func (e *PercentExpense) GetPaidBy() *User {
	return e.PaidBy
}

type ExpenseService struct{}

func (es *ExpenseService) CreateExpense(expenseType string, amount float64, paidBy *User, splits []Split, metadata *ExpenseMetadata) Expense {
	switch expenseType {
	case "EQUAL":
		totalSplits := len(splits)
		splitAmount := amount / float64(totalSplits)
		for _, split := range splits {
			split.SetAmount(splitAmount)
		}
		return &EqualExpense{
			Amount:   amount,
			PaidBy:   paidBy,
			Splits:   splits,
			Metadata: metadata,
		}
	case "EXACT":
		totalAmount := amount
		sumSplitAmount := 0.0
		for _, split := range splits {
			if exactSplit, ok := split.(*ExactSplit); ok {
				sumSplitAmount += exactSplit.Amount
			}
		}
		if totalAmount != sumSplitAmount {
			return nil
		}
		return &ExactExpense{
			Amount:   amount,
			PaidBy:   paidBy,
			Splits:   splits,
			Metadata: metadata,
		}
	case "PERCENT":
		for _, split := range splits {
			if percentSplit, ok := split.(*PercentSplit); ok {
				split.SetAmount(amount * percentSplit.Percent / 100)
			}
		}
		return &PercentExpense{
			Amount:   amount,
			PaidBy:   paidBy,
			Splits:   splits,
			Metadata: metadata,
		}
	}
	return nil
}

type ExpenseManager struct {
	Expenses     []Expense
	UserMap      map[string]*User
	BalanceSheet map[string]map[string]float64
}

func NewExpenseManager() *ExpenseManager {
	return &ExpenseManager{
		Expenses:     []Expense{},
		UserMap:      make(map[string]*User),
		BalanceSheet: make(map[string]map[string]float64),
	}
}

func (em *ExpenseManager) AddUser(user *User) {
	em.UserMap[user.ID] = user
	em.BalanceSheet[user.ID] = make(map[string]float64)
}

func (em *ExpenseManager) AddExpense(expenseType string, amount float64, paidBy string, splits []Split, metadata *ExpenseMetadata) {
	expenseService := &ExpenseService{}
	expense := expenseService.CreateExpense(expenseType, amount, em.UserMap[paidBy], splits, metadata)
	if expense != nil {
		em.Expenses = append(em.Expenses, expense)

		// Update balances
		for _, split := range expense.GetSplits() {
			splitAmount := split.GetAmount()
			payer := expense.GetPaidBy().ID
			receiver := split.GetUser().ID

			em.BalanceSheet[payer][receiver] -= splitAmount
			em.BalanceSheet[receiver][payer] += splitAmount
		}

		// Print Summary for Each Event
		em.PrintEventSummary(metadata.Name)
	}
}

func (em *ExpenseManager) PrintEventSummary(eventName string) {
	fmt.Printf("\nSummary of Bill Settlement for %s:\n", eventName)
	printed := make(map[string]bool)
	for user, balances := range em.BalanceSheet {
		for otherUser, amount := range balances {
			if amount != 0 && !printed[user+otherUser] {
				em.PrintBalance(user, otherUser, amount)
				printed[user+otherUser] = true
				printed[otherUser+user] = true
			}
		}
	}
}

func (em *ExpenseManager) ShowBalance(userId string) {
	balances := em.BalanceSheet[userId]
	isEmpty := true
	for user, amount := range balances {
		if amount != 0 {
			isEmpty = false
			em.PrintBalance(userId, user, amount)
		}
	}

	if isEmpty {
		fmt.Println("No balances")
	}
}

func (em *ExpenseManager) ShowBalances() {
	fmt.Println("\n==========================================================Overall Summary of Bill Settlement:==========================================================")
	printed := make(map[string]bool)
	isEmpty := true
	for user, balances := range em.BalanceSheet {
		for otherUser, amount := range balances {
			if amount != 0 && !printed[user+otherUser] {
				isEmpty = false
				em.PrintBalance(user, otherUser, amount)
				printed[user+otherUser] = true
				printed[otherUser+user] = true
			}
		}
	}

	if isEmpty {
		fmt.Println("No balances")
	}
}

func (em *ExpenseManager) PrintBalance(user1, user2 string, amount float64) {
	user1Name := em.UserMap[user1].Name
	user2Name := em.UserMap[user2].Name
	if amount < 0 {
		fmt.Printf("%s owes %s: %.2f\n", user1Name, user2Name, -amount)
	} else if amount > 0 {
		fmt.Printf("%s owes %s: %.2f\n", user2Name, user1Name, amount)
	}
}

func main() {
	func() {
		user1 := NewUser("1", "Alice")
		user2 := NewUser("2", "Bob")
		user3 := NewUser("3", "Charlie")

		expenseManager := NewExpenseManager()

		expenseManager.AddUser(user1)
		expenseManager.AddUser(user2)
		expenseManager.AddUser(user3)

		// Dinner expense
		metadata := &ExpenseMetadata{Name: "Dinner"}
		equalSplits := []Split{
			&EqualSplit{User: user1},
			&EqualSplit{User: user2},
			&EqualSplit{User: user3},
		}
		expenseManager.AddExpense("EQUAL", 90, "1", equalSplits, metadata)

		// Cruise expense
		metadata.Name = "Cruise"
		expenseManager.AddExpense("EXACT", 150, "2", []Split{
			&ExactSplit{User: user1, Amount: 50},
			&ExactSplit{User: user2, Amount: 50},
			&ExactSplit{User: user3, Amount: 50},
		}, metadata)

		// Breakfast expense
		metadata.Name = "Breakfast"
		expenseManager.AddExpense("PERCENT", 50, "1", []Split{
			&PercentSplit{User: user1, Percent: 50},
			&PercentSplit{User: user2, Percent: 50},
		}, metadata)

		// Lunch expense
		metadata.Name = "Lunch"
		expenseManager.AddExpense("EQUAL", 120, "1", []Split{
			&EqualSplit{User: user1},
			&EqualSplit{User: user2},
		}, metadata)

		// Display Overall Summary
		expenseManager.ShowBalances()
	}()
}
