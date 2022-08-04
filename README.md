# Description
simple crud application for moveis written in golang
# Build Application
```bash
make build
```
# Start Application
you should pass database uri to app using env variables.
```bash
POSTGRES_URI=URI
```
uri format: 
```bash
postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]
```
you can also pass a port which you want app runs on it (default port is 8080):
```bash
./movie_app -port [port]
```