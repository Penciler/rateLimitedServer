# Run server
```bash
cd src
./main_linux or ./main_windows (tested on ubuntu 16.04)
```
# Test
```bash
cd src
go test -run Main
```
# Build
```bash
cd src
go build -o $AppName
```
# Introduction
This is a server with simple rate limiting capability. Leaky bucket algorithm is used to limit request from different IP address.
The default setting is 60 requests per minute and can handle burst size of 60 requests. 
An IP to limter map is used to track request history of each IP. The data is in memory for quick access. And since the record older than 1 minute is irrelevant for default request limit, we don't need a database for permanent storage. There is also a routine that clean record older than 1 minute.
The test case would send request with two IPs:

192.168.12.1,  60 request, 10 millisecond interval. The purpose is to test the upper limit of the server.  
192.168.12.2,  61 request, 10 millisecond interval. The purpose is to test if the server report error when over limit.

# References 
How to Rate Limit HTTP Requests:
https://www.alexedwards.net/blog/how-to-rate-limit-http-requests

Go by Example: Rate Limiting:
https://gobyexample.com/rate-limiting