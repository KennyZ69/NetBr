exe: src/main.go
	go build src/main.go

all: 
	exe

clean: 
	rm ./main
	rm ~/.config/netBr/config.json

