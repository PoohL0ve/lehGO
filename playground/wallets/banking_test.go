package wallets

import "testing"

func TestWallet(t *testing.T) {

	// Avoid duplication for deposit and withdraw
	// Helper functions can also be in their own function
	assertBalance := func(t testing.TB, wallet Wallet, want Bitcoin) {
		t.Helper()

		got := wallet.Balance()

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}

	}

	assertError := func(t testing.TB, got, want error) {
		t.Helper()
		if got == nil {
			t.Fatal("wanted an error but didn't get one")
		}

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	}
	t.Run("Deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("Withdraw with funds", func(t *testing.T) {
		wallet := Wallet{balance: Bitcoin(20)}
		err := wallet.Withdraw(Bitcoin(15))

		assertNoError(t, err)
		assertBalance(t, wallet, Bitcoin(5))
	})

	t.Run("withdraw insufficient funds", func(t *testing.T) {
		startingBalance := Bitcoin(20)
		wallet := Wallet{startingBalance}
		err := wallet.Withdraw(Bitcoin(100))

		assertError(t, err, ErrInsufficientFunds)
		assertBalance(t, wallet, startingBalance)
	})

	t.Run("Bitcoin Test", func(t *testing.T) {
		btc := Bitcoin(10)
		got := btc.String()
		want := "10 BTC"

		if got != want {
			t.Errorf("got: %s, want: %s", got, want)
		}
	})
}

// Unchecked errors
func assertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatal("got an error but didn't want one")
	}
}
