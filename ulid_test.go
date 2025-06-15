package ulid

import (
	"testing"
	"time"
)

func TestULIDEncodingDecoding(t *testing.T) {
	timestamp := uint64(time.Now().UnixMilli())
	var randomness [randomnessBytes]byte
	for i := range randomness {
		randomness[i] = byte(i * 13) // Some test pattern
	}

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

	if decoded.randomness != randomness {
		t.Errorf("Randomness mismatch: got %v, expected %v", decoded.randomness, randomness)
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
	// Set lastRandomness to maximum value (all 0xFF)
	for i := range lastRandomness {
		lastRandomness[i] = 0xFF
	}
	mutex.Unlock()

	_, err := NewTime(maxTimestamp) // Call NewTime with max timestamp
	if err == nil {
		t.Errorf("Expected error for randomness overflow")
	}
}

// Benchmark functions
func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = New()
	}
}

func BenchmarkParse(b *testing.B) {
	ulidStr, _ := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Parse(ulidStr)
	}
}

func BenchmarkString(b *testing.B) {
	ulid, _ := Parse("01ARZ3NDEKTSV4RRFFQ69G5FAV")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ulid.String()
	}
}
