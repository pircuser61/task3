### Вопросы:

#### Большинство утилит не видит файлы

просто go vet в конрне проекта не видит файлы

D:\projects\go_db>go vet

` no Go files in D:\projects\go_db`

D:\projects\go_db>go vet .
` no Go files in D:\projects\go_db`

D:\projects\go_db>go vet go_db
` package go_db is not in std (C:\Program Files\Go\src\go_db)`

#### Ошибка при rollback

```
    defer func() {
    	err := tx.Rollback()
    	if err != nil {
    		// Что делать ?
    	}
    }()
```

#### Перекрытие ошибки

Может быть ошибка не из-за транзакции
Может быть ошибка commit
Может быть ошибка rollback

```
func (i PostgressStore) EmployeeCreate(ctx context.Context, empl models.Employee) (newiD uint32, err error) {


	defer func() {
		closeErr := tx.Close() // Close calls Rollback if the tx has not already been committed or rolled back.
		if err == nil {
			err = closeErr
		}
	}()
```

#### staticchek срабатывает на вспомогательные поля для ORM

internal\models\models.go:4:2: field tableName is unused (U1000)

### go fmt

go fmt ./...

### go vet

go vet ./cmd/server/main.go

` cmd\server\main.go:20:12: fmt.Printf call has arguments but no formatting directives`

go vet ./...

```
# go_db/cmd/mongo_test
cmd\mongo_test\mongodb.go:154:4: go.mongodb.org/mongo-driver/bson/primitive.E struct literal uses unkeyed fields
cmd\mongo_test\mongodb.go:156:6: go.mongodb.org/mongo-driver/bson/primitive.E struct literal uses unkeyed fields
cmd\mongo_test\mongodb.go:157:6: go.mongodb.org/mongo-driver/bson/primitive.E struct literal uses unkeyed fields
cmd\mongo_test\mongodb.go:157:23: go.mongodb.org/mongo-driver/bson/primitive.E struct literal uses unkeyed fields
# go_db/cmd/server
cmd\server\main.go:20:12: fmt.Printf call has arguments but no formatting directives
```

### test

#### bench

go test -bench ./internal/storage/postgress/sql
`no Go files in D:\projects\go_db`

go test -bench=./internal/storage/postgress/sql
`ok      go_db/internal/storage/postgress/sql    0.232s`

#### cover

go test ./... --cover

```
?       go_db/cmd/db_client     [no test files]
?       go_db/cmd/migrations/goose      [no test files]
?       go_db/cmd/migrations/migrate    [no test files]
?       go_db/cmd/mongo_test    [no test files]
# go_db/cmd/server
cmd\server\main.go:20:12: fmt.Printf call has arguments but no formatting directives
?       go_db/config    [no test files]
?       go_db/internal/models   [no test files]
?       go_db/internal/storage  [no test files]
?       go_db/internal/storage/go-redis [no test files]
?       go_db/internal/storage/mongo    [no test files]
?       go_db/internal/storage/postgress        [no test files]
?       go_db/internal/storage/postgress/go-pg  [no test files]
?       go_db/internal/storage/postgress/gorm   [no test files]
?       go_db/internal/storage/postgress/pgx    [no test files]
?       go_db/internal/storage/postgress/pgxpool        [no test files]
?       go_db/internal/storage/postgress/sqlx   [no test files]
?       go_db/internal/storage/redigo   [no test files]
ok      go_db/internal/storage/postgress/sql    (cached)        coverage: 44.9% of statements
FAIL
```

go test ./internal/storage/postgress/sql --cover
`ok go_db/internal/storage/postgress/sql (cached) coverage: 44.9% of statements`

#### coverprofile

go test ./internal/storage/postgress/sql --coverprofile=coverage.html
`go_db/internal/storage/postgress/sql    0.255s  coverage: 44.9% of statements`
go tool cover -html coverage.html

### errcheck

go install github.com/kisielk/errcheck@latest

D:\projects\go_db> errcheck ./...

```
cmd\db_client\main.go:118:27:   dbInstanse.EmployeeDelete(ctx, ux.Id)
cmd\mongo_test\mongodb.go:68:11:        coll.Drop(ctx)
cmd\mongo_test\mongodb.go:128:22:       singleResult.Decode(&p4)
cmd\server\main.go:91:21:       http.ListenAndServe(config.AppAddr, nil)
internal\storage\go-redis\redis.go:40:12:       i.cl.Close()
internal\storage\postgress\go-pg\go-pg.go:75:12:        i.db.Close()
internal\storage\postgress\go-pg\go-pg.go:87:16:        defer tx.Close() // Close calls Rollback if the tx has not already been committed or rolled back.
internal\storage\postgress\go-pg\go-pg.go:95:12:        tx.Commit()
internal\storage\postgress\pgx\pgx.go:44:14:    i.conn.Close(ctx)
internal\storage\postgress\sql\sql.go:61:12:    i.db.Close()
internal\storage\postgress\sql\sql.go:82:19:    defer tx.Rollback()
internal\storage\postgress\sql\sql.go:89:11:    tx.Commit()
internal\storage\postgress\sql\sql.go:144:18:   defer rows.Close()
internal\storage\postgress\sqlx\sqlx.go:50:12:  i.db.Close()
internal\storage\postgress\sqlx\sqlx.go:119:13: i.db.Select(&result, queries.QueryList)
internal\storage\redigo\redis.go:47:14: i.conn.Close()
```

### staticcheck

| Tool                  | Description                                                            |
| --------------------- | ---------------------------------------------------------------------- |
| keyify                | Transforms an unkeyed struct literal into a keyed one.                 |
| staticcheck           | Go static analysis, detecting bugs, performance issues, and much more. |
| structlayout          | Displays the layout (field sizes and padding) of structs.              |
| structlayout-optimize | Reorders struct fields to minimize the amount of padding.              |
| structlayout-pretty   | Formats the output of structlayout with ASCII art.                     |

~~go install honnef.co/go/tools/cmd/staticcheck@2022.1~~
go install honnef.co/go/tools/cmd/staticcheck@latest

staticcheck ./...

```
cmd\mongo_test\mongodb.go:189:6: func transaction is unused (U1000)
cmd\server\main.go:20:13: Printf call needs 0 args but has 1 args (SA5009)
internal\models\models.go:4:2: field tableName is unused (U1000)
internal\storage\postgress\sql\sql.go:87:3: empty branch (SA9003)
```

### golangci-lint

```
go install/go get installation isn't recommended because of the following points:

    Some users use -u flag for go get, which upgrades our dependencies. Resulting configuration wasn't tested and isn't guaranteed to work.
    go.mod replacement directive doesn't apply. It means a user will be using patched version of golangci-lint if we use such replacements.
    We've encountered a lot of issues with Go modules hashes.
    It allows installation from master branch which can't be considered stable.
    It's slower than binary installation.
```

~~choco install golangci-lint~~

~~On Windows, you can run the above commands with Git Bash, which comes with Git for Windows.~~
~~curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.57.2~~
Скачал непосредственно с гитхаба бинарник и сохранил в {$PATH}
golangci-lint --version
`golangci-lint has version 1.57.2 built with go1.22.1 from 77a8601a on 2024-03-28T19:01:11Z`

golangci-lint run

```
internal\models\models.go:4:2: field `tableName` is unused (unused)
        tableName struct{} `pg:"employee"` // for go-pg
        ^
cmd\mongo_test\mongodb.go:189:6: func `transaction` is unused (unused)
func transaction(ctx context.Context, cl *mongo.Client) error {
     ^
cmd\mongo_test\mongodb.go:161:4: composites: go.mongodb.org/mongo-driver/bson/primitive.E struct literal uses unkeyed fields (govet)
                        {"$group",
                        ^
cmd\mongo_test\mongodb.go:163:6: composites: go.mongodb.org/mongo-driver/bson/primitive.E struct literal uses unkeyed fields (govet)
                                        {"_id", "$name"},
                                        ^
cmd\mongo_test\mongodb.go:164:6: composites: go.mongodb.org/mongo-driver/bson/primitive.E struct literal uses unkeyed fields (govet)
                                        {"count", bson.D{{"$sum", 1}}},
                                        ^
cmd\server\main.go:20:12: printf: fmt.Printf call has arguments but no formatting directives (govet)
        fmt.Printf("bad format for vet", 12)
                  ^
internal\storage\postgress\sql\sql.go:87:3: SA9003: empty branch (staticcheck)
                if err != nil {
```

git tag -a "v0.0.1" -m "for golangci-lint"
golangci-lint run --new-from-rev "v0.0.1"

```internal\storage\postgress\sql\sql.go:87:3: SA9003: empty branch (staticcheck)
                if err != nil {
                ^
cmd\server\main.go:20:12: printf: fmt.Printf call has arguments but no formatting directives (govet)
        fmt.Printf("bad format for vet", 12)
```

golangci-lint linters
golangci-lint run --new-from-rev --enable-all

### git hooks

https://medium.com/@radlinskii/writing-the-pre-commit-git-hook-for-go-files-810f8d5f1c6f

\go_db\.git\hooks\pre-commit:

```
#!/bin/bash
echo "PRE COMMIT"
STAGED_GO_FILES=$(git diff --cached --name-only | grep ".go$")

if [[ "$STAGED_GO_FILES" = "" ]]; then
    exit 0
fi

PASS=true

for FILE in $STAGED_GO_FILES
do
  go vet $FILE
  if [[ $? != 0 ]]; then
    PASS=false
  fi
done

if ! $PASS; then
  echo "COMMIT FAILED"
  exit 1
else
  echo "COMMIT SUCCEEDED"
fi

exit 0
```
