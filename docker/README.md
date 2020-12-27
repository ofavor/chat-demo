# docker

Docker related files

## Run Dev

```
docker-compose -f docker-compose-dev.yml up -d
```

## Logs

```
docker logs -f docker_gateway_1
docker logs -f docker_chat-srv_1
```

## CLI MySQL

```
docker exec -it docker_mysql_1 mysql -uroot -proot
```

Or

```
docker run --network docker_default -it --rm mysql:5.7 mysql -uroot -p -h mysql
```

## CLI Redis

```
docker exec -it docker_redis_1 redis-cli
```

Or

```
docker run --network docker_default -it --rm redis redis-cli -h redis
```