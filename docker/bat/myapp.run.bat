docker run --name app --rm --network mynet -p 81:8080 -e POSTGRES_HOST="pg" -e POSTGRES_PORT="5432" -e REDIS_ADDR="redis:6379" iapp


::не сработало, хотя по идее оба контейнера в дефолтной сети - доступа нет даже по IP
::docker run --name app --rm -p 81:8080  -e POSTGRES_HOST="172.21.0.2" -e POSTGRES_PORT="5432" iapp

:: не сработало, само приложение перестает быть доступным снаружи 
:: docker run --name app --rm -p 81:8080 -e POSTGRES_HOST="172.21.0.1" -e POSTGRES_PORT="5433" iapp


::docker create --name app --rm -p 81:8080  -e POSTGRES_HOST="pg" -e POSTGRES_PORT="5432"  iapp
::docker network connect mynet app
::docker network disconnect mynet app
::docker start app
