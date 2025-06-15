module benchmarks

go 1.24.1

toolchain go1.24.4

replace github.com/cloudresty/ulid => ../

require (
	github.com/cloudresty/ulid v0.0.0
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/oklog/ulid/v2 v2.1.0
)
