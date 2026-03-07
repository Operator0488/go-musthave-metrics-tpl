.PHONY: test1 test2 test3 test4 test5 test6 test7 test8 test9
test1:
	./metricstest-darwin-arm64 -test.v -test.run=^TestIteration1$$ \
      -binary-path=cmd/server/server \
      -source-path=.

test2:
	./metricstest-darwin-arm64 -test.v -test.run=^TestIteration2 \
	  -agent-binary-path=cmd/agent/agent \
      -binary-path=cmd/server/server \
      -source-path=.

test3:
	./metricstest-darwin-arm64 -test.v -test.run=^TestIteration3 \
	  -agent-binary-path=cmd/agent/agent \
      -binary-path=cmd/server/server \
      -source-path=.

test4:
	./metricstest-darwin-arm64 -test.v -test.run=^TestIteration4$ \
      -agent-binary-path=cmd/agent/agent \
      -binary-path=cmd/server/server \
      -server-port=8080 \
      -source-path=.

test5:
	./metricstest-darwin-arm64 -test.v -test.run=^TestIteration5$ \
      -agent-binary-path=cmd/agent/agent \
      -binary-path=cmd/server/server \
      -server-port=8080 \
      -source-path=.

test6:
	./metricstest-darwin-arm64 -test.v -test.run=^TestIteration6$ \
      -agent-binary-path=cmd/agent/agent \
      -binary-path=cmd/server/server \
      -server-port=8080 \
      -source-path=.

test7:
	./metricstest-darwin-arm64 -test.v -test.run=^TestIteration7$ \
    	-agent-binary-path=cmd/agent/agent \
        -binary-path=cmd/server/server \
      	-server-port=8080 \
        -source-path=.

test8:
	./metricstest-darwin-arm64 -test.v -test.run=^TestIteration8$ \
        	-agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
          	-server-port=8080 \
            -source-path=.

test9:
	@go build -o cmd/agent/agent cmd/agent/*.go
	@go build -o cmd/server/server cmd/server/*.go
	./metricstest-darwin-arm64 -test.v -test.run=^TestIteration9$ \
    	-agent-binary-path=cmd/agent/agent \
        -binary-path=cmd/server/server \
        -file-storage-path=./file \
      	-server-port=8080 \
        -source-path=.

test10:
	@docker-compose up -d
	./metricstest-darwin-arm64 -test.v -test.run=^TestIteration10[AB]$ \
    	-agent-binary-path=cmd/agent/agent \
        -binary-path=cmd/server/server \
        -file-storage-path=./file \
        -database-dsn='postgres://postgres:yourpasswords@localhost:5432/postgres?sslmode=disable' \
      	-server-port=8080 \
        -source-path=.

test11:
	@docker-compose up -d
	@go build -o cmd/agent/agent cmd/agent/*.go
	@go build -o cmd/server/server cmd/server/*.go
	./metricstest-darwin-arm64 -test.v -test.run=^TestIteration11$ \
    	-agent-binary-path=cmd/agent/agent \
        -binary-path=cmd/server/server \
        -database-dsn='postgres://postgres:yourpasswords@localhost:5432/postgres?sslmode=disable' \
      	-server-port=8080 \
        -source-path=.

test12:
	@docker-compose up -d
	@go build -o cmd/agent/agent cmd/agent/*.go
	@go build -o cmd/server/server cmd/server/*.go
	./metricstest-darwin-arm64 -test.v -test.run=^TestIteration12$ \
    	-agent-binary-path=cmd/agent/agent \
        -binary-path=cmd/server/server \
        -database-dsn='postgres://postgres:yourpasswords@localhost:5432/postgres?sslmode=disable' \
      	-server-port=8080 \
        -source-path=.

test13:
	@docker-compose up -d
	@go build -o cmd/agent/agent cmd/agent/*.go
	@go build -o cmd/server/server cmd/server/*.go
	./metricstest-darwin-arm64 -test.v -test.run=^TestIteration13$ \
    	-agent-binary-path=cmd/agent/agent \
        -binary-path=cmd/server/server \
        -database-dsn='postgres://postgres:yourpasswords@localhost:5432/postgres?sslmode=disable' \
      	-server-port=8080 \
        -source-path=.


test9-log:
	@mkdir -p logs
	@$(MAKE) test9 2>&1 | tee logs/test9.log

test8-log:
	@mkdir -p logs
	@$(MAKE) test8 2>&1 | tee logs/test8.log

test7-log:
	@mkdir -p logs
	@$(MAKE) test7 2>&1 | tee logs/test7.log

.PHONY: buildServer
buildServer:
	@go build -o cmd/server/server cmd/server/*.go

.PHONY: buildAgent
buildAgent:
	@go build -o cmd/agent/agent cmd/agent/*.go

.PHONY: test
test:
	@go test ./... -coverprofile cover.out && go tool cover -func cover.out && go tool cover -html cover.out

migration:
	@goose -dir migrations/postgres create metrics sql
