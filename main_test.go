package main

import "testing"

// Test creating users
func TestCreateUser(t *testing.T) {
	user := NewUser("1", "Alice")
	if user.ID != "1" {
		t.Errorf("Expected user ID '1', got %s", user.ID)
	}
	if user.Name != "Alice" {
		t.Errorf("Expected user name 'Alice', got %s", user.Name)
	}
}

// Test adding users to ExpenseManager
func TestAddUserToManager(t *testing.T) {
	manager := NewExpenseManager()
	user := NewUser("1", "Alice")

	manager.AddUser(user)

	if len(manager.UserMap) != 1 {
		t.Errorf("Expected 1 user in the manager, got %d", len(manager.UserMap))
	}
	if _, exists := manager.UserMap["1"]; !exists {
		t.Errorf("User with ID '1' not found in manager")
	}
}

// Test equal expense creation and balance sheet update
func TestEqualExpense(t *testing.T) {
	manager := NewExpenseManager()

	user1 := NewUser("1", "Alice")
	user2 := NewUser("2", "Bob")
	user3 := NewUser("3", "Charlie")

	manager.AddUser(user1)
	manager.AddUser(user2)
	manager.AddUser(user3)

	metadata := &ExpenseMetadata{Name: "Dinner"}
	equalSplits := []Split{
		&EqualSplit{User: user1},
		&EqualSplit{User: user2},
		&EqualSplit{User: user3},
	}

	// Add equal expense
	manager.AddExpense("EQUAL", 90, "1", equalSplits, metadata)

	// Check the balance sheet
	if manager.BalanceSheet[user1.ID][user2.ID] != -30 {
		t.Errorf("Expected Alice owes Bob: -30, got %.2f", manager.BalanceSheet[user1.ID][user2.ID])
	}
	if manager.BalanceSheet[user1.ID][user3.ID] != -30 {
		t.Errorf("Expected Alice owes Charlie: -30, got %.2f", manager.BalanceSheet[user1.ID][user3.ID])
	}
	if manager.BalanceSheet[user2.ID][user1.ID] != 30 {
		t.Errorf("Expected Bob owes Alice: 30, got %.2f", manager.BalanceSheet[user2.ID][user1.ID])
	}
	if manager.BalanceSheet[user2.ID][user3.ID] != 30 {
		t.Errorf("Expected Bob owes Charlie: 30, got %.2f", manager.BalanceSheet[user2.ID][user3.ID])
	}
	if manager.BalanceSheet[user3.ID][user1.ID] != 30 {
		t.Errorf("Expected Charlie owes Alice: 30, got %.2f", manager.BalanceSheet[user3.ID][user1.ID])
	}
	if manager.BalanceSheet[user3.ID][user2.ID] != 30 {
		t.Errorf("Expected Charlie owes Bob: 30, got %.2f", manager.BalanceSheet[user3.ID][user2.ID])
	}
}

// Test exact expense creation and balance sheet update
func TestExactExpense(t *testing.T) {
	manager := NewExpenseManager()

	user1 := NewUser("1", "Alice")
	user2 := NewUser("2", "Bob")
	user3 := NewUser("3", "Charlie")

	manager.AddUser(user1)
	manager.AddUser(user2)
	manager.AddUser(user3)

	metadata := &ExpenseMetadata{Name: "Cruise"}
	exactSplits := []Split{
		&ExactSplit{User: user1, Amount: 50},
		&ExactSplit{User: user2, Amount: 50},
		&ExactSplit{User: user3, Amount: 50},
	}

	// Add exact expense
	manager.AddExpense("EXACT", 150, "2", exactSplits, metadata)

	// Check the balance sheet
	if manager.BalanceSheet[user1.ID][user2.ID] != 50 {
		t.Errorf("Expected Alice owes Bob: 50, got %.2f", manager.BalanceSheet[user1.ID][user2.ID])
	}
	if manager.BalanceSheet[user1.ID][user3.ID] != 50 {
		t.Errorf("Expected Alice owes Charlie: 50, got %.2f", manager.BalanceSheet[user1.ID][user3.ID])
	}
	if manager.BalanceSheet[user2.ID][user1.ID] != -50 {
		t.Errorf("Expected Bob owes Alice: -50, got %.2f", manager.BalanceSheet[user2.ID][user1.ID])
	}
	if manager.BalanceSheet[user2.ID][user3.ID] != 50 {
		t.Errorf("Expected Bob owes Charlie: 50, got %.2f", manager.BalanceSheet[user2.ID][user3.ID])
	}
	if manager.BalanceSheet[user3.ID][user1.ID] != -50 {
		t.Errorf("Expected Charlie owes Alice: -50, got %.2f", manager.BalanceSheet[user3.ID][user1.ID])
	}
	if manager.BalanceSheet[user3.ID][user2.ID] != -50 {
		t.Errorf("Expected Charlie owes Bob: -50, got %.2f", manager.BalanceSheet[user3.ID][user2.ID])
	}
}

// Test percent expense creation and balance sheet update
func TestPercentExpense(t *testing.T) {
	manager := NewExpenseManager()

	user1 := NewUser("1", "Alice")
	user2 := NewUser("2", "Bob")

	manager.AddUser(user1)
	manager.AddUser(user2)

	metadata := &ExpenseMetadata{Name: "Breakfast"}
	percentSplits := []Split{
		&PercentSplit{User: user1, Percent: 50},
		&PercentSplit{User: user2, Percent: 50},
	}

	// Add percent expense
	manager.AddExpense("PERCENT", 50, "1", percentSplits, metadata)

	// Check the balance sheet
	if manager.BalanceSheet[user1.ID][user2.ID] != -25 {
		t.Errorf("Expected Alice owes Bob: -25, got %.2f", manager.BalanceSheet[user1.ID][user2.ID])
	}
	if manager.BalanceSheet[user2.ID][user1.ID] != 25 {
		t.Errorf("Expected Bob owes Alice: 25, got %.2f", manager.BalanceSheet[user2.ID][user1.ID])
	}
}

// Test ExpenseService for creating equal expense
func TestExpenseServiceCreateEqualExpense(t *testing.T) {
	expenseService := ExpenseService{}

	user1 := NewUser("1", "Alice")
	user2 := NewUser("2", "Bob")

	equalSplits := []Split{
		&EqualSplit{User: user1},
		&EqualSplit{User: user2},
	}

	metadata := &ExpenseMetadata{Name: "Lunch"}
	expense := expenseService.CreateExpense("EQUAL", 60, user1, equalSplits, metadata)

	// Validate expense
	if expense == nil {
		t.Errorf("Expected a valid expense, got nil")
	}
	if expense.GetAmount() != 60 {
		t.Errorf("Expected expense amount 60, got %.2f", expense.GetAmount())
	}
	if len(expense.GetSplits()) != 2 {
		t.Errorf("Expected 2 splits, got %d", len(expense.GetSplits()))
	}
}

// Test ExpenseService for creating percent expense
func TestExpenseServiceCreatePercentExpense(t *testing.T) {
	expenseService := ExpenseService{}

	user1 := NewUser("1", "Alice")
	user2 := NewUser("2", "Bob")

	percentSplits := []Split{
		&PercentSplit{User: user1, Percent: 50},
		&PercentSplit{User: user2, Percent: 50},
	}

	metadata := &ExpenseMetadata{Name: "Snack"}
	expense := expenseService.CreateExpense("PERCENT", 40, user1, percentSplits, metadata)

	// Validate expense
	if expense == nil {
		t.Errorf("Expected a valid expense, got nil")
	}
	if expense.GetAmount() != 40 {
		t.Errorf("Expected expense amount 40, got %.2f", expense.GetAmount())
	}
	if len(expense.GetSplits()) != 2 {
		t.Errorf("Expected 2 splits, got %d", len(expense.GetSplits()))
	}
}
