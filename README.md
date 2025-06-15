# ULID (Universally Unique Lexicographically Sortable Identifier) for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/cloudresty/ulid.svg)](https://pkg.go.dev/github.com/cloudresty/ulid)
[![Go Tests](https://github.com/cloudresty/ulid/actions/workflows/test.yaml/badge.svg)](https://github.com/cloudresty/ulid/actions/workflows/test.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/cloudresty/ulid)](https://goreportcard.com/report/github.com/cloudresty/ulid)
[![GitHub Tag](https://img.shields.io/github/v/tag/cloudresty/ulid?label=Version)](https://github.com/cloudresty/ulid/tags)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

&nbsp;

This package provides a robust and efficient Go implementation of the ULID (Universally Unique Lexicographically Sortable Identifier) specification, as defined in [github.com/ulid/spec](https://github.com/ulid/spec). ULIDs are designed to be universally unique, lexicographically sortable, and more compact than UUIDs.

&nbsp;

## Key Features

* **128-bit Compatibility with UUID:** Seamless integration with systems that use UUIDs.
* **High Throughput:** Generates 1.21e+24 unique ULIDs per millisecond, suitable for high-demand applications.
* **Lexicographical Sortability:** Enables efficient sorting and indexing in databases and other systems.
* **Compact Representation:** Encoded as a 26-character string using Crockford's Base32, compared to the 36-character UUID.
* **Crockford's Base32 Encoding:** Improves readability and efficiency by excluding ambiguous characters (I, L, O, U).
* **Lowercase by Default:** New in v1.1+ - generates lowercase ULIDs for better readability while maintaining case-insensitive parsing.
* **Case Insensitive Parsing:** Accepts both uppercase and lowercase ULIDs for backward compatibility.
* **URL Safety:** Contains no special characters, making it safe for use in URLs and web applications.
* **Monotonicity:** Ensures correct sorting order even when multiple ULIDs are generated within the same millisecond.
* **Thread Safety:** Safe for concurrent use in multi-threaded applications.
* **High Performance:** Optimized implementation with ~10-50x performance improvements over previous versions.
* **Zero Dependencies:** No external dependencies beyond Go standard library.

&nbsp;

## ULID Structure

A ULID consists of two components:

```text
 01AN4Z07BY     79KA1307SR9X4MV3

|-----------|  |----------------|
  Timestamp        Randomness
    48bits           80bits
```

* **Timestamp (48 bits):** Represents the UNIX timestamp in milliseconds, allowing for time-based sorting and uniqueness. This component provides time representation up to the year 10889 AD.
* **Randomness (80 bits):** A cryptographically secure random value that ensures uniqueness even within the same millisecond.

&nbsp;

## Performance

This implementation is highly optimized for speed and efficiency:

* **~10-50x faster** than big.Int-based implementations
* **~4.5M+ ULIDs/second** generation rate on modern hardware
* **~200ns per ULID** average generation time
* **Zero external dependencies** - only Go standard library
* **Reduced memory allocations** through optimized byte operations
* **Custom Crockford Base32** encoding for maximum performance
* **Efficient monotonicity** handling with byte array operations

### Benchmark Results

```plaintext
BenchmarkNew-10                 5000000    213 ns/op    48 B/op    2 allocs/op
BenchmarkParse-10              10000000    120 ns/op    32 B/op    1 allocs/op
BenchmarkString-10             20000000     85 ns/op    32 B/op    1 allocs/op
```

&nbsp;

## Installation

To install the `ULID` package, use the following command:

```bash
go get github.com/cloudresty/ulid
```

&nbsp;

## Usage

```go
package main

import (
    "fmt"
    "log"

    ulid "github.com/cloudresty/ulid"
)

func main() {

    // Generate a new ULID
    ulidStr, err := ulid.New()
    if err != nil {
        log.Fatalf("Error generating ULID: %v", err)
    }
    fmt.Println("Generated ULID:", ulidStr)

    // Parse a ULID string
    parsedUlid, err := ulid.Parse(ulidStr)
    if err != nil {
        log.Fatalf("Error parsing ULID: %v", err)
    }
    fmt.Println("Parsed ULID time:", parsedUlid.GetTime())

    // Generate a ULID with a specific timestamp
    timestamp := uint64(1678886400000) // Example timestamp (milliseconds)
    ulidStr2, err := ulid.NewTime(timestamp)
    if err != nil {
        log.Fatalf("Error generating ULID with time: %v", err)
    }
    fmt.Println("ULID with specific timestamp:", ulidStr2)

}
```

&nbsp;

## Examples

ULIDs are highly versatile and can be used in various applications, including JSON APIs, NoSQL databases, and SQL databases.

### JSON Example

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"

    ulid "github.com/cloudresty/ulid"
)

type User struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

func main() {

    ulidStr, err := ulid.New()
    if err != nil {
        log.Fatalf("Error generating ULID: %v", err)
    }

    user := User{
        ID:   ulidStr,
        Name: "John Doe",
    }

    userJSON, err := json.Marshal(user)
    if err != nil {
        log.Fatalf("Error marshaling JSON: %v", err)
    }

    fmt.Println(string(userJSON))

}
```

&nbsp;

### MongoDB (NoSQL) Example

When using MongoDB, you can store ULIDs as strings. MongoDB's indexing and sorting capabilities will work seamlessly with ULIDs.

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

    ulid "github.com/cloudresty/ulid"
)

type Product struct {
    ID   string `bson:"_id"`
    Name string `bson:"name"`
}

func main() {

    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }
    defer func() {
        if err = client.Disconnect(context.TODO()); err != nil {
            panic(err)
        }
    }()

    collection := client.Database("testdb").Collection("products")

    ulidStr, err := ulid.New()
    if err != nil {
        log.Fatalf("Error generating ULID: %v", err)
    }

    product := Product{
        ID:   ulidStr,
        Name: "Laptop",
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err = collection.InsertOne(ctx, product)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Product inserted with ID:", product.ID)

    // Find the product
    var foundProduct Product
    err = collection.FindOne(ctx, bson.M{"_id": product.ID}).Decode(&foundProduct)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Found product:", foundProduct)

}
```

&nbsp;

### PostgreSQL (SQL) Example

ULIDs can also be used as primary keys in SQL databases like PostgreSQL. You can store them as `VARCHAR(26)` columns.

```go
package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq" // PostgreSQL driver

    ulid "github.com/cloudresty/ulid"
)

type Order struct {
    ID     string
    UserID int
    Amount float64
}

func main() {

    connStr := "user=postgres password=password dbname=testdb sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    ulidStr, err := ulid.New()
    if err != nil {
        log.Fatalf("Error generating ULID: %v", err)
    }

    order := Order{
        ID:     ulidStr,
        UserID: 123,
        Amount: 99.99,
    }

    _, err = db.Exec("CREATE TABLE IF NOT EXISTS orders (id VARCHAR(26) PRIMARY KEY, user_id INTEGER, amount FLOAT)")
    if err != nil {
        log.Fatal(err)
    }

    _, err = db.Exec("INSERT INTO orders (id, user_id, amount) VALUES ($1, $2, $3)", order.ID, order.UserID, order.Amount)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Order inserted with ID:", order.ID)

    // Find the order
    var foundOrder Order
    err = db.QueryRow("SELECT id, user_id, amount FROM orders WHERE id = $1", order.ID).Scan(&foundOrder.ID, &foundOrder.UserID, &foundOrder.Amount)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Found order:", foundOrder)

}
```

&nbsp;

## API Reference

**`func New() (string, error)`**

Generates a new ULID string using the current UNIX timestamp in milliseconds.

```go
ulidStr, err := ulid.New()
if err != nil {
    // Handle error
}
```

&nbsp;

**`func NewTime(timestamp uint64) (string, error)`**

Generates a new ULID string using the provided UNIX timestamp in milliseconds.

```go
ulidStr, err := ulid.NewTime(1678886400000)
if err != nil {
    // Handle error
}
```

&nbsp;

**`func Parse(s string) (ULID, error)`**

Parses a ULID string and returns a `ULID` struct. Returns an error if the string is invalid.

```go
parsedUlid, err := ulid.Parse("01ARZ3NDEKTSV4RRFFQ69G5FAV")
if err != nil {
    // Handle error
}
```

&nbsp;

**`func (u ULID) String() string`**

Returns the canonical 26-character string representation of the `ULID`.

```go
ulidStr := parsedUlid.String()
```

&nbsp;

## Error Handling

The package returns errors for:

* Invalid ULID string formats.
* Timestamps exceeding the maximum allowed value.
* Randomness generation failures.
* Randomness overflow during monotonic generation.

&nbsp;

## Thread Safety

The `New()` and `NewTime()` functions are thread-safe, ensuring safe concurrent use.

&nbsp;

## Monotonicity Details

When generating ULIDs within the same millisecond, the package ensures monotonicity by incrementing the randomness component. If the randomness component reaches its maximum value, the timestamp is incremented.

&nbsp;

## Contributing

Contributions are welcome! Please submit [pull requests](https://github.com/cloudresty/ulid/pulls) or bug reports through GitHub.

&nbsp;

## License

This project is licensed under the MIT License. See the LICENSE file for details.

&nbsp;

---

Made with ♥️ by [Cloudresty](https://cloudresty.com).
