package ulid

import (
	"crypto/rand"
	"errors"
	"sync"
	"time"
	"unsafe"
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
	// Pre-computed encoding/decoding tables aligned for cache efficiency
	encodeTable = [32]byte{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'j', 'k',
		'm', 'n', 'p', 'q', 'r', 's', 't', 'v', 'w', 'x',
		'y', 'z',
	}
	decodeTable [256]byte

	// Monotonicity state with CPU cache alignment
	lastTime       uint64
	lastRandomness [randomnessBytes]byte
	mutex          sync.Mutex
)

func init() {
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

// ultraFastEncode uses highly optimized base32 encoding with SIMD-style operations
func ultraFastEncode(data [totalBytes]byte) string {
	// Stack allocation for result - no heap allocation
	var result [encodedLength]byte

	// Ultra-optimized encoding using 64-bit operations and parallel processing
	// This approach minimizes CPU cycles by processing multiple bytes simultaneously

	// Process first 8 bytes as a single 64-bit word
	word1 := uint64(data[0])<<56 | uint64(data[1])<<48 | uint64(data[2])<<40 | uint64(data[3])<<32 |
		uint64(data[4])<<24 | uint64(data[5])<<16 | uint64(data[6])<<8 | uint64(data[7])

	// Extract 13 characters from 64 bits (65 bits total, 3 bits overflow)
	result[0] = encodeTable[word1>>59]
	result[1] = encodeTable[(word1>>54)&0x1F]
	result[2] = encodeTable[(word1>>49)&0x1F]
	result[3] = encodeTable[(word1>>44)&0x1F]
	result[4] = encodeTable[(word1>>39)&0x1F]
	result[5] = encodeTable[(word1>>34)&0x1F]
	result[6] = encodeTable[(word1>>29)&0x1F]
	result[7] = encodeTable[(word1>>24)&0x1F]
	result[8] = encodeTable[(word1>>19)&0x1F]
	result[9] = encodeTable[(word1>>14)&0x1F]
	result[10] = encodeTable[(word1>>9)&0x1F]
	result[11] = encodeTable[(word1>>4)&0x1F]
	result[12] = encodeTable[((word1&0x0F)<<1)|(uint64(data[8])>>7)]

	// Process remaining 8 bytes
	word2 := uint64(data[8]&0x7F)<<57 | uint64(data[9])<<49 | uint64(data[10])<<41 | uint64(data[11])<<33 |
		uint64(data[12])<<25 | uint64(data[13])<<17 | uint64(data[14])<<9 | uint64(data[15])<<1

	result[13] = encodeTable[word2>>59]
	result[14] = encodeTable[(word2>>54)&0x1F]
	result[15] = encodeTable[(word2>>49)&0x1F]
	result[16] = encodeTable[(word2>>44)&0x1F]
	result[17] = encodeTable[(word2>>39)&0x1F]
	result[18] = encodeTable[(word2>>34)&0x1F]
	result[19] = encodeTable[(word2>>29)&0x1F]
	result[20] = encodeTable[(word2>>24)&0x1F]
	result[21] = encodeTable[(word2>>19)&0x1F]
	result[22] = encodeTable[(word2>>14)&0x1F]
	result[23] = encodeTable[(word2>>9)&0x1F]
	result[24] = encodeTable[(word2>>4)&0x1F]
	result[25] = encodeTable[(word2<<1)&0x1F]

	// Zero-copy string conversion using unsafe
	return unsafe.String(&result[0], encodedLength)
}

// ultraFastDecode decodes with minimal validation and optimized bit operations
func ultraFastDecode(s string) ([totalBytes]byte, error) {
	var result [totalBytes]byte

	if len(s) != encodedLength {
		return result, errors.New("invalid ULID length")
	}

	// Branch-free validation using lookup table
	// First pass: validate all characters
	for i := range encodedLength {
		c := s[i]
		if int(c) >= 256 || decodeTable[c] == 0xFF {
			return result, errors.New("invalid character in ULID")
		}
	}

	// Optimized decoding in 8-character chunks
	// First 8 chars -> 5 bytes
	v0, v1, v2, v3, v4, v5, v6, v7 := decodeTable[s[0]], decodeTable[s[1]], decodeTable[s[2]], decodeTable[s[3]],
		decodeTable[s[4]], decodeTable[s[5]], decodeTable[s[6]], decodeTable[s[7]]

	acc := uint64(v0)<<35 | uint64(v1)<<30 | uint64(v2)<<25 | uint64(v3)<<20 |
		uint64(v4)<<15 | uint64(v5)<<10 | uint64(v6)<<5 | uint64(v7)

	result[0] = byte(acc >> 32)
	result[1] = byte(acc >> 24)
	result[2] = byte(acc >> 16)
	result[3] = byte(acc >> 8)
	result[4] = byte(acc)

	// Next 8 chars -> 5 bytes
	v0, v1, v2, v3, v4, v5, v6, v7 = decodeTable[s[8]], decodeTable[s[9]], decodeTable[s[10]], decodeTable[s[11]],
		decodeTable[s[12]], decodeTable[s[13]], decodeTable[s[14]], decodeTable[s[15]]

	acc = uint64(v0)<<35 | uint64(v1)<<30 | uint64(v2)<<25 | uint64(v3)<<20 |
		uint64(v4)<<15 | uint64(v5)<<10 | uint64(v6)<<5 | uint64(v7)

	result[5] = byte(acc >> 32)
	result[6] = byte(acc >> 24)
	result[7] = byte(acc >> 16)
	result[8] = byte(acc >> 8)
	result[9] = byte(acc)

	// Next 8 chars -> 5 bytes
	v0, v1, v2, v3, v4, v5, v6, v7 = decodeTable[s[16]], decodeTable[s[17]], decodeTable[s[18]], decodeTable[s[19]],
		decodeTable[s[20]], decodeTable[s[21]], decodeTable[s[22]], decodeTable[s[23]]

	acc = uint64(v0)<<35 | uint64(v1)<<30 | uint64(v2)<<25 | uint64(v3)<<20 |
		uint64(v4)<<15 | uint64(v5)<<10 | uint64(v6)<<5 | uint64(v7)

	result[10] = byte(acc >> 32)
	result[11] = byte(acc >> 24)
	result[12] = byte(acc >> 16)
	result[13] = byte(acc >> 8)
	result[14] = byte(acc)

	// Last 2 chars -> 1 byte
	result[15] = decodeTable[s[24]]<<3 | decodeTable[s[25]]>>2

	return result, nil
}

// String returns the canonical string representation of the ULID.
func (u ULID) String() string {
	var data [totalBytes]byte

	// Encode timestamp (big-endian) - unrolled for speed
	data[0] = byte(u.timestamp >> 40)
	data[1] = byte(u.timestamp >> 32)
	data[2] = byte(u.timestamp >> 24)
	data[3] = byte(u.timestamp >> 16)
	data[4] = byte(u.timestamp >> 8)
	data[5] = byte(u.timestamp)

	// Copy randomness - compiler will optimize this
	copy(data[timestampBytes:], u.randomness[:])

	return ultraFastEncode(data)
}

// Parse parses a ULID string and returns a ULID struct.
func Parse(s string) (ULID, error) {
	data, err := ultraFastDecode(s)
	if err != nil {
		return ULID{}, err
	}

	// Extract timestamp (big-endian) - unrolled for speed
	timestamp := uint64(data[0])<<40 | uint64(data[1])<<32 | uint64(data[2])<<24 |
		uint64(data[3])<<16 | uint64(data[4])<<8 | uint64(data[5])

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
	for i := range randomnessBytes {
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
// Hyper-optimized version that avoids all unnecessary allocations
func NewTime(timestamp uint64) (string, error) {
	if timestamp > maxTimestamp {
		return "", errors.New("timestamp out of range")
	}

	randomness, err := generateRandomness()
	if err != nil {
		return "", err
	}

	// Critical section optimized for minimal lock time
	mutex.Lock()
	if timestamp == lastTime {
		// Inline comparison for maximum speed
		needIncrement := true
		for i := 0; i < randomnessBytes && needIncrement; i++ {
			if randomness[i] > lastRandomness[i] {
				needIncrement = false
			} else if randomness[i] < lastRandomness[i] {
				needIncrement = true
				break
			}
		}

		if needIncrement {
			// Fast copy and increment
			copy(randomness[:], lastRandomness[:])

			// Unrolled increment for maximum speed
			randomness[9]++
			if randomness[9] == 0 {
				randomness[8]++
				if randomness[8] == 0 {
					randomness[7]++
					if randomness[7] == 0 {
						randomness[6]++
						if randomness[6] == 0 {
							randomness[5]++
							if randomness[5] == 0 {
								randomness[4]++
								if randomness[4] == 0 {
									randomness[3]++
									if randomness[3] == 0 {
										randomness[2]++
										if randomness[2] == 0 {
											randomness[1]++
											if randomness[1] == 0 {
												randomness[0]++
												if randomness[0] == 0 {
													// Overflow - increment timestamp
													timestamp++
													if timestamp > maxTimestamp {
														mutex.Unlock()
														return "", errors.New("timestamp out of range due to randomness exhaustion")
													}
													randomness, err = generateRandomness()
													if err != nil {
														mutex.Unlock()
														return "", err
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	lastTime = timestamp
	lastRandomness = randomness
	mutex.Unlock()

	// Direct encoding without intermediate ULID struct allocation
	var data [totalBytes]byte

	// Unrolled timestamp encoding
	data[0] = byte(timestamp >> 40)
	data[1] = byte(timestamp >> 32)
	data[2] = byte(timestamp >> 24)
	data[3] = byte(timestamp >> 16)
	data[4] = byte(timestamp >> 8)
	data[5] = byte(timestamp)

	// Copy randomness
	copy(data[6:], randomness[:])

	return ultraFastEncode(data), nil
}
