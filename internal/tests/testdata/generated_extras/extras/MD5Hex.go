package extras

import "encoding/hex"

func (m *MD5) Hex() string {
	return hex.EncodeToString(m[:])
}
