package memory

import (
	"sync"
	"testing"
)

func Test_threadSafety_nonceProvider_Provide(t *testing.T) {
	t.Run("100 concurrent goroutines increment", func(t *testing.T) {
		concurrency := 1_000
		sum := 100_000_000
		work := sum / concurrency

		p := &nonceProvider{}

		var wg sync.WaitGroup
		wg.Add(concurrency)

		for i := 0; i < concurrency; i++ {
			go func() {
				for j := 0; j < work; j++ {
					_ = p.Provide()
				}

				wg.Done()
			}()
		}

		wg.Wait()

		if p.nonce != uint64(sum) {
			t.Errorf("current nonce is = %v, want %v", p.nonce, sum)
		}
	})
}
