# GORM DM Driver
GORM DM Driver for connect Dameng DB and manager Dameng DB.

# Support Dependency
- DM v8
- Golang 1.18+
- gorm 1.25+


# Quick Start
## How to install
```go
go get -d github.com/sineycoder/gorm-dm
```
## How to Use
```go
package main

import (
	dm "github.com/sineycoder/gorm-dm"
	"gorm.io/gorm"
)

func main() {
	// dm://user:password@127.0.0.1:1521?autoCommit=true
	url := dm.BuildDsn("127.0.0.1", 1521, "user", "password", nil)
	db, err := gorm.Open(dm.Open(url), &gorm.Config{})
	if err != nil {
		// panic error or log error info
	}

	// do somethings
}
```
