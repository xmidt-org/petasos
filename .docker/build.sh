GOOS=linux GOARCH=amd64 go build ../src/petasos/

docker build -t petasos:local .
