# netcat
A more friendly and simple netcat, implemented using golang。

```bash
go install github.com/10khour/netcat@latest
```

Send file
```
tmp $ netcat --port 9999 --host 10.80.184.2 < bigfile
Count: 208 MB Speed: 1.9 MB/S
```
Recevie file
```
tmp $ netcat -l 9999 > bigfile
Count: 526 MB Speed: 298 MB/S
```
