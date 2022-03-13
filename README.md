# Croupier

## Description
Personal single-player poker server



## Run it
```bash
PLAYER_NAME=`whoami` \
docker run \
  -p 5000:5000 \
  --network=host \
  --rm \
  -i \
  devbytom/croupier:latest


```

## License

[MIT](./LICENSE)
