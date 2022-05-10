package timeutils

import "time"

func StopTicker(t *time.Ticker) {
	t.Stop()

	select {
	case <-t.C:
	default:
	}
}
