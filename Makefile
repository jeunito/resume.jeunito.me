all: build package

clean:
	rm -v main main.zip

build:
	env GOOS=linux GOARCH=amd64 go build -o main .

package: build
	zip -j main.zip main
	mv main.zip terraform/
