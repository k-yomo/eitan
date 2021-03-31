# eitan

## Setup
1. Install tools
```sh
$ make setup
```

2. Fill .env
``` 
$ cp .env.sample .env
$ vim .env
```

3. Add below hosts to `/etc/hosts`
```
127.0.0.1 account.local.eitan-flash.com
127.0.0.1 api.local.eitan-flash.com
127.0.0.1 local.eitan-flash.com
```

## Run servers
- Run local servers and DB on docker
```sh
$ make run
```

- Run everything on docker
```sh
$ make run-dc
```