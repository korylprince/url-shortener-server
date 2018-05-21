# About

This is the backend for a URL shortening service. The frontend can be found [here](https://github.com/korylprince/url-shortener-client).

# Install

```bash
go get github.com/korylprince/url-shortener-server
./build.sh
```

# Configuration

The server is configured with environment variables:

```bash
SHORTENER_SESSIONEXPIRATION="60" # In minutes
SHORTENER_DATABASEPATH="/path/to/urls.db"
SHORTENER_URLIDLENGTH="6" # Length of random URL id. Recommended to leave at 6
SHORTENER_LDAPSERVER="ldap.example.com"
SHORTENER_LDAPPORT="389"
SHORTENER_LDAPBASEDN="OU=Container,DC=example,DC=net"
SHORTENER_LDAPGROUP="Group" # This will be the group's CN
SHORTENER_LDAPSECURITY="starttls" # none, tls, or starttls
SHORTENER_TLSCERT="/path/to/cert.pem"
SHORTENER_TLSKEY="/path/to/key.pem"
SHORTENER_LISTENADDR=":8080"
SHORTENER_PREFIX="/short" # Used to prefix all URLs
SHORTENER_DEBUG="false" # Show extra debugging information in server logs and client API
```

For more information see [config.go](https://github.com/korylprince/url-shortener-server/blob/master/config.go).

# Docker

You can use the pre-built Docker container, [korylprince/url-shortener-server:v1.0.0](https://hub.docker.com/r/korylprince/url-shortener-server/).

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
    korylprince/url-shortener-server:v1.0.0
```

#### Use LDAP StartTLS + HTTP TLS + LDAP Group

```bash
docker run -d --name="url-shortener" \
    -p 443:443 \
    -v /data:/data \
    -e SHORTENER_DATABASEPATH="/data/urls.db" \
    -e SHORTENER_LDAPSERVER="ldap.example.com" \
    -e SHORTENER_LDAPPORT="389" \
    -e SHORTENER_LDAPBASEDN="OU=Container,DC=example,DC=net" \
    -e SHORTENER_LDAPGROUP="Example Group" \
    -e SHORTENER_LDAPSECURITY="starttls" \
    -e SHORTENER_LISTENADDR=":443" \
    -e SHORTENER_TLSCERT="/data/chained_cert.pem" \
    -e SHORTENER_TLSKEY="/data/key.pem" \
    --restart="always" \
    korylprince/url-shortener-server:v1.0.0
```
