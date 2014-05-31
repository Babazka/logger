Logger
======

Simple HTTP logger which writes HTTP queries to text files.

Usage
-----
Compile:
```
logger$ make
go install logger
```

Run:
```
logger$ bin/logger --listen=":5588" --log-root=/www/logger >/dev/null 2>&1 &
[1] 5400
logger$ curl "http://localhost:5588/foo/bar?baz=123&wtf=42"
OK
logger$ curl -d "h123=331&jjjj=oooo" http://localhost:5588/foo/bar?buu
OK
logger$ tail /www/logger/foo/bar
2014-05-31T20:06:53.489 GET /foo/bar?baz=123&wtf=42

.
2014-05-31T20:07:34.441 POST /foo/bar?buu
h123=331&jjjj=oooo
.
```
