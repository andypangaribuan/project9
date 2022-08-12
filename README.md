<!-- 
Find errors not caught by the compilers.
This command vets the package in the current directory.
  $ go vet
Download all dependencies
  $ go mod download
Remove unused dependencies
  $ go mod tidy

IntelliJ Shortcuts
✦ opt + (↑ OR ↓)    ︎〉extend selection
✦ ctrl + opt + I    ︎〉auto-indent lines
✦ ctrl + shift + E  ︎〉next highlighted error (custom)
✦ ctrl + tab        ︎〉switcher
✦ cmd + .  ︎         〉collapse or expand block
✦ ctrl + shift + -  ︎〉collapse all
✦ ctrl + shift + +  ︎〉expand all
-->

# Project 9

## _One-stop Solution For All Your GO Project_

Project 9 is a small and light but powerful go framework.
Also, it has proven to be saving a lot of development hours.
If you need performance and good productivity, you will love Project 9.

## Installation

```sh
go get github.com/andypangaribuan/project9
```
Then do initialization before using the p9.
```sh
project9.Initialize()
```

## Interfaces

| Code      | Description      |
|-----------|------------------|
| p9.Check  | Value checker    |
| p9.Conf   | P9 Configuration |
| p9.Conv   | Converter        |
| p9.Crypto | Cryptography     |
| p9.Db     | Database         |
| p9.Err    | Error wrap       |
| p9.Http   | API Caller       |
| p9.Json   | Json             |
| p9.Log    | Logging          |
| p9.Server | TCP Server       |
| p9.Util   | Utilities        |

> f9 : access direct function
