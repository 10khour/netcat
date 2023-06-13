# netcat
Netcat write in go , display network progress .

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