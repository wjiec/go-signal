# go-signal
[![GoDoc](https://godoc.org/github.com/wjiec/go-signal?status.svg)](https://godoc.org/github.com/wjiec/go-signal)

Package signal provides simple, semantic manipulation of the operating system's signal processing.


### Installation

```bash
go get -u github.com/wjiec/go-signal
```


### Quick Start

Listens to the user's signal to exit the program and performs cleanup
```go
func main() {
	f, _ := os.Open("path/to/your/config")
	signal.Once(syscall.SIGTERM).Notify(context.TODO(), func(sig os.Signal) {
		_ = f.Close()
	})
}
```

Listening for `SIGUSR1` signals from users and performing services reload
```go
var srv Reloadable

func main() {
	signal.When(syscall.SIGUSR1).Notify(context.TODO(), func(sig os.Signal) {
		_ = srv.Reload()
	})
}
```

Create a context object using the specified signals and cancel the current context when the signal arrived
```go
var db *sql.DB

func main() {
	ctx, cancel := signal.With(context.TODO(), syscall.SIGTERM)
	defer cancel()

	_, _ = db.QueryContext(ctx, "select id,username,password from `user`")
}
```


### License

Released under the [MIT License](LICENSE).
