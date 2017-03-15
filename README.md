# RCert
Used when reverse IP lookup does not help to find the name of the server and the domain ( e.g. for virtual hosting), and there is an SSL socket listening on the remote server.

Sometimes it is possible to get additonal clues from the SSL certificate's extended Subject Alternative Name (SAN) records.

Things *rcert* looks for:
- DNS names
- Email addresses
- Permitted domains

All these are presented visually for inspection. Also, greppable format. 

Timeouts on TLS connection are decent but the overall speed is not great at the moment - sequential connectivity.
So if you want to pre-scan the hosts to see if they are listening on HTTP/S that will speed things up.

The tool does *not* fetch content of pages it requests. It breaks after the TLS has been established and the certificate can be examined.

* TODO: Async connectivity.

### Usage:
Accepts a file with one IP per line

$ `go run rcert.go -ipfile=./randips.open443`

or 

$ `go build -o rcert.osx`

$ `./rcert.osx   -ipfile=./randips.open443 `
### Cross-compile for Go:

bash-3.2$ `GOOS=darwin GOARCH=386 go build -o rcert.osx`

bash-3.2$ `GOOS=linux GOARCH=386 go build -o rcert.ux`

bash-3.2$ `GOOS=windows GOARCH=386 go build -o rcert.exe`

* Provided binaries:
rcert.exe: PE32 executable (console) Intel 80386 (stripped to external PDB), for MS Windows
rcert.osx: Mach-O executable i386
rcert.ux:  ELF 32-bit LSB executable, Intel 80386, version 1 (SYSV), statically linked, not stripped

![rcert](https://github.com/dsnezhkov/rcert/raw/master/screenshot_448.png "rcert run")
