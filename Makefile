.PHONY:  prepare 

DST=$(wildcard typex/*_json.go)
SRC=$(patsubst %_json.go,%.go,$(DST))

all:  prepare

prepare: $(DST)

$(DST): $(SRC)
	@which easyjson >/dev/null || go install github.com/mailru/easyjson/...@latest
	@which stringer >/dev/null || go install golang.org/x/tools/cmd/stringer@latest
	@find typex -name "*_json.go" -delete
	@go generate ./...