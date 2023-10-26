module baran/grocer

go 1.20

replace baran/handler => ./handler
replace baran/model => ./model
replace baran/main => ./cmd

replace baran/database => ./database

require (
	baran/handler v0.0.0-00010101000000-000000000000
	baran/model v0.0.0-00010101000000-000000000000
	baran/main v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
)

require (
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
)
