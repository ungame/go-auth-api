# Golang Authentication API with Fiber and MongoDB

## Run MongoDB

```bash
    docker run -it --rm --name mongodb_container -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=admin -v mongodata:/data/db -d -p 27017:27017 mongo

    docker exec -it mongodb_container /bin/bash

    mongo -u admin -p admin --authenticationDatabase admin

    use authapi

    db.createUser({user: 'user', pwd: 'password', roles:[{'role': 'readWrite', 'db': 'authapi'}]});

    mongo -u user -p password --authenticationDatabase authapi

    use authapi

    show collections
```

## Create Your Env File on Root Directory

```bash
# example
JWT_SECRET_KEY=secret
DATABASE_USER=user
DATABASE_PASS=password
DATABASE_HOST=127.0.0.1
DATABASE_PORT=27017
DATABASE_NAME=authapi
```

## Run API

```bash
    go run main.go -port=8080
```