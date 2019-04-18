all: build package

clean:
	rm -v generator | rm -v main main.zip

embed:
	embedfiles -pkg main files/

build:
	env GOOS=linux GOARCH=amd64 go build -o main .

package: build
	zip -j main.zip main
