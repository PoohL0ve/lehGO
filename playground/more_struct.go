package main

import "fmt"

// Compounding struct (Address)
type Address struct {
	City    string
	Country string
}

// Parent/Outer struct
type User struct {
	Username string
	Email    string
	Address  // Embedded struct (Anonymous field: no explicit field name)
}

// Value Receiver: Receives a copy. Original struct remains untouched.
func (u User) DisplayProfile() {
	// City and Country are "promoted" fields. We don't need to write u.Address.City!
	fmt.Printf("%s lives in %s, %s\n", u.Username, u.City, u.Country)
}

// Pointer Receiver: Can modify the actual struct fields in-place.
func (u *User) UpdateEmail(newEmail string) {
	u.Email = newEmail
}

func main() {
	// Initializing with field names
	user := User{
		Username: "gopher123",
		Email:    "go@lang.org",
		Address: Address{
			City:    "San Francisco",
			Country: "USA",
		},
	}

	// 1. Call value receiver method
	user.DisplayProfile()

	// 2. Call pointer receiver method (Go automatically takes the address under the hood: (&user).UpdateEmail)
	user.UpdateEmail("new_email@lang.org")
	fmt.Println("Updated Email:", user.Email)
}
