# Splitwise_app
## Expense Manager

This Go project implements a Bill Splitting and Expense Management system that allows users to split expenses across different events such as Dinner, Cruise, Breakfast, and Lunch. The system can handle equal, exact, and percentage-based splits and calculates the balance for each user. It also provides summaries of bill settlements per event and overall summaries for all events.

### Features

- **Multiple Split Types:** Handle equal, exact, and percentage splits for expenses.
- **Expense Events:** Track expenses across different events (Dinner, Cruise, etc.).
- **Balance Management:** Automatically updates and calculates balances between users.
- **Event and Overall Summary:** Generate a summary of bill settlements for each event and an overall summary.

### Components

- **User:** Represents a participant in the expense split.
- **Split:** Interface for different split types (EqualSplit, ExactSplit, PercentSplit).
- **Expense:** Interface for different types of expenses (EqualExpense, ExactExpense, PercentExpense).
- **ExpenseManager:** Manages users, expenses, and balances.
- **ExpenseService:** Creates and validates expenses based on the type (Equal, Exact, Percent).

### Structure

```bash
|-- main.go            # Main application logic
|-- README.md          # Project documentation
|-- expense_test.go    # Unit tests for the Expense Manager
```

### Requirements

- **Go version:** 1.18 or later
- **Operating System:** Cross-platform (Windows, macOS, Linux)

This will execute the example scenarios and print the summaries of the bill settlements.

### Example Usage

The following code demonstrates how you can add users and expenses to the system and get the event-wise and overall bill settlements:

```go
package main

func main() {
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
}
```

### Example Output

```yaml
Summary of Bill Settlement for Dinner:
Alice owes Bob: 30.00
Alice owes Charlie: 30.00

Summary of Bill Settlement for Cruise:
Alice owes Charlie: 30.00
Bob owes Alice: 20.00
Bob owes Charlie: 50.00

Summary of Bill Settlement for Breakfast:
Alice owes Charlie: 30.00
Bob owes Charlie: 50.00
Alice owes Bob: 5.00

Summary of Bill Settlement for Lunch:
Alice owes Bob: 65.00
Alice owes Charlie: 30.00
Bob owes Charlie: 50.00

==========================================================Overall Summary of Bill Settlement:==========================================================  
Bob owes Charlie: 50.00
Alice owes Bob: 65.00
Alice owes Charlie: 30.00
```

### Unit Tests

Unit tests for the Expense Manager are provided in the `expense_test.go` file. You can run the tests with the following command:

```bash
go test -v
```
git config --global core.autocrlf true
