package ulid

import (
	"math/big"
	"testing"
	"time"
)

func TestULIDEncodingDecoding(t *testing.T) {
	timestamp := uint64(time.Now().UnixMilli())
	randomness := new(big.Int).SetInt64(123456789)

	ulid := ULID{
		timestamp:  timestamp,
		randomness: randomness,
	}

	encoded := ulid.String()
	decoded, err := Parse(encoded)

	if err != nil {
		t.Fatalf("Error parsing ULID: %v", err)
	}

	if decoded.timestamp != timestamp {
		t.Errorf("Timestamp mismatch: got %d, expected %d", decoded.timestamp, timestamp)
	}

	if decoded.randomness.Cmp(randomness) != 0 {
		t.Errorf("Randomness mismatch: got %s, expected %s", decoded.randomness.String(), randomness.String())
	}
}

func TestULIDMonotonicity(t *testing.T) {
	timestamp := uint64(time.Now().UnixMilli())

	ulid1, err := NewTime(timestamp)
	if err != nil {
		t.Fatalf("Error generating ULID 1: %v", err)
	}

	ulid2, err := NewTime(timestamp)
	if err != nil {
		t.Fatalf("Error generating ULID 2: %v", err)
	}

	if ulid2 <= ulid1 {
		t.Errorf("Monotonicity failed: ULID 2 is not greater than ULID 1")
	}
}

func TestULIDInvalidParsing(t *testing.T) {
	_, err := Parse("invalid-ulid-string")
	if err == nil {
		t.Errorf("Expected error for invalid ULID string")
	}

	_, err = Parse("0123456789ABCDEFGHJKMNPQRSTUV") // incorrect length
	if err == nil {
		t.Errorf("Expected error for invalid ULID string length")
	}
}

func TestULIDGetTime(t *testing.T) {
	now := time.Now()
	timestamp := uint64(now.UnixMilli())
	ulidStr, err := NewTime(timestamp)
	if err != nil {
		t.Fatalf("Error generating ULID: %v", err)
	}
	parsedUlid, err := Parse(ulidStr)
	if err != nil {
		t.Fatalf("Error parsing ULID: %v", err)
	}

	parsedTime := time.UnixMilli(int64(parsedUlid.GetTime()))

	if !parsedTime.Equal(now.Truncate(time.Millisecond)) {
		t.Errorf("GetTime failed: got %v, expected %v", parsedTime, now.Truncate(time.Millisecond))
	}
}

func TestTimestampOverflow(t *testing.T) {
	_, err := NewTime(maxTimestamp + 1)
	if err == nil {
		t.Errorf("Expected error for timestamp overflow")
	}
}

func TestRandomnessOverflow(t *testing.T) {

	mutex.Lock()
	lastTime = maxTimestamp // Set lastTime to max timestamp
	lastRandomness = new(big.Int).Set(maxRandomness)
	mutex.Unlock()

	_, err := NewTime(maxTimestamp) // Call NewTime with max timestamp
	if err == nil {
		t.Errorf("Expected error for randomness overflow")
	}
}
