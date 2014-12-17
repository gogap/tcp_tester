
TCP TESTER
==========
Tcp tester is a simple tool for send some test data to tcp server for mock a request, it could be a validator of the server protocol.

### Usage

```bash
$ go run main.go -f test.dat -n 2
[00001]: pong
[00002]: pong
```


```bash
go run main.go --help
Usage of tcp_tester
  -b=1024: the buffer size (shorthand)
  -buffer=1024: the buffer size
  -f="": the file to send to tcp server (shorthand)
  -file="": the file to send to tcp server
  -i="127.0.0.1": server ip (shorthand)
  -ip="127.0.0.1": server ip
  -o="": the content response form server will write this file (shorthand)
  -output="": the content response form server will write this file
  -p=6788: server port (shorthand)
  -port=6788: server port
  -t=100: timeout to wait server response (shorthand)
  -timeout=100: timeout to wait server response
  -times=1: test times
  -n=1: test times (shorthand)
```

