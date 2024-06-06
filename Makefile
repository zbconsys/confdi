run: clean merge

clean:
	rm merged.toml | echo ""

merge:
	go run main.go -default ./config-example/default.toml -override ./config-example/override.toml -merged merged.toml

test:
	go clean -testcache && go test -v ./...