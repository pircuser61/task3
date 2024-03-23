#migrate create -ext sql init --создать файлы для миграции
migrate -database postgres://user:1234@localhost:5433/empl?sslmode=disable -path ./ up
