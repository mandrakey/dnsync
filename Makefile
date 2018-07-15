all : amd64 386

amd64 :
	GOARCH=amd64 go build -o ./bin/dnsync_amd64

386 :
	GOARCH=386 go build -o ./bin/dnsync_386

test :
	go test ./...

clean :
	rm bind/bindconfig_test2.conf
	rm -rf bin
