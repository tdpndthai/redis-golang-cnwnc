### Download Go:
https://golang.org/dl/go1.15.windows-amd64.msi

### Using go modules:
```
go mod init <name_modules>
```

### Build project
```
go build
```

### Run project
```
go run main.go
```
or
```
./<name_modules> (./redis-golang-cnwnc)
```

### Server redis
```
RUN docker pull redis
RUN docker run --name my-redis -p 6379:6379 -d redis
```

### Using package  
https://github.com/go-redis/redis
