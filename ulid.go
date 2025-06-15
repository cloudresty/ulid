package ulid

import (
	"crypto/rand"
	"errors"
	"sync"
	"time"
)

const (
	encodedLength   = 26
	timestampBits   = 48
	randomnessBits  = 80
	randomnessBytes = 10 // 80 bits = 10 bytes
	timestampBytes  = 6  // 48 bits = 6 bytes
	totalBytes      = timestampBytes + randomnessBytes
	maxTimestamp    = (1 << timestampBits) - 1
)

// Crockford Base32 alphabet in lowercase for better readability
const crockfordAlphabet = "0123456789abcdefghjkmnpqrstvwxyz"

var (
	// Pre-computed encoding/decoding tables for performance
	encodeTable [32]byte
	decodeTable [256]byte

	// Monotonicity state
	lastTime       uint64
	lastRandomness [randomnessBytes]byte
	mutex          sync.Mutex
)

func init() {
	// Initialize encoding table
	copy(encodeTable[:], crockfordAlphabet)

	// Initialize decoding table with invalid values
	for i := range decodeTable {
		decodeTable[i] = 0xFF
	}

	// Map valid characters to their values (case insensitive)
	for i, c := range crockfordAlphabet {
		decodeTable[c] = byte(i)
		if c >= 'a' && c <= 'z' {
			decodeTable[c-'a'+'A'] = byte(i) // uppercase variants
		}
	}

	// Handle ambiguous characters as per Crockford spec
	decodeTable['I'] = decodeTable['1']
	decodeTable['i'] = decodeTable['1']
	decodeTable['L'] = decodeTable['1']
	decodeTable['l'] = decodeTable['1']
	decodeTable['O'] = decodeTable['0']
	decodeTable['o'] = decodeTable['0']
	decodeTable['U'] = decodeTable['V']
	decodeTable['u'] = decodeTable['v']
}

// ULID represents a Universally Unique Lexicographically Sortable Identifier.
type ULID struct {
	timestamp  uint64
	randomness [randomnessBytes]byte
}

// fastEncode encodes 16 bytes to 26-character string using Crockford Base32
func fastEncode(data [totalBytes]byte) string {
	result := make([]byte, encodedLength)

	// Convert 16 bytes (128 bits) to base32 (5 bits per char = 26 chars)
	var acc uint64
	var bits uint
	j := 0

	for i := 0; i < totalBytes; i++ {
		acc = (acc << 8) | uint64(data[i])
		bits += 8

		for bits >= 5 {
			bits -= 5
			result[j] = encodeTable[(acc>>bits)&0x1F]
			j++
		}
	}

	// Handle remaining bits
	if bits > 0 {
		result[j] = encodeTable[(acc<<(5-bits))&0x1F]
	}

	return string(result)
}

// fastDecode decodes 26-character string to 16 bytes
func fastDecode(s string) ([totalBytes]byte, error) {
	var result [totalBytes]byte

	if len(s) != encodedLength {
		return result, errors.New("invalid ULID length")
	}

	var acc uint64
	var bits uint
	j := 0

	for i := 0; i < encodedLength; i++ {
		if int(s[i]) >= len(decodeTable) {
			return result, errors.New("invalid character in ULID")
		}

		val := decodeTable[s[i]]
		if val == 0xFF {
			return result, errors.New("invalid character in ULID")
		}

		acc = (acc << 5) | uint64(val)
		bits += 5

		if bits >= 8 && j < totalBytes {
			bits -= 8
			result[j] = byte(acc >> bits)
			j++
		}
	}

	return result, nil
}

// String returns the canonical string representation of the ULID.
func (u ULID) String() string {
	var data [totalBytes]byte

	// Encode timestamp (big-endian)
	for i := 0; i < timestampBytes; i++ {
		data[i] = byte(u.timestamp >> (8 * (timestampBytes - 1 - i)))
	}

	// Copy randomness
	copy(data[timestampBytes:], u.randomness[:])

	return fastEncode(data)
}

// Parse parses a ULID string and returns a ULID struct.
func Parse(s string) (ULID, error) {
	data, err := fastDecode(s)
	if err != nil {
		return ULID{}, err
	}

	// Extract timestamp (big-endian)
	timestamp := uint64(0)
	for i := 0; i < timestampBytes; i++ {
		timestamp = (timestamp << 8) | uint64(data[i])
	}

	// Extract randomness
	var randomness [randomnessBytes]byte
	copy(randomness[:], data[timestampBytes:])

	return ULID{
		timestamp:  timestamp,
		randomness: randomness,
	}, nil
}

// GetTime returns the timestamp of the ULID in milliseconds.
func (u ULID) GetTime() uint64 {
	return u.timestamp
}

// generateRandomness generates cryptographically secure random bytes
func generateRandomness() ([randomnessBytes]byte, error) {
	var randomness [randomnessBytes]byte
	_, err := rand.Read(randomness[:])
	return randomness, err
}

// incrementRandomness increments the randomness component by 1
// Returns true if overflow occurred
func incrementRandomness(r *[randomnessBytes]byte) bool {
	for i := randomnessBytes - 1; i >= 0; i-- {
		r[i]++
		if r[i] != 0 {
			return false // No overflow
		}
	}
	return true // Overflow occurred
}

// compareRandomness compares two randomness arrays
// Returns: -1 if a < b, 0 if a == b, 1 if a > b
func compareRandomness(a, b [randomnessBytes]byte) int {
	for i := 0; i < randomnessBytes; i++ {
		if a[i] < b[i] {
			return -1
		}
		if a[i] > b[i] {
			return 1
		}
	}
	return 0
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

	// Handle monotonicity
	if timestamp == lastTime {
		// Check if we need to increment randomness for monotonicity
		if compareRandomness(randomness, lastRandomness) <= 0 {
			randomness = lastRandomness
			if incrementRandomness(&randomness) {
				// Randomness overflow, increment timestamp
				timestamp++
				if timestamp > maxTimestamp {
					return "", errors.New("timestamp out of range due to randomness exhaustion")
				}
				// Generate new randomness for new timestamp
				randomness, err = generateRandomness()
				if err != nil {
					return "", err
				}
			}
		}
	}

	lastTime = timestamp
	lastRandomness = randomness

	ulid := ULID{
		timestamp:  timestamp,
		randomness: randomness,
	}

	return ulid.String(), nil
}
