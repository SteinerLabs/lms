module github.com/SteinerLabs/lms/backend/services/user

go 1.24.4

require (
	github.com/SteinerLabs/lms/backend/shared v0.0.0
	github.com/jackc/pgx/v5 v5.7.5
	github.com/nats-io/nats.go v1.44.0
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/nats-io/nkeys v0.4.11 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
)

replace github.com/SteinerLabs/lms/backend/shared => ../../shared
