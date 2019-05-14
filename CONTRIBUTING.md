# Contributions
Contributions are always welcome. If there is a bug or a feature that you would like to address please feel free to tackle it.

If you fork the project, take a look at [Working with forks in Go](http://blog.sgmansfield.com/2016/06/working-with-forks-in-go) which is a great article about how you can get around the importing problems that sometimes affect Go projects.

```bash
git clone https://github.com/kochie/zeit-api-go.git

cd zeit-api-go

git remote rename origin upstream

git remote add origin https://github.com/[username]/zeit-api-go
```

Also please run `golint`, `gofmt`, and `go test` before submitting a pull request. It might be usefult to create a git hook for this.
