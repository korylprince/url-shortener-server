# About

This is the backend for a URL shortening service. The frontend can be found [here](https://github.com/korylprince/url-shortener-client).

# Building

## Client

The client is vendored in the repo so is only needed to be built if changes are made. See `build-client.sh` for a client build script. [nvm](https://github.com/nvm-sh/nvm) is required to build the client.

## Server


```bash
$ cd /path/to/build/directory
$ GOBIN="$(pwd)" go install "github.com/korylprince/url-shortener-server@<tagged version>"
```

# Configuration

The server is configured with environment variables:

```bash
SHORTENER_SESSIONEXPIRATION="60" # In minutes
SHORTENER_DATABASEPATH="/path/to/urls.db"
SHORTENER_URLIDLENGTH="6" # Length of random URL id. Recommended to leave at 6
SHORTENER_APPTITLE="My Shortener" # Set to change name of app in client
SHORTENER_LDAPSERVER="ldap.example.com"
SHORTENER_LDAPPORT="389"
SHORTENER_LDAPBASEDN="OU=Container,DC=example,DC=net"
SHORTENER_LDAPGROUP="Group" # This will be the group's CN
SHORTENER_LDAPADMINGROUP="Admin Group" # Group allowed to access Admin Interface
SHORTENER_LDAPSECURITY="starttls" # none, tls, or starttls
SHORTENER_TLSCERT="/path/to/cert.pem"
SHORTENER_TLSKEY="/path/to/key.pem"
SHORTENER_LISTENADDR=":8080"
SHORTENER_PREFIX="/short" # Used to prefix all URLs
```

For more information see [config.go](https://github.com/korylprince/url-shortener-server/blob/master/config.go).

# Docker

You can use the pre-built Docker container, [ghcr.io/korylprince/url-shortener-server](https://github.com/korylprince/url-shortener-server/pkgs/container/url-shortener-server).

### Examples

#### No Security

```bash
docker run -d --name="url-shortener" \
    -p 80:80 \
    -v /data:/data \
    -e SHORTENER_DATABASEPATH="/data/urls.db" \
    -e SHORTENER_LDAPSERVER="ldap.example.com" \
    -e SHORTENER_LDAPPORT="389" \
    -e SHORTENER_LDAPBASEDN="OU=Container,DC=example,DC=net" \
    -e SHORTENER_LDAPSECURITY="none" \
    -e SHORTENER_LISTENADDR=":80" \
    --restart="always" \
    ghcr.io/korylprince/url-shortener-server:latest
```

#### Use Custom App Title + LDAP StartTLS + HTTP TLS + LDAP Groups

```bash
docker run -d --name="url-shortener" \
    -p 443:443 \
    -v /data:/data \
    -e SHORTENER_DATABASEPATH="/data/urls.db" \
    -e SHORTENER_APPTITLE="My Shortener" \
    -e SHORTENER_LDAPSERVER="ldap.example.com" \
    -e SHORTENER_LDAPPORT="389" \
    -e SHORTENER_LDAPBASEDN="OU=Container,DC=example,DC=net" \
    -e SHORTENER_LDAPGROUP="Example Group" \
    -e SHORTENER_LDAPADMINGROUP="Example Admin Group" \
    -e SHORTENER_LDAPSECURITY="starttls" \
    -e SHORTENER_LISTENADDR=":443" \
    -e SHORTENER_TLSCERT="/data/chained_cert.pem" \
    -e SHORTENER_TLSKEY="/data/key.pem" \
    --restart="always" \
    ghcr.io/korylprince/url-shortener-server:latest
```
