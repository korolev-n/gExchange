module github.com/korolev-n/gExchange/wallet

go 1.24.4

require (
	github.com/go-chi/chi/v5 v5.2.1
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/joho/godotenv v1.5.1
	github.com/korolev-n/gExchange/exchanger v0.0.0-00010101000000-000000000000
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.39.0
	google.golang.org/grpc v1.73.0
)

require (
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace github.com/korolev-n/gExchange/exchanger => ../exchanger
