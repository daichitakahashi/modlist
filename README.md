# modlist
`modlist` command lists your go modules/packages.

## How it works
`modlist` loads `go.work` and list Go modules. And load module names from each `go.mod` files of modules.

If `--packages` option is enabled, list all packages from your modules instead of module name.
It uses `go list` command

## Usage
```shell
$ cat go.work
go 1.20

use (
        .
        ./module1
        ./module2
)
$ tree .
.
├── go.mod
├── go.work
├── module1
│   ├── go.mod
│   ├── package1
│   │   └── source.go
│   ├── package2
│   │   └── source.go
│   └── source.go
├── module2
│   ├── go.mod
│   ├── internal
│   │   └── pkg
│   │       └── source.go
│   └── source.go
├── package1
│   └── source.go
├── package2
│   └── source.go
└── source.go

9 directories, 12 files
```

List your modules.
```shell
$ modlist
example.com/multi
example.com/multi/module1
example.com/multi/module2
```

List your packages.
```shell
$ modlist --packages
example.com/multi
example.com/multi/module1
example.com/multi/module1/package1
example.com/multi/module1/package2
example.com/multi/module2
example.com/multi/module2/internal/pkg
example.com/multi/package1
example.com/multi/package2
```

Test your all packages and get coverages without internal package.
```shell
$ go test -coverpkg=`modlist -p -e="*internal*" --separator=","` `modlist -p -s`
ok      example.com/multi/module2       0.115s  coverage: [no statements]
ok      example.com/multi/package1      0.221s  coverage: [no statements]
ok      example.com/multi/module1/package2      0.325s  coverage: [no statements]
ok      example.com/multi       0.429s  coverage: [no statements]
ok      example.com/multi/module1       0.815s  coverage: [no statements]
ok      example.com/multi/module2/internal/pkg  0.928s  coverage: [no statements]
ok      example.com/multi/package2      1.333s  coverage: [no statements]
ok      example.com/multi/module1/package1      1.446s  coverage: [no statements]
```

## Options
|option|shorthand|default|description|
|---|---|---|---|
|`--packages`|`-p`|false|List all packages instead of modules|
|`--shuffle`|`-s`|false|Shuffle items|
|`--match`|`-m`||Exclude unmatch items|
|`--exclude`|`-e`||Exclude match items|
|`--separator`||"\n"|Change separator|
|`--directory`|`-d`|false|Show module/package paths instead of their names|
|`--golangci-lint-skip-dirs`||false|If configuration file of [golangci-lint](https://golangci-lint.run/usage/configuration/) exists, read `run.skip-dirs` and `run.skip-dirs-use-default` option|
