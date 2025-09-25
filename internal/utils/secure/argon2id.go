package secure

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Params struct {
	Memory  uint32 
	Time    uint32 
	Threads uint8
	SaltLen uint32
	KeyLen  uint32
}

var Default = &Params{
	Memory:  64 * 1024, 
	Time:    3,
	Threads: 2,
	SaltLen: 16,
	KeyLen:  32,
}

func Hash(plain string, p *Params) (string, error) {
	if strings.TrimSpace(plain) == "" {
		return "", errors.New("empty password")
	}
	if p == nil { p = Default }

	salt := make([]byte, p.SaltLen)
	if _, err := rand.Read(salt); err != nil { return "", err }

	key := argon2.IDKey([]byte(plain), salt, p.Time, p.Memory, p.Threads, p.KeyLen)
	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		p.Memory, p.Time, p.Threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	), nil
}

func Verify(encoded, plain string) bool {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 { return false }

	var m uint32; var t uint32; var p uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &m, &t, &p); err != nil { return false }

	salt, err := base64.RawStdEncoding.DecodeString(parts[4]); if err != nil { return false }
	want, err := base64.RawStdEncoding.DecodeString(parts[5]); if err != nil { return false }

	got := argon2.IDKey([]byte(plain), salt, t, m, p, uint32(len(want)))
	return subtle.ConstantTimeCompare(want, got) == 1
}
