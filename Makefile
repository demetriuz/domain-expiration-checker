.PHONY: bin

bin:
	docker build -f Dockerfile.artifacts -t domain-exp-checker.bin .
	- @docker rm -f domain-exp-checker.bin 2>/dev/null || exit 0
	docker run -d --name=domain-exp-checker.bin domain-exp-checker.bin
	docker cp domain-exp-checker.bin:/artifacts/domain-expiration-checker.linux-amd64 domain-expiration-checker.linux-amd64
	docker cp domain-exp-checker.bin:/artifacts/domain-expiration-checker.windows-amd64.exe domain-expiration-checker.windows-amd64.exe
	docker cp domain-exp-checker.bin:/artifacts/domain-expiration-checker.darwin-amd64 domain-expiration-checker.darwin-amd64