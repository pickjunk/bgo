package bgo

// CtxKey bgo's context key
// https://medium.com/@matryer/context-keys-in-go-5312346a868d
type CtxKey string

func (c CtxKey) String() string {
	return "bgo context key: " + string(c)
}
