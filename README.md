# gouro

Solves the issue of duplicate or uninteresting urls for pentest or bugbounty

It doesn't make any http requests to the URLs and removes:
- human written content e.g. blog posts
- urls with same path but parameter value difference
- incremental urls e.g. `/cat/1/` and `/cat/2/`
- image, js, css and other static files

#### Install

```
go get -u github.com/felipemelchior/gouro
```

#### Usage

![gouro-demo](https://i.ibb.co/T0vYyR1/Captura-de-Tela-2022-03-24-s-11-38-58.png)

Good to use in oneliners

```
echo "domain.com" | assetfinder -subs-only | waybackurls | gouro | gf xss
```
## Credits
```
Go implementation of https://github.com/s0md3v/uro/
```
