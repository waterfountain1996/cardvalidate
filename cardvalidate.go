// Package cardvalidate provides utilities for validating credit card information.
package cardvalidate

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/waterfountain1996/cardvalidate/issuer"
)

// Time layout for credit card's expiration date.
const expDateLayout = "01/2006"

var (
	ErrMalformedDate        = errors.New("cardvalidate: malformed expiration date")
	ErrCardExpired          = errors.New("cardvalidate: credit card has expired")
	ErrMalformedNumber      = errors.New("cardvalidate: malformed card number")
	ErrUnknownIssuer        = errors.New("cardvalidate: unknown card issuer")
	ErrInvalidAccountNumber = errors.New("cardvalidate: invalid account number")
)

// Validate validates credit card number and its expiration date.
// Validate only checks that cardNumber is structured accoring to ISO/IEC 7812, not
// whether it's an actual valid account number.
func Validate(cardNumber, expDate string) error {
	return validate(cardNumber, expDate, time.Now().UTC())
}

// Validate validates credit card number and its expiration date against currentDate.
func validate(cardNumber, expDate string, currentDate time.Time) error {
	if !validCardNumber(cardNumber) {
		return ErrMalformedNumber
	}

	if i := issuer.Identify(cardNumber); i == issuer.Unknown {
		return ErrUnknownIssuer
	}

	if !luhnCheck(cardNumber) {
		return ErrInvalidAccountNumber
	}

	exp, err := time.Parse(expDateLayout, expDate)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrMalformedDate, expDate)
	}

	if exp.Before(currentDate) || exp.Equal(currentDate) {
		return ErrCardExpired
	}

	return nil
}

// validCardNumber checks if cardNumber only consists of digits.
func validCardNumber(cardNumber string) bool {
	if len(cardNumber) < 8 || len(cardNumber) > 19 {
		return false
	}

	for _, r := range cardNumber {
		if !isDigit(r) {
			return false
		}
	}
	return true
}

// isDigit checks if r is an ASCII digit.
func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

// luhnCheck does Luhn's check (mod 10 check) validates credit card number.
func luhnCheck(cardNumber string) bool {
	digits := []rune(cardNumber)
	slices.Reverse(digits)

	sum := 0
	for i, d := range digits {
		n := int(d - '0')

		multiplier := 1
		if i%2 != 0 {
			multiplier = 2
		}

		x := n * multiplier
		if x > 9 {
			x -= 9
		}

		sum += x
	}
	return sum%10 == 0
}
