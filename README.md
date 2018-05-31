based on [this gist](https://gist.github.com/bastman/5b57ddb3c11942094f8d0a97d461b430) i reated a small docker image to clean up docker hosts

## usage

run it with with docker.sock as volume

`docker run -v /var/run/docker.sock:/var/run/docker.sock cdreier/docker-cleanup:latest`

or with docker-compose

```
version: '3'

services:
  app: 
    image: cdreier/docker-cleanup:latest
    # restart: always
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
```

##