.PHONY:blockchain clean

app: clean blockchain

#编译命令
blockchain:
	go build -o blockchain *.go

clean:
	rm blockchain *.db

test:
	go test -v *_test.go

docker:
	#docker build xxx