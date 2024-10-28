testpost:
	cd tester && go run cmd/main.go --numtests 300000 -requesttype "POST" -payload '{"message":"World"}' -endpoint "http://localhost:8080/hello"

testget:
	cd tester && go run cmd/main.go --numtests 300000 -requesttype "GET" -endpoint "http://localhost:8080/hello"


serverwitherror:
	cd server && go run server/cmd/main.go --qt-tokens 100000 --time-frame-seconds 1 --simulate-slow-requests --seed-for-simulate-slow-requests 100 --simulate-errors  --seed-for-simulate-errors 100

serverwithouterror:
	cd server && go run cmd/main.go --qt-tokens 100000 --time-frame-seconds 1