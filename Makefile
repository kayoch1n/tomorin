build: clean
	CGO_ENABLED=0 go build -o ./bin/tomorin

clean:
	rm bin -rf
