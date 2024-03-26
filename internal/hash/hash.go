package hash

// Hasher is an interface for types which are able to produce a hash based on
// their internal state.
type Hasher interface {
	Hash() (string, error)
}
