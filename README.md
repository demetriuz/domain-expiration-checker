# domain-expiration-checker

**WARNING: Not for production usage, development in progress**

Usage:
```
Usage of ./domain-expiration-checker:
  -d value
        Domains
  -t int
        Expire Threshold Days (default 30)

./domain-expiration-checker -d ya.ru -d yandex.ru -t 340
ya.ru: 2019-09-01 00:00:00 +0000 UTC
```

# TODO:
- [ ] configs for `whois_backends`
- [x] use gorutines
- [ ] backends should configure app
- [ ] tests
- [ ] build infrastructure
- [ ] get domains from STDIN