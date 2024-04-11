  docker run --rm --name redis -d --network mynet -p 6379:6379 -e ALLOW_EMPTY_PASSWORD=yes bitnami/redis:latest
