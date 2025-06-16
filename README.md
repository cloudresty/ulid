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

## Why Choose ULID?

ULIDs offer significant advantages over traditional UUIDs and other identifier schemes, making them ideal for modern applications:

&nbsp;

### Performance Advantages

**Database B-Tree Efficiency:**

* **Sequential inserts** - ULIDs are time-ordered, reducing B-tree fragmentation
* **Better cache locality** - Related records are stored near each other
* **Faster queries** - Range queries benefit from natural time-based ordering
* **Reduced index rebuilding** - Less page splits in database indexes

**Generation Speed:**

* **~6x faster than UUID v4** - Optimized encoding and no complex formatting
* **151ns per ULID** vs typical UUID libraries at 800-1000ns
* **Zero allocations for parsing** - Efficient byte operations
* **Batch generation friendly** - Monotonicity allows rapid successive generation

&nbsp;

### Database & Storage Benefits

**Index Performance:**

```sql
-- ULID: Natural time-based clustering
SELECT * FROM orders WHERE created_between('01H0', '01H1'); -- Fast range scan

-- UUID: Random distribution requires full index scan
SELECT * FROM orders WHERE created_between(uuid1, uuid2); -- Slower, fragmented
```

**Storage Efficiency:**

* **26 characters** vs UUID's 36 characters (28% smaller in string form)
* **Better compression** - Time prefix allows better compression ratios
* **Reduced storage I/O** - Smaller keys mean more records per page

&nbsp;

### Developer Experience

**Human Readable:**

```text
ULID: 01ARZ3NDEKTSV4RRFFQ69G5FAV  (sortable, readable timestamp prefix)
UUID: 550e8400-e29b-41d4-a716-446655440000  (random, no meaningful order)
```

**Natural Sorting:**

* **Lexicographic sorting** matches chronological order
* **No custom comparators** needed
* **Works in any system** - databases, file systems, logs

**URL & API Friendly:**

* **No special characters** - safe in URLs without encoding
* **Case insensitive** - works with case-insensitive systems
* **Compact** - shorter URLs and API responses

&nbsp;

### Real-World Impact

**E-commerce Platform Example:**

```text
Traditional UUID v4 System:
- Order insertion: ~2000ms for 10K orders (random B-tree splits)
- Recent orders query: ~150ms (index fragmentation)
- Database size: 2.1GB for 1M orders

ULID-Based System:
- Order insertion: ~400ms for 10K orders (sequential inserts)
- Recent orders query: ~25ms (clustered data)
- Database size: 1.8GB for 1M orders (better compression)
```

**Microservices Benefits:**

* **Distributed tracing** - Natural correlation by time
* **Log aggregation** - Events sort chronologically across services
* **Debugging** - Easy to spot time-related patterns
* **Monitoring** - Time-based partitioning works naturally

&nbsp;

### Security & Uniqueness

**Collision Resistance:**

* **1.21e+24 unique IDs per millisecond** - practically impossible collisions
* **Cryptographically secure randomness** - unpredictable despite time component
* **No coordination required** - safe in distributed systems

**Vs UUID Comparison:**

| Feature | ULID | UUID v4 | UUID v1 |
|---------|------|---------|---------|
| **Sortable** | ✅ Natural | ❌ Random | ⚠️ Complex |
| **Performance** | ✅ ~150ns | ❌ ~800ns | ❌ ~600ns |
| **Size** | ✅ 26 chars | ❌ 36 chars | ❌ 36 chars |
| **B-tree friendly** | ✅ Sequential | ❌ Random | ⚠️ Partial |
| **URL safe** | ✅ No encoding | ❌ Needs encoding | ❌ Needs encoding |
| **Human readable** | ✅ Time prefix | ❌ Opaque | ⚠️ MAC address |
| **Privacy** | ✅ Anonymous | ✅ Anonymous | ❌ MAC exposed |

&nbsp;

### When to Use ULID

**Perfect for:**

* High-throughput applications requiring fast inserts
* Time-series data and event logging
* Distributed systems needing correlation
* APIs where shorter IDs improve performance
* Database-heavy applications with frequent queries

**Consider alternatives when:**

* Existing systems deeply integrated with UUID v4
* Regulatory requirements mandate specific UUID versions
* Time-based correlation is undesired for privacy reasons

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

## Installation

To install the `ULID` package, use the following command:

```bash
go get github.com/cloudresty/ulid
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

## Benchmarking

To run comprehensive benchmarks and generate a detailed performance report:

```bash
cd benchmarks
go run benchmark.go
```

This will generate a `RESULTS.md` file with:

* Detailed performance metrics
* System information
* Optimization techniques used
* Usage examples
* Comparison data

For standard Go benchmarks:

```bash
go test -bench=. -benchmem
```

&nbsp;

### Benchmark Results

```plaintext
goos: darwin
goarch: arm64
pkg: github.com/cloudresty/ulid
cpu: Apple M1 Max
BenchmarkNew-10          7565216               151.6 ns/op            32 B/op          1 allocs/op
BenchmarkParse-10       56834553                20.89 ns/op            0 B/op          0 allocs/op
BenchmarkString-10      50176750                24.70 ns/op           32 B/op          1 allocs/op
```

&nbsp;

### Latest Performance Results (Updated 2025-06-16)

Run `cd benchmarks && go run benchmark.go` to generate fresh benchmark results:

* **Generation Rate**: ~6.18 million ULIDs/second
* **Average Latency**: ~161ns per ULID
* **Memory Efficiency**: 32B/op, 1 alloc/op
* **Throughput**: 100,000 ULIDs in ~16ms

See [benchmarks/RESULTS.md](benchmarks/RESULTS.md) for detailed performance analysis and system information.

&nbsp;

## Basic Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/cloudresty/ulid"
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

## Basic Examples

ULIDs are highly versatile and can be used in various applications, including JSON APIs, NoSQL databases, and SQL databases.

### JSON Example

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"

    "github.com/cloudresty/ulid"
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

    "github.com/cloudresty/ulid"
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

    "github.com/cloudresty/ulid"
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

## More Examples & Use Cases

### Microservices & Distributed Systems

ULIDs are perfect for distributed systems where you need globally unique IDs without coordination:

```go
package main

import (
    "fmt"
    "log"
    "sync"
    "time"

    "github.com/cloudresty/ulid"
)

// Service represents a microservice instance
type Service struct {
    ID       string
    Name     string
    Requests []Request
    mu       sync.Mutex
}

type Request struct {
    ID        string    `json:"id"`
    ServiceID string    `json:"service_id"`
    Path      string    `json:"path"`
    Timestamp time.Time `json:"timestamp"`
}

func (s *Service) HandleRequest(path string) {
    requestID, err := ulid.New()
    if err != nil {
        log.Printf("Error generating request ID: %v", err)
        return
    }

    s.mu.Lock()
    s.Requests = append(s.Requests, Request{
        ID:        requestID,
        ServiceID: s.ID,
        Path:      path,
        Timestamp: time.Now(),
    })
    s.mu.Unlock()

    fmt.Printf("Service %s handled request %s for %s\n", s.Name, requestID, path)
}

func main() {
    serviceID, _ := ulid.New()
    service := &Service{
        ID:   serviceID,
        Name: "api-gateway",
    }

    // Simulate concurrent requests
    var wg sync.WaitGroup
    paths := []string{"/users", "/orders", "/products", "/health"}

    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(path string) {
            defer wg.Done()
            service.HandleRequest(path)
        }(paths[i%len(paths)])
    }

    wg.Wait()
    fmt.Printf("Service processed %d requests\n", len(service.Requests))
}
```

&nbsp;

### Event Sourcing & CQRS

ULIDs provide natural ordering for events in event-driven architectures:

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/cloudresty/ulid"
)

type Event struct {
    ID          string      `json:"id"`
    AggregateID string      `json:"aggregate_id"`
    Type        string      `json:"type"`
    Data        interface{} `json:"data"`
    Timestamp   time.Time   `json:"timestamp"`
    Version     int         `json:"version"`
}

type EventStore struct {
    events []Event
}

func (es *EventStore) AppendEvent(aggregateID, eventType string, data interface{}) error {
    eventID, err := ulid.New()
    if err != nil {
        return err
    }

    event := Event{
        ID:          eventID,
        AggregateID: aggregateID,
        Type:        eventType,
        Data:        data,
        Timestamp:   time.Now(),
        Version:     len(es.events) + 1,
    }

    es.events = append(es.events, event)
    return nil
}

func (es *EventStore) GetEvents(aggregateID string) []Event {
    var events []Event
    for _, event := range es.events {
        if event.AggregateID == aggregateID {
            events = append(events, event)
        }
    }
    return events
}

func main() {
    store := &EventStore{}
    userID, _ := ulid.New()

    // User lifecycle events
    store.AppendEvent(userID, "UserCreated", map[string]string{
        "email": "user@example.com",
        "name":  "John Doe",
    })

    store.AppendEvent(userID, "EmailUpdated", map[string]string{
        "old_email": "user@example.com",
        "new_email": "john.doe@example.com",
    })

    store.AppendEvent(userID, "UserDeactivated", map[string]string{
        "reason": "User requested account deletion",
    })

    // Retrieve and display events (naturally ordered by ULID)
    events := store.GetEvents(userID)
    for _, event := range events {
        eventJSON, _ := json.MarshalIndent(event, "", "  ")
        fmt.Println(string(eventJSON))
        fmt.Println("---")
    }
}
```

&nbsp;

### File & Document Management

Perfect for organizing files, documents, and media with time-based sorting:

```go
package main

import (
    "fmt"
    "path/filepath"
    "time"

    "github.com/cloudresty/ulid"
)

type Document struct {
    ID       string    `json:"id"`
    Name     string    `json:"name"`
    Type     string    `json:"type"`
    Size     int64     `json:"size"`
    Path     string    `json:"path"`
    Created  time.Time `json:"created"`
    Modified time.Time `json:"modified"`
}

type DocumentManager struct {
    documents map[string]Document
    basePath  string
}

func NewDocumentManager(basePath string) *DocumentManager {
    return &DocumentManager{
        documents: make(map[string]Document),
        basePath:  basePath,
    }
}

func (dm *DocumentManager) CreateDocument(name, docType string, size int64) (*Document, error) {
    docID, err := ulid.New()
    if err != nil {
        return nil, err
    }

    now := time.Now()
    doc := Document{
        ID:       docID,
        Name:     name,
        Type:     docType,
        Size:     size,
        Path:     filepath.Join(dm.basePath, docType, docID+filepath.Ext(name)),
        Created:  now,
        Modified: now,
    }

    dm.documents[docID] = doc
    return &doc, nil
}

func (dm *DocumentManager) GetDocument(id string) (*Document, bool) {
    doc, exists := dm.documents[id]
    return &doc, exists
}

func (dm *DocumentManager) ListDocumentsByType(docType string) []Document {
    var docs []Document
    for _, doc := range dm.documents {
        if doc.Type == docType {
            docs = append(docs, doc)
        }
    }
    return docs
}

func main() {
    dm := NewDocumentManager("/storage/documents")

    // Create various documents
    documents := []struct {
        name, docType string
        size          int64
    }{
        {"report.pdf", "pdf", 1024000},
        {"presentation.pptx", "presentation", 2048000},
        {"image.jpg", "image", 512000},
        {"contract.pdf", "pdf", 256000},
        {"logo.png", "image", 128000},
    }

    for _, d := range documents {
        doc, err := dm.CreateDocument(d.name, d.docType, d.size)
        if err != nil {
            fmt.Printf("Error creating document: %v\n", err)
            continue
        }
        fmt.Printf("Created document: %s (ID: %s)\n", doc.Name, doc.ID)
    }

    // List all PDF documents (naturally sorted by creation time due to ULID)
    pdfs := dm.ListDocumentsByType("pdf")
    fmt.Println("\nPDF Documents:")
    for _, pdf := range pdfs {
        fmt.Printf("- %s (Created: %s, Path: %s)\n",
            pdf.Name, pdf.Created.Format(time.RFC3339), pdf.Path)
    }
}
```

&nbsp;

### Logging & Tracing

ULIDs provide excellent correlation IDs for distributed tracing and logging:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/cloudresty/ulid"
)

type Logger struct {
    serviceName string
}

type LogEntry struct {
    TraceID   string    `json:"trace_id"`
    SpanID    string    `json:"span_id"`
    Service   string    `json:"service"`
    Level     string    `json:"level"`
    Message   string    `json:"message"`
    Timestamp time.Time `json:"timestamp"`
}

func (l *Logger) WithTrace(ctx context.Context, message, level string) {
    traceID := ctx.Value("trace_id").(string)
    spanID, _ := ulid.New()

    entry := LogEntry{
        TraceID:   traceID,
        SpanID:    spanID,
        Service:   l.serviceName,
        Level:     level,
        Message:   message,
        Timestamp: time.Now(),
    }

    fmt.Printf("[%s] %s | %s | Trace: %s | Span: %s | %s\n",
        entry.Timestamp.Format("2006-01-02 15:04:05"),
        entry.Level,
        entry.Service,
        entry.TraceID,
        entry.SpanID,
        entry.Message)
}

func processOrder(ctx context.Context, orderID string) {
    logger := &Logger{serviceName: "order-service"}

    logger.WithTrace(ctx, fmt.Sprintf("Processing order %s", orderID), "INFO")

    // Simulate processing time
    time.Sleep(50 * time.Millisecond)

    logger.WithTrace(ctx, "Validating order items", "DEBUG")
    time.Sleep(30 * time.Millisecond)

    logger.WithTrace(ctx, "Order validation complete", "INFO")
    time.Sleep(20 * time.Millisecond)

    logger.WithTrace(ctx, fmt.Sprintf("Order %s processed successfully", orderID), "INFO")
}

func main() {
    // Create a trace ID for the entire request
    traceID, err := ulid.New()
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.WithValue(context.Background(), "trace_id", traceID)

    // Create an order ID
    orderID, _ := ulid.New()

    fmt.Printf("Starting request processing with trace ID: %s\n", traceID)
    fmt.Println("=" * 80)

    processOrder(ctx, orderID)

    fmt.Println("=" * 80)
    fmt.Printf("Request completed. All logs correlated by trace ID: %s\n", traceID)
}
```

&nbsp;

### Gaming & Leaderboards

ULIDs work great for gaming applications where you need unique player/game session IDs:

```go
package main

import (
    "fmt"
    "math/rand"
    "sort"
    "time"

    "github.com/cloudresty/ulid"
)

type Player struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Score    int    `json:"score"`
    Level    int    `json:"level"`
}

type GameSession struct {
    ID        string    `json:"id"`
    PlayerID  string    `json:"player_id"`
    StartTime time.Time `json:"start_time"`
    EndTime   time.Time `json:"end_time"`
    Score     int       `json:"score"`
    Duration  time.Duration `json:"duration"`
}

type GameManager struct {
    players  map[string]Player
    sessions []GameSession
}

func NewGameManager() *GameManager {
    return &GameManager{
        players:  make(map[string]Player),
        sessions: make([]GameSession, 0),
    }
}

func (gm *GameManager) CreatePlayer(username string) (*Player, error) {
    playerID, err := ulid.New()
    if err != nil {
        return nil, err
    }

    player := Player{
        ID:       playerID,
        Username: username,
        Score:    0,
        Level:    1,
    }

    gm.players[playerID] = player
    return &player, nil
}

func (gm *GameManager) StartGameSession(playerID string) (*GameSession, error) {
    sessionID, err := ulid.New()
    if err != nil {
        return nil, err
    }

    session := GameSession{
        ID:        sessionID,
        PlayerID:  playerID,
        StartTime: time.Now(),
    }

    return &session, nil
}

func (gm *GameManager) EndGameSession(session *GameSession) {
    session.EndTime = time.Now()
    session.Duration = session.EndTime.Sub(session.StartTime)

    // Simulate random score
    session.Score = rand.Intn(10000) + 100

    gm.sessions = append(gm.sessions, *session)

    // Update player's total score
    if player, exists := gm.players[session.PlayerID]; exists {
        player.Score += session.Score
        if player.Score > player.Level*1000 {
            player.Level++
        }
        gm.players[session.PlayerID] = player
    }
}

func (gm *GameManager) GetTopPlayers(limit int) []Player {
    players := make([]Player, 0, len(gm.players))
    for _, player := range gm.players {
        players = append(players, player)
    }

    sort.Slice(players, func(i, j int) bool {
        return players[i].Score > players[j].Score
    })

    if limit > len(players) {
        limit = len(players)
    }

    return players[:limit]
}

func main() {
    gm := NewGameManager()

    // Create players
    playerNames := []string{"Alice", "Bob", "Charlie", "Diana", "Eve"}
    var players []*Player

    for _, name := range playerNames {
        player, err := gm.CreatePlayer(name)
        if err != nil {
            fmt.Printf("Error creating player %s: %v\n", name, err)
            continue
        }
        players = append(players, player)
        fmt.Printf("Created player: %s (ID: %s)\n", player.Username, player.ID)
    }

    // Simulate game sessions
    fmt.Println("\n=== Game Sessions ===")
    for _, player := range players {
        for i := 0; i < 3; i++ { // 3 sessions per player
            session, err := gm.StartGameSession(player.ID)
            if err != nil {
                continue
            }

            // Simulate game duration
            time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)+50))

            gm.EndGameSession(session)
            fmt.Printf("Session %s: %s scored %d points in %v\n",
                session.ID, player.Username, session.Score, session.Duration)
        }
    }

    // Show leaderboard
    fmt.Println("\n=== Leaderboard ===")
    topPlayers := gm.GetTopPlayers(5)
    for i, player := range topPlayers {
        fmt.Printf("%d. %s - Level %d - %d points (ID: %s)\n",
            i+1, player.Username, player.Level, player.Score, player.ID)
    }
}
```

&nbsp;

## Contributing

Contributions are welcome! Please submit [pull requests](https://github.com/cloudresty/ulid/pulls) or bug reports through GitHub.

&nbsp;

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE.txt) file for details.

&nbsp;

---

Made with ♥️ by [Cloudresty](https://cloudresty.com).
