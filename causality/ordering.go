package causality

// CausalOrder indicates causal ordering between two events, usually expressed
// in a form like  A -> B (event A is causually ordered before B) and tracked
// using some form of logical clock (i.e. Lamport or Vector clock).
//
//  A good primer on causal ordering is http://scattered-thoughts.net/blog/2012/08/16/causal-ordering/
//
type CausalOrder uint8

const (
	// OrderEqual indicates that events A and B happened simultaneously in the same
	// timeline.
	OrderEqual CausalOrder = iota

	// OrderGreater indicates that clock A > clock B, and therefore
	// event B -> event A
	OrderGreater

	// OrderLess indicates that clock B > clock A, and therefore
	// event A -> event B
	OrderLess

	// OrderNone indicates that two actors come from concurrent timelines. This
	// is possible because we are modelling a partial, rather than total, ordering.
	OrderNone
)

var (
	// CausalOrderDescriptions is a map of human readable versions of the
	// CausualOrder constants
	CausalOrderDescriptions = map[CausalOrder]string{}
)

func init() {
	CausalOrderDescriptions[OrderEqual] = "OrderEqual"
	CausalOrderDescriptions[OrderGreater] = "OrderGreater"
	CausalOrderDescriptions[OrderLess] = "OrderLess"
	CausalOrderDescriptions[OrderNone] = "OrderNone"
}

func (c CausalOrder) String() string {
	return CausalOrderDescriptions[c]
}
