package main
import "fmt"

type Order struct {
    ID       int
    Customer string
    IsPremium bool // true for premium/quality order, false for normal
}

func main() {
    // Sample list of orders
    orders := []Order{
        {ID: 1, Customer: "Alice", IsPremium: false},
        {ID: 2, Customer: "Bob", IsPremium: true},
        {ID: 3, Customer: "Charlie", IsPremium: false},
        {ID: 4, Customer: "Diana", IsPremium: true},
    }

    // Process each order using a loop
    for _, order := range orders {
        // Implement Selection
        if order.IsPremium {
            fmt.Printf("Processing PREMIUM order #%d for %s.\n", order.ID, order.Customer)
        } else {
            fmt.Printf("Processing standard order #%d for %s.\n", order.ID, order.Customer)
        }
    }
}