// Package issuer defines a list of known credit card issuers.
package issuer

import (
	"fmt"
)

// Credit card issuer.
type Issuer int

const (
	Unknown Issuer = iota
	AmericanExpress
	DinersClub
	Discover
	JCB
	MasterCard
	UnionPay
	Visa
)

// String implements fmt.Stringer
func (i Issuer) String() string {
	switch i {
	case AmericanExpress:
		return "American Express"
	case DinersClub:
		return "Diners Club"
	case Discover:
		return "Discover"
	case JCB:
		return "JCB"
	case MasterCard:
		return "MasterCard"
	case UnionPay:
		return "UnionPay"
	case Visa:
		return "Visa"
	default:
		return fmt.Sprintf("Unknown(%d)", i)
	}
}

// iinTrie is a prefix tree that contains all known credit card issuers' identification numbers.
var iinTrie = newTrie()

func init() {
	items := []struct {
		Issuer Issuer
		Prefix intRange // IIN range.
		Length intRange // Credit card number length.
	}{
		{AmericanExpress, newSingleIntRange(34), newSingleIntRange(15)},
		{AmericanExpress, newSingleIntRange(37), newSingleIntRange(15)},
		{DinersClub, newSingleIntRange(30), newSingleIntRange(14)},
		{DinersClub, newSingleIntRange(36), newSingleIntRange(14)},
		{DinersClub, newSingleIntRange(38), newSingleIntRange(14)},
		{DinersClub, newSingleIntRange(39), newSingleIntRange(14)},
		{Discover, newSingleIntRange(6011), newSingleIntRange(16)},
		{Discover, newIntRange(644, 649), newSingleIntRange(16)},
		{Discover, newSingleIntRange(65), newSingleIntRange(16)},
		{JCB, newIntRange(3528, 3589), newSingleIntRange(16)},
		{MasterCard, newIntRange(51, 55), newSingleIntRange(16)},
		{MasterCard, newIntRange(2221, 2720), newSingleIntRange(16)},
		{UnionPay, newSingleIntRange(62), newIntRange(13, 19)},
		{Visa, newSingleIntRange(4), newSingleIntRange(16)},
	}

	for _, item := range items {
		iinTrie.Put(item.Issuer, item.Prefix, item.Length)
	}
}

// Identify tries to identify the issuer of a given credit card number based on the
// list of known IINs and card number length.
// It does not do any validation and assumes that cardNumber contains only ASCII digits and
// may panic on non-digit characters.
func Identify(cardNumber string) Issuer {
	issuer, length := iinTrie.Get(cardNumber)
	if length.Contains(len(cardNumber)) {
		return issuer
	}
	return Unknown
}
