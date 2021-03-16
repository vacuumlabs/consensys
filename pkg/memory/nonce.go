package memory

import (
	"sync/atomic"
	"vax/pkg/model/dao"
)

type nonceProvider struct {
	nonce uint64
}

func NewMonotonicNonceProvider() dao.NonceProvider {
	return &nonceProvider{}
}

func (p *nonceProvider) Provide() uint64 {
	return atomic.AddUint64(&p.nonce, 1)
}
