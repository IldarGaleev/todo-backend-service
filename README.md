# TO DO List service

Backend gRPC service

## Environment variables

|Key               |Values              |Default|Description
|:----------------:|--------------------|:-----:|---------------------------
|`ENV_MODE`        |`local`,`dev`,`prod`|`prod` |Production mode
|`PORT`            |`int`               |`9090` |gRPC server tcp port
|`DSN`             |`str`               |       |database connection string
|`SECRET_KEY`      |`bytes`             |       |private key for JWT
|`SECRETS_MAX_AGE` |`duration`          |`24h`  |JWT token max age

## Cmd

<table>
<tr>
  <th>Cmd</th>
  <th>Flags</th>
  <th>Description</th>
</tr>
<tr>
<td><code>todo\main</code></td>
<td>
-
</td>
<td>run backend server</td>
</tr>
<tr>
<td><code>utils\createuser</code></td>
<td><ul><li><code>-username</code></li>
<li><code>-password</code></li></ul></td>
<td>create new user</td>
</tr>
</table>
