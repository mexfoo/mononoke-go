<p align="center">
        <img height="100px" src="https://git.mex.network/assets/img/logo.svg" />
</p>

<h2 align="center">mononoke-go<br>MMORPG authentication server in Go</h2>

### Intro
In search for a larger project to learn Go, I decided to rewrite my original authentication server "mononoke", previously written in C++.  

### Features
- Supports MySQL/MariaDB, PostgreSQL, SQL Server and sqlite3
- Entirely written in Go, gotta go speedy
- Easy to work with
- Easily extendable
- Runs with docker
- config via environment variables or `config.yml` file

### Configuration  
You can pass the configuration either via `config.yml` file or using environment variables (e.g. when using docker).  
`config.yml`
```yaml
database:
  dialect: sqlite3 # possible values sqlite3, mysql, sqlserver, postgres
  connection: data/mononoke.db # or DSN for the dialect
  defaultsalt: 2010 # salt for password hashing

defaultuser:
  name: test # the username of the default user
  password: test # the username of the default user

server:
  authclient:
    listenip: 127.0.0.1 # use 0.0.0.0 for external access
    listenport: 4500 # default port
    useencryption: true # default for Auth <-> Client
    encryptionkey: test  # use proper encryption key

  authgame:
    listenip: 127.0.0.1 # use 0.0.0.0 for external access, usually not necessary
    listenport: 4502 # default port
    useencryption: false # default for Auth <-> Game
    encryptionkey: test  # use proper encryption key 

  defaultdeskey: password # use proper DES key
  agerestriction: 18 # default

loggerlevel: Info # possible values Info, Debug, Error, Warning
loggerType: Text # possible values Text (default), JSON
```  
  
If you want to use environment variables instead, all variables have the `MONONOKE` prefix. Possible environment variables are named the same as the `config.yml` configuration, just pass an `_` between each level:
```bash
MONONOKE_DATABASE_DIALECT=sqlite3
MONONOKE_DATABASE_CONNECTION=data/mononoke.db
MONONOKE_SERVER_AUTHCLIENT_LISTENIP=127.0.0.1
MONONOKE_SERVER_AUTHCLIENT_ENCRYPTIONKEY=test
MONONOKE_SERVER_AUTHGAME_LISTENIP=127.0.0.1
MONONOKE_SERVER_DEFAULTDESKEY=password
MONONOKE_LOGGERLEVEL=Info
```