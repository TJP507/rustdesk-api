# RustDesk API

[English Doc](README_EN.md)

This project implements the RustDesk API in Go and includes a Web Admin panel and a Web Client.


<div align=center>
<img src="https://img.shields.io/badge/golang-1.22-blue"/>
<img src="https://img.shields.io/badge/gin-v1.9.0-lightBlue"/>
<img src="https://img.shields.io/badge/gorm-v1.25.7-green"/>
<img src="https://img.shields.io/badge/swag-v1.16.3-yellow"/>
<img src="https://goreportcard.com/badge/github.com/lejianwen/rustdesk-api/v2"/>
<img src="https://github.com/lejianwen/rustdesk-api/actions/workflows/build.yml/badge.svg"/>
</div>

## Works best paired with [lejianwen/rustdesk-server].
> [lejianwen/rustdesk-server] is a fork of the official RustDesk Server repository.
> 1. Fixes the connection timeout issue when using the API
> 2. Allows enforcing login before a connection can be initiated
> 3. Supports client WebSocket



# Features

- PC Client API
    - Personal edition API
    - Login
    - Address book
    - Groups
    - OAuth login
      - Supports `github`, `google`, and `OIDC` login
      - Supports authorization login via the `web admin panel`
      - Supports `LDAP` (AD and OpenLDAP have been tested), when LDAP is configured on the API Server
    - i18n
- Web Admin
    - User management
    - Device management
    - Address book management
    - Tag management
    - Group management
    - OAuth management
    - Configure LDAP via config file or environment variables
    - Login logs
    - Connection logs
    - File transfer logs
    - Quick access to the web client
    - i18n
    - Share with guests via the web client
    - Server control (a selection of official simple commands — [WIKI](https://github.com/lejianwen/rustdesk-api/wiki/Rustdesk-Command))
- Web Client
    - Automatically fetches the API server
    - Automatically fetches the ID server and KEY
    - Automatically fetches the address book
    - Guests can remotely access devices directly via a temporary share link
- CLI
    - Reset the administrator password

## Functionality


### API Service
The core PC client API endpoints are fully implemented. The Personal edition API is also supported and can be enabled or disabled via the config file key `rustdesk.personal` or the environment variable `RUSTDESK_API_RUSTDESK_PERSONAL`.

<table>
    <tr>
      <td width="50%" align="center" colspan="2"><b>Login</b></td>
    </tr>
    <tr>
        <td width="50%" align="center" colspan="2"><img src="docs/pc_login.png"></td>
    </tr>
     <tr>
      <td width="50%" align="center"><b>Address Book</b></td>
      <td width="50%" align="center"><b>Groups</b></td>
    </tr>
    <tr>
        <td width="50%" align="center"><img src="docs/pc_ab.png"></td>
        <td width="50%" align="center"><img src="docs/pc_gr.png"></td>
    </tr>
</table>

### Web Admin:

* Uses a decoupled frontend/backend architecture to provide a user-friendly management interface, primarily for administration and display. The frontend source code is at [rustdesk-api-web](https://github.com/lejianwen/rustdesk-api-web).

* The admin panel is accessible at `http://<your server>[:port]/_admin/`
* On first installation, the administrator username is `admin`. The password is printed to the console and can be changed via the [CLI](#CLI).

  ![img.png](./docs/init_admin_pwd.png)

1. Administrator interface
   ![web_admin](docs/web_admin.png)
2. Regular user interface
   ![web_user](docs/web_admin_user.png)

3. Each user can have multiple address books and can share them with other users.
4. Groups are fully customizable for easier management. Two group types are currently supported: `Shared Group` and `Regular Group`.
5. The web client can be opened directly for convenient use, or it can be shared with guests who can then remotely access devices directly through the web client.
6. OAuth — supports `Github`, `Google`, and `OIDC`. You must create an `OAuth App` and configure it in the admin panel.
    - For `Google` and `Github`, `Issuer` and `Scopes` do not need to be filled in.
    - For `OIDC`, `Issuer` is required. `Scopes` is optional and defaults to `openid,profile,email`. Ensure that `sub`, `email`, and `preferred_username` can be retrieved.
    - Create a `GitHub OAuth App` at `Settings` -> `Developer settings` -> `OAuth Apps` -> `New OAuth App`, at [https://github.com/settings/developers](https://github.com/settings/developers).
    - Set the `Authorization callback URL` to `http://<your server[:port]>/api/oidc/callback`, for example `http://127.0.0.1:21114/api/oidc/callback`.
7. Login logs
8. Connection logs
9. File transfer logs
10. Server control

  - `Simple mode` — a UI-driven interface for common simple commands that can be executed directly from the admin panel.
    ![rustdesk_command_simple](./docs/rustdesk_command_simple.png)

  - `Advanced mode` — execute commands directly in the admin panel.
      * Supports official commands
      * Allows adding custom commands
      * Allows executing custom commands

 
11. **LDAP support** — when LDAP is configured on the API Server (AD and OpenLDAP have been tested), users can log in using their LDAP credentials. See https://github.com/lejianwen/rustdesk-api/issues/114. If LDAP authentication fails, the system falls back to local user authentication.

### Web Client:

1. If you are already logged into the admin panel, the web client will log in automatically.
2. If you are not logged into the admin panel, click the login button in the top-right corner — the API server is already pre-configured.
3. After logging in, the ID server and KEY are automatically synchronized.
4. After logging in, the address book is automatically saved into the web client for convenient access.


### Auto-generated Documentation: API documentation is generated with Swag to help developers understand and use the API.

1. Admin panel docs: `<your server[:port]>/admin/swagger/index.html`
2. PC client docs: `<your server[:port]>/swagger/index.html`
   ![api_swag](docs/api_swag.png)

### CLI

```bash
# Show help
./apimain -h
```

#### Reset the administrator password
```bash
./apimain reset-admin-pwd <pwd>
```

## Installation and Setup

### Configuration

* [Configuration file](./conf/config.yaml)
* Refer to `conf/config.yaml` and adjust the relevant settings as needed.
* If `gorm.type` is set to `sqlite`, MySQL configuration is not required.
* If no language is set, the default is `zh-CN`.

### Environment Variables
Environment variables correspond one-to-one with the settings in `conf/config.yaml`. All variable names are prefixed with `RUSTDESK_API`.
The table below is not exhaustive — refer to `conf/config.yaml` for a full list of options.

| Variable Name                                          | Description                                                                                                                                          | Example                      |
|--------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------|
| TZ                                                     | Timezone                                                                                                                                             | Asia/Shanghai                |
| RUSTDESK_API_LANG                                      | Language                                                                                                                                             | `en`,`zh-CN`                 |
| RUSTDESK_API_APP_WEB_CLIENT                            | Enable the web client; 1: enabled, 0: disabled; enabled by default                                                                                  | 1                            |
| RUSTDESK_API_APP_REGISTER                              | Enable user registration; `true`, `false`; default `false`                                                                                           | `false`                      |
| RUSTDESK_API_APP_SHOW_SWAGGER                          | Show Swagger documentation; `1` to show, `0` to hide; default `0` (hidden)                                                                           | `1`                          |
| RUSTDESK_API_APP_TOKEN_EXPIRE                          | Token validity duration                                                                                                                              | `168h`                       |
| RUSTDESK_API_APP_DISABLE_PWD_LOGIN                     | Disable password login; `true`, `false`; default `false`                                                                                             | `false`                      |
| RUSTDESK_API_APP_REGISTER_STATUS                       | Default status for newly registered users; 1 = enabled, 2 = disabled; default 1                                                                     | `1`                          |
| RUSTDESK_API_APP_CAPTCHA_THRESHOLD                     | CAPTCHA trigger threshold; -1 = disabled, 0 = always enabled, >0 = enabled after that many failed logins; default `3`                               | `3`                          |
| RUSTDESK_API_APP_BAN_THRESHOLD                         | IP ban trigger threshold; 0 = disabled, >0 = ban IP after that many failed logins; default `0`                                                       | `0`                          |
| -----ADMIN CONFIG-----                                 | ----------                                                                                                                                           | ----------                   |
| RUSTDESK_API_ADMIN_TITLE                               | Admin panel title                                                                                                                                    | `RustDesk Api Admin`         |
| RUSTDESK_API_ADMIN_HELLO                               | Admin panel welcome message; supports `html`                                                                                                         |                              |
| RUSTDESK_API_ADMIN_HELLO_FILE                          | Admin panel welcome message file; useful when the content is long. Overrides `RUSTDESK_API_ADMIN_HELLO`.<br>                                         | `./conf/admin/hello.html`    |
| -----GIN CONFIG-----                                   | ----------                                                                                                                                           | ----------                   |
| RUSTDESK_API_GIN_TRUST_PROXY                           | Comma-separated list of trusted proxy IPs; trusts all by default                                                                                     | 192.168.1.2,192.168.1.3      |
| -----GORM CONFIG-----                                  | ----------                                                                                                                                           | ---------------------------  |
| RUSTDESK_API_GORM_TYPE                                 | Database type: `sqlite` or `mysql`; default `sqlite`                                                                                                 | sqlite                       |
| RUSTDESK_API_GORM_MAX_IDLE_CONNS                       | Maximum number of idle database connections                                                                                                          | 10                           |
| RUSTDESK_API_GORM_MAX_OPEN_CONNS                       | Maximum number of open database connections                                                                                                          | 100                          |
| RUSTDESK_API_RUSTDESK_PERSONAL                         | Enable the Personal edition API; 1: enabled, 0: disabled; enabled by default                                                                         | 1                            |
| -----MYSQL CONFIG-----                                 | ----------                                                                                                                                           | ----------                   |
| RUSTDESK_API_MYSQL_USERNAME                            | MySQL username                                                                                                                                       | root                         |
| RUSTDESK_API_MYSQL_PASSWORD                            | MySQL password                                                                                                                                       | 111111                       |
| RUSTDESK_API_MYSQL_ADDR                                | MySQL address                                                                                                                                        | 192.168.1.66:3306            |
| RUSTDESK_API_MYSQL_DBNAME                              | MySQL database name                                                                                                                                  | rustdesk                     |
| RUSTDESK_API_MYSQL_TLS                                 | Enable TLS; accepted values: `true`, `false`, `skip-verify`, `custom`                                                                               | `false`                      |
| -----RUSTDESK CONFIG-----                              | ----------                                                                                                                                           | ----------                   |
| RUSTDESK_API_RUSTDESK_ID_SERVER                        | RustDesk ID server address                                                                                                                           | 192.168.1.66:21116           |
| RUSTDESK_API_RUSTDESK_RELAY_SERVER                     | RustDesk relay server address                                                                                                                        | 192.168.1.66:21117           |
| RUSTDESK_API_RUSTDESK_API_SERVER                       | RustDesk API server address                                                                                                                          | http://192.168.1.66:21114    |
| RUSTDESK_API_RUSTDESK_KEY                              | RustDesk key                                                                                                                                         | 123456789                    |
| RUSTDESK_API_RUSTDESK_KEY_FILE                         | File containing the RustDesk key                                                                                                                     | `./conf/data/id_ed25519.pub` |
| RUSTDESK_API_RUSTDESK_WEBCLIENT<br/>_MAGIC_QUERYONLINE | Enable the new online-status query method in web client v2; `1`: enabled, `0`: disabled; disabled by default                                         | `0`                          |
| RUSTDESK_API_RUSTDESK_WS_HOST                          | Custom WebSocket host                                                                                                                                | `wss://192.168.1.123:1234`   |
| ----PROXY CONFIG-----                                  | ----------                                                                                                                                           | ----------                   |
| RUSTDESK_API_PROXY_ENABLE                              | Enable proxy: `false`, `true`                                                                                                                        | `false`                      |
| RUSTDESK_API_PROXY_HOST                                | Proxy address                                                                                                                                        | `http://127.0.0.1:1080`      |
| ----JWT CONFIG----                                     | --------                                                                                                                                             | --------                     |
| RUSTDESK_API_JWT_KEY                                   | Custom JWT key; leave empty to disable JWT.<br/>If you are not using `MUST_LOGIN` from `lejianwen/rustdesk-server`, it is recommended to leave this empty. |                              |
| RUSTDESK_API_JWT_EXPIRE_DURATION                       | JWT validity duration                                                                                                                                | `168h`                       |


### Running

#### Running with Docker

1. Run directly with Docker. Configuration can be adjusted by mounting a config file at `/app/conf/config.yaml`, or by overriding settings with environment variables.

    ```bash
    docker run -d --name rustdesk-api -p 21114:21114 \
    -v /data/rustdesk/api:/app/data \
    -e TZ=Asia/Shanghai \
    -e RUSTDESK_API_LANG=zh-CN \
    -e RUSTDESK_API_RUSTDESK_ID_SERVER=192.168.1.66:21116 \
    -e RUSTDESK_API_RUSTDESK_RELAY_SERVER=192.168.1.66:21117 \
    -e RUSTDESK_API_RUSTDESK_API_SERVER=http://192.168.1.66:21114 \
    -e RUSTDESK_API_RUSTDESK_KEY=<key> \
    lejianwen/rustdesk-api
    ```

2. Using `docker compose` — refer to the [WIKI](https://github.com/lejianwen/rustdesk-api/wiki).

#### Download a Release and Run Directly

[Download page](https://github.com/lejianwen/rustdesk-api/releases)

#### Install from Source

1. Clone the repository
   ```bash
   git clone https://github.com/lejianwen/rustdesk-api.git
   cd rustdesk-api
   ```

2. Install dependencies

    ```bash
    go mod tidy
    # Install swag — skip this if you do not need to generate documentation
    go install github.com/swaggo/swag/cmd/swag@latest
    ```

3. Build the admin frontend. The frontend source code is at [rustdesk-api-web](https://github.com/lejianwen/rustdesk-api-web).
   ```bash
   cd resources
   mkdir -p admin
   git clone https://github.com/lejianwen/rustdesk-api-web
   cd rustdesk-api-web
   npm install
   npm run build
   cp -ar dist/* ../admin/
   ```
4. Run
    ```bash
    # Run directly
    go run cmd/apimain.go
    # Or generate the API and run using generate_api.go
    go generate generate_api.go
    ```
   > Note: When using `go run` or a compiled binary, the `conf` and `resources` directories
   > must exist in the current working directory. If you run from a different directory,
   > you can specify absolute paths using the `-c` flag and the
   > `RUSTDESK_API_GIN_RESOURCES_PATH` environment variable, for example:
   > ```bash
   > RUSTDESK_API_GIN_RESOURCES_PATH=/opt/rustdesk-api/resources ./apimain -c /opt/rustdesk-api/conf/config.yaml
   > ```
5. Build — if you want to build yourself, change to the project root directory, then run `build.bat` on Windows or `build.sh` on Linux. The compiled executable will be placed in the `release` directory. Run it directly.

6. Open a browser and navigate to `http://<your server[:port]>/_admin/`. The default username and password are both `admin` — change the password promptly.


#### Running with the `lejianwen/server-s6` Image

- The connection timeout issue is resolved
- Login can be enforced before a connection is initiated
- GitHub: https://github.com/lejianwen/rustdesk-server

```yaml
 networks:
   rustdesk-net:
     external: false
 services:
   rustdesk:
     ports:
       - 21114:21114
       - 21115:21115
       - 21116:21116
       - 21116:21116/udp
       - 21117:21117
       - 21118:21118
       - 21119:21119
     image: lejianwen/rustdesk-server-s6:latest
     environment:
       - RELAY=<relay_server[:port]>
       - ENCRYPTED_ONLY=1
       - MUST_LOGIN=N
       - TZ=Asia/Shanghai
       - RUSTDESK_API_RUSTDESK_ID_SERVER=<id_server[:21116]>
       - RUSTDESK_API_RUSTDESK_RELAY_SERVER=<relay_server[:21117]>
       - RUSTDESK_API_RUSTDESK_API_SERVER=http://<api_server[:21114]>
       - RUSTDESK_API_KEY_FILE=/data/id_ed25519.pub
       - RUSTDESK_API_JWT_KEY=xxxxxx # jwt key
     volumes:
       - /data/rustdesk/server:/data
       - /data/rustdesk/api:/app/data # mount the database
     networks:
       - rustdesk-net
     restart: unless-stopped
       
```


## Other Resources

- [WIKI](https://github.com/lejianwen/rustdesk-api/wiki)
- [Connection timeout issue](https://github.com/lejianwen/rustdesk-api/issues/92)
- [Change client ID](https://github.com/abdullah-erturk/RustDesk-ID-Changer)
- [Web client source](https://hub.docker.com/r/keyurbhole/flutter_web_desk)


## Acknowledgements

Thanks to everyone who has contributed!

<a href="https://github.com/lejianwen/rustdesk-api/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=lejianwen/rustdesk-api" />
</a>

## Thank you for your support! If this project has been useful to you, please give it a ⭐️ — it means a lot!

[lejianwen/rustdesk-server]: https://github.com/lejianwen/rustdesk-server