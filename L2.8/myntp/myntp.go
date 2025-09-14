package myntp

import (
	"time"

	"github.com/beevik/ntp"
)

// GetNTPTime возвращает текущее время через NTP либо ошибку
func GetNTPTime() (time.Time, error) {
	t, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
