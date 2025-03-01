exe: src/main.go
	go build src/main.go

all: exe
	exe
	touch ~/.config/netBr/config.json

clean: 
	rm ./main
	rm ~/.config/netBr/config.json

