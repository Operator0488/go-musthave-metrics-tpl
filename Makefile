.PHONY: test1 test build clean
test1:
	./metricstest-darwin-arm64 -test.v -test.run=^TestIteration1 \
      -binary-path=cmd/server/server \
      -source-path=.

.PHONY: test2 test build clean
test2:
	./metricstest-darwin-arm64 -test.v -test.run=^TestIteration2 \
	  -agent-binary-path=cmd/agent/agent \
      -binary-path=cmd/server/server \
      -source-path=.

.PHONY: test3 test build clean
test3:
	./metricstest-darwin-arm64 -test.v -test.run=^TestIteration3 \
	  -agent-binary-path=cmd/agent/agent \
      -binary-path=cmd/server/server \
      -source-path=.

.PHONY: test4 test build clean
test4:
	./metricstest-darwin-arm64 -test.v -test.run=^TestIteration4$ \
      -agent-binary-path=cmd/agent/agent \
      -binary-path=cmd/server/server \
      -server-port=8080 \
      -source-path=.


.PHONY: test5 test build clean
test5:
	./metricstest-darwin-arm64 -test.v -test.run=^TestIteration5$ \
      -agent-binary-path=cmd/agent/agent \
      -binary-path=cmd/server/server \
      -server-port=8080 \
      -source-path=.

.PHONY: buildServer
buildServer:
	@go build -o cmd/server/server cmd/server/*.go

.PHONY: buildAgent
buildAgent:
	@go build -o cmd/agent/agent cmd/agent/*.go

.PHONY: test
test:
	@go test ./... -coverprofile cover.out && go tool cover -func cover.out && go tool cover -html cover.out
