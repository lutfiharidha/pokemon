tidy:
	@go mod tidy
battle: tidy
	@go run main.go battle
pokedex-s: tidy
	@go run main.go pokedex-s
pokedex: tidy
	@go run main.go pokedex
log-i: tidy
	@go run main.go log-i
log-d: tidy
	@go run main.go log-d
log-si: tidy
	@go run main.go log-si
help: tidy
	@go run main.go help