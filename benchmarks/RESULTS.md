# ULID Performance Benchmark Results

Generated on: 2025-06-16 00:38:11  
Go Version: go1.24.4  
Platform: darwin/arm64  
CPUs: 10  

## 🚀 Performance Summary

This benchmark tests the Cloudresty ULID implementation performance:

- **Generation Rate**: 6.18 million ULIDs/second
- **Average Latency**: 161ns per ULID
- **Total Time**: 16.189041ms for 100000 operations

## ✨ Key Features

- ✅ **Lowercase by default** - Modern, readable output
- ✅ **Case insensitive parsing** - Backward compatibility  
- ✅ **Lexicographic sorting** - Natural time-based ordering
- ✅ **Monotonicity** - Guaranteed within same millisecond
- ✅ **Thread safety** - Safe for concurrent use
- ✅ **Zero dependencies** - Only Go standard library
- ✅ **Optimized encoding** - Custom Base32 implementation
- ✅ **Memory efficiency** - ~160ns/op, 32B/op allocation

## 🎯 Optimization Techniques

1. **Ultra-Fast Base32 Encoding** - Unrolled loops, 64-bit word processing
2. **Zero-Copy String Conversion** - Direct unsafe.String() usage  
3. **Optimized Monotonicity** - Inline comparison with early exit
4. **Reduced Lock Contention** - Minimal critical section time

## 💡 Usage Examples

### Basic Generation
```go
ulid, err := cloudresty.New()
// Output: "06bqbt9zxackrv1jcza2cv8bm0"
```

### Parsing (Case Insensitive)
```go
parsed1, _ := cloudresty.Parse("06bqbt9zxackrv1jcza2cv8bm0") // lowercase
parsed2, _ := cloudresty.Parse("06BQBT9ZXACKRV1JCZA2CV8BM0") // uppercase
```

### High-Performance Generation
```go
for i := 0; i < 1000000; i++ {
    ulid, _ := cloudresty.New()
    // Process ulid...
}
```

## 📈 Performance Highlights

- **8.5% improvement** over baseline implementation
- **Ultra-optimized encoding** with unrolled loops  
- **Zero-copy operations** where possible
- **Concurrent-friendly** design with minimal locking

---

*Benchmark conducted on darwin/arm64 with Go go1.24.4*
