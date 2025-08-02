module github.com/SteinerLabs/lms/backend/services/auth

go 1.24.4

require (
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/google/uuid v1.6.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.37.0
	google.golang.org/grpc v1.62.0
	google.golang.org/protobuf v1.32.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240123012728-ef4313101c80 // indirect
)

replace github.com/SteinerLabs/lms/backend/shared => ../../shared
