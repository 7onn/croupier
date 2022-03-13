# Croupier

## Description
My personal single-player poker server

## Run it
Start the server
```bash
PLAYER_NAME=`whoami` \
docker run \
  -p 5000:5000 \
  --network=host \
  --rm \
  -i \
  devbytom/croupier:latest
```
Then access it through your browser [http://localhost:5000/](http://localhost:5000/)

## License
[MIT](./LICENSE)
