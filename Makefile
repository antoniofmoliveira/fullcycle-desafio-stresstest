testpost:
	go run cmd/main.go --numtests 400000 -requesttype "POST" -payload '{"message":"World"}' -endpoint "http://localhost:8080/hello"

testget:
	go run cmd/main.go --numtests 300000 -requesttype "GET" -endpoint "http://localhost:8080/hello"