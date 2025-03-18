package ulid

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"math/big"
	"sync"
	"time"
)

const (
	encodedLength  = 26
	timestampBits  = 48
	randomnessBits = 80
	maxTimestamp   = (1 << timestampBits) - 1
)

var (
	encoding       = base32.NewEncoding("0123456789ABCDEFGHJKMNPQRSTVWXYZ").WithPadding(base32.NoPadding)
	lastTime       uint64
	lastRandomness *big.Int
	mutex          sync.Mutex
	maxRandomness  *big.Int
)

func init() {
	maxRandomness = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), randomnessBits), big.NewInt(1))
}

// ULID represents a Universally Unique Lexicographically Sortable Identifier.
type ULID struct {
	timestamp  uint64
	randomness *big.Int
}

// String returns the canonical string representation of the ULID.
func (u ULID) String() string {
	timestampBytes := make([]byte, 6)
	for i := 5; i >= 0; i-- {
		timestampBytes[i] = byte(u.timestamp & 0xFF)
		u.timestamp >>= 8
	}

	randomnessBytes := u.randomness.Bytes()

	// Pad randomnessBytes to 10 bytes if needed.
	if len(randomnessBytes) < 10 {
		paddedRandomnessBytes := make([]byte, 10)
		copy(paddedRandomnessBytes[10-len(randomnessBytes):], randomnessBytes)
		randomnessBytes = paddedRandomnessBytes
	}

	combinedBytes := append(timestampBytes, randomnessBytes...)
	encoded := encoding.EncodeToString(combinedBytes)

	return encoded
}

// Parse parses a ULID string and returns a ULID struct.
func Parse(s string) (ULID, error) {
	if len(s) != encodedLength {
		return ULID{}, errors.New("invalid ULID length")
	}

	decoded, err := encoding.DecodeString(s)
	if err != nil {
		return ULID{}, err
	}

	timestampBytes := decoded[:6]
	randomnessBytes := decoded[6:]

	timestamp := uint64(0)
	for _, b := range timestampBytes {
		timestamp = (timestamp << 8) | uint64(b)
	}

	randomness := new(big.Int).SetBytes(randomnessBytes)

	return ULID{
		timestamp:  timestamp,
		randomness: randomness,
	}, nil
}

// GetTime returns the timestamp of the ULID in milliseconds.
func (u ULID) GetTime() uint64 {
	return u.timestamp
}

func generateRandomness() (*big.Int, error) {
	randomBytes := make([]byte, 10)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	randomness := new(big.Int).SetBytes(randomBytes)

	return randomness.Mod(randomness, new(big.Int).Add(maxRandomness, big.NewInt(1))), nil

}

// New returns a new ULID.
func New() (string, error) {
	return NewTime(uint64(time.Now().UnixMilli()))
}

// NewTime returns a new ULID with the given timestamp in milliseconds.
func NewTime(timestamp uint64) (string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if timestamp > maxTimestamp {
		return "", errors.New("timestamp out of range")
	}

	randomness, err := generateRandomness()
	if err != nil {
		return "", err
	}

	if timestamp == lastTime {
		if lastRandomness != nil && randomness.Cmp(lastRandomness) <= 0 {
			randomness = new(big.Int).Add(lastRandomness, big.NewInt(1))

			if randomness.Cmp(maxRandomness) > 0 {
				//Randomness wrapped around. Increment time.
				timestamp++
				randomness, err = generateRandomness()
				if err != nil {
					return "", err
				}
				if timestamp > maxTimestamp {
					return "", errors.New("timestamp out of range, due to randomness exhaustion")
				}
			}
		}
	}

	lastTime = timestamp
	lastRandomness = new(big.Int).Set(randomness)

	ulid := ULID{
		timestamp:  timestamp,
		randomness: randomness,
	}

	return ulid.String(), nil
}
