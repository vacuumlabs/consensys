package dao

type NonceProvider interface {
	Provide() uint64
}
