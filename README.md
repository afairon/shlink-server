# Shlink URL shortener server

## Main Features

- URL Shortening
- Visitor Counting

## Example of configuration file

```yaml
database:
  host: 127.0.0.1
  port: "27017"
  db: shlink
log:
  logName: logs/shlink.log
  maxSize: 10
  maxBackups: 2
  maxAge: 7
server:
  host: 127.0.0.1
  port: "8080"
  base: http://127.0.0.1:8080
```

## Example of cli

```

  _____ _     _ _       _       _____                          
 / ____| |   | (_)     | |     / ____|                         
| (___ | |__ | |_ _ __ | | __ | (___   ___ _ ____   _____ _ __ 
 \___ \| '_ \| | | '_ \| |/ /  \___ \ / _ \ '__\ \ / / _ \ '__|
 ____) | | | | | | | | |   <   ____) |  __/ |   \ V /  __/ |   
|_____/|_| |_|_|_|_| |_|_|\_\ |_____/ \___|_|    \_/ \___|_|   
                                                                
(C) Copyright 2018 Shlink. All Rights Reserved.

Usage:
  shlink-server [command]

Available Commands:
  config      Shows config
  help        Help about any command
  server      Start shlink http server
  version     Shows binary version

Flags:
  -h, --help        help for shlink-server
      --no-banner   Don't display banner

Use "shlink-server [command] --help" for more information about a command.

```