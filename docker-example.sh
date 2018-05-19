#!/bin/sh

docker run -d --name="url-shortener" \
    -p 80:80 \
    -v /data:/data \
    -e SHORTENER_DATABASEPATH="/data/urls.db" \
    -e SHORTENER_LDAPSERVER="ldap.example.com" \
    -e SHORTENER_LDAPPORT="389" \
    -e SHORTENER_LDAPBASEDN="OU=Container,DC=example,DC=net" \
    -e SHORTENER_LDAPGROUP="CN=Group,OU=Container,DC=example,DC=net" \
    -e SHORTENER_LDAPSECURITY="starttls" \
    -e SHORTENER_LISTENADDR=":80" \
    --restart="always" \
    korylprince/url-shortener-server:v1.0.0
