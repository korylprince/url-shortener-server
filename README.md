# About

This is the backend for a URL shortening service. The frontend can be found [here](https://github.com/korylprince/url-shortener-client).

# Install

```
go get github.com/korylprince/url-shortener-server
```

# Configuration

The server is configured with environment variables:

```
SHORTENER_SESSIONEXPIRATION="60"
SHORTENER_DATABASEPATH="/path/to/urls.db"
SHORTENER_URLIDLENGTH="6"
SHORTENER_LDAPSERVER="ldap.example.com"
SHORTENER_LDAPPORT="389"
SHORTENER_LDAPBASEDN="OU=Container,DC=example,DC=net"
SHORTENER_LDAPGROUP="CN=Group,OU=Container,DC=example,DC=net"
SHORTENER_LDAPSECURITY="starttls"
SHORTENER_LISTENADDR=":8080"
SHORTENER_PREFIX="/short"
SHORTENER_DEBUG="true"
```

For more information see [config.go](https://github.com/korylprince/url-shortener-server/blob/master/config.go).
