docker create --name pg -p 5433:5432 -e POSTGRES_DB=Empl -e POSTGRES_USER=user -e POSTGRES_PASSWORD=1234 postgres:latest
::docker network connect mynet pg
docker start pg

