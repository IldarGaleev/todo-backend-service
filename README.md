# TO DO List service

Backend gRPC service

## Environment variables

|Key       |Values              |Default|Description
|:--------:|--------------------|:-----:|---------------------------|
|`ENV_MODE`|`local`,`dev`,`prod`|`prod` |Production mode            |
|`PORT`    |`int`               |`9090` |gRPC server tcp port       |
|`DSN`     |`str`               |       |database connection string (**unimplement**) |