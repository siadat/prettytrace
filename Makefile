test:
	go test -count=1 -v ./...
	go run ./examples > /dev/null
	@echo "Done"
