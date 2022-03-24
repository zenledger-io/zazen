package timeutils

import "time"

func StopTimer(t *time.Timer) {
	if t.Stop() {
		return
	}

	select {
	case <-t.C:
	default:
	}
}
