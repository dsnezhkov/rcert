# RCert
Used when reverse IP lookup does not help to find the name of the server and the domain ( e.g. for virtual hosting), and there is an SSL socket listening on the remote server.

Sometimes it is possible to get additonal clues from the SSL certificate's extended Subject Alternative Name (SAN) records.

Things *rcert* looks for:
- DNS names
- Email addresses
- Permitted domains

All these are presented visually for inspection. Also, greppable format. 

### Usage:
$ `go run rcert.go -ipfile=./randips.open443`

or 

$ `go build -o rcert.osx`

$ `./rcert.osx   -ipfile=./randips.open443 `
### X-compile:


bash-3.2$ `GOOS=darwin GOARCH=386 go build -o rcert.osx`

bash-3.2$ `GOOS=linux GOARCH=386 go build -o rcert.ux`

bash-3.2$ `GOOS=windows GOARCH=386 go build -o rcert.exe`

![rcert](/relative/path/to/img.jpg?raw=true "rcert run")
