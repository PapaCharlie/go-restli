package d2

import (
	"testing"
	"time"
)

func TestOnceWg(t *testing.T) {
	wg := new(OnceWg)
	go func() {
		time.Sleep(100 * time.Millisecond)
		wg.Done()
	}()
	wg.Wait()
	wg.Done()
	wg.Wait()
}
