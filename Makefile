obu: 
	@go build -o bin/obu obu/main.go
	@echo "\033[32m[+] Build obu successfully\033[0m"
	@./bin/obu

receiver:
	@go build -o bin/receiver ./receiver
	@echo "\033[32m[+] Build receiver successfully\033[0m"
	@./bin/receiver

calculator:
	@go build -o bin/calculator ./calculator
	@echo "\033[32m[+] Build calculator successfully\033[0m"
	@./bin/calculator

.PHONY: obu , receiver , calculator