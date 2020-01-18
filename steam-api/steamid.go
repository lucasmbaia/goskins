package steam

import (
	"strconv"
)

type SteamID uint64

func (s *SteamID) ToString() string {
	return strconv.FormatUint(uint64(*s), 10)
}
