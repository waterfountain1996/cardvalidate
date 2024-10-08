package issuer

import "strconv"

// intRange is helper for storing numeric ranges.
type intRange struct {
	Start int
	End   int
}

// newIntRange returns a new inRange.
func newIntRange(start, end int) intRange {
	return intRange{
		Start: start,
		End:   end,
	}
}

// newSingleIntRange returns a new intRange where Start is equal to End.
func newSingleIntRange(start int) intRange {
	return intRange{
		Start: start,
		End:   start,
	}
}

// Contains checks if n fits inside the range.
func (r intRange) Contains(n int) bool {
	return r.Start <= n && n <= r.End
}

// trie is a prefix-tree structure that acts as a mapping of card number prefixes to issuers.
// Since credit card numbers are composed of digits only, we use an array instead of a map for
// trie's children where leaf index is ASCII digit - '0', therefore Get will panic on out-of-bounds
// access if called with a key that contains non-digit characters.
type trie struct {
	issuer   Issuer
	length   intRange
	children [10]*trie
}

// newTrie allocates a new trie.
func newTrie() *trie {
	return &trie{}
}

// Get takes a credit card number and returns Issuer that matched its prefix plus the valid
// card number length range for this IIN.
func (t *trie) Get(key string) (Issuer, intRange) {
	node := t
	for _, r := range key {
		node = node.children[r-'0']
		if node == nil {
			return Unknown, intRange{}
		}

		if node.issuer != Unknown {
			break
		}
	}
	return node.issuer, node.length
}

// Put adds a new IIN into the trie.
func (t *trie) Put(issuer Issuer, prefix, length intRange) {
	for key := prefix.Start; key <= prefix.End; key++ {
		t.put(strconv.Itoa(key), issuer, length)
	}
}

// put adds a new leaf for given key into the trie and panics if it's a duplicate.
func (t *trie) put(key string, issuer Issuer, length intRange) {
	node := t
	for _, r := range key {
		child := node.children[r-'0']
		if child == nil {
			child = newTrie()
			node.children[r-'0'] = child
		}
		node = child
	}
	if node.issuer != Unknown {
		panic("duplicate trie entry")
	}
	node.issuer = issuer
	node.length = length
}
