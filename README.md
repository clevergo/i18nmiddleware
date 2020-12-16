# I18N Middleware for CleverGo
[![Build Status](https://img.shields.io/travis/clevergo/i18nmiddleware?style=flat-square)](https://travis-ci.org/clevergo/i18nmiddleware)
[![Coverage Status](https://img.shields.io/coveralls/github/clevergo/i18nmiddleware?style=flat-square)](https://coveralls.io/github/clevergo/i18nmiddleware)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/clevergo.tech/i18nmiddleware?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/clevergo/i18nmiddleware?style=flat-square)](https://goreportcard.com/report/github.com/clevergo/i18nmiddleware)
[![Release](https://img.shields.io/github/release/clevergo/i18nmiddleware.svg?style=flat-square)](https://github.com/clevergo/i18nmiddleware/releases)
[![Downloads](https://img.shields.io/endpoint?url=https://pkg.clevergo.tech/api/badges/downloads/total/clevergo.tech/i18nmiddleware&style=flat-square)](https://pkg.clevergo.tech/clevergo.tech/i18nmiddleware)
[![Chat](https://img.shields.io/badge/chat-telegram-blue?style=flat-square)](https://t.me/clevergotech)
[![Community](https://img.shields.io/badge/community-forum-blue?style=flat-square&color=orange)](https://forum.clevergo.tech)

```shell
$ go get -u clevergo.tech/i18nmiddleware
```

## Usage

```go
package main

import (
	"net/http"

	"clevergo.tech/clevergo"
	"clevergo.tech/i18nmiddleware"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func main() {
	app := clevergo.New()
	bundle := i18n.NewBundle(language.English)
	bundle.ParseMessageFileBytes([]byte(`{"home": "Home"}`), "en.json")
	bundle.ParseMessageFileBytes([]byte(`{"home": "主页"}`), "zh-CN.json")
	app.Use(i18nmiddleware.New(bundle))
	app.Get("/", func(c *clevergo.Context) error {
		localizer := i18nmiddleware.Localizer(c)
		s, _, _ := localizer.LocalizeWithTag(&i18n.LocalizeConfig{
			MessageID: "home",
		})
		return c.String(http.StatusOK, s)
	})
	app.Run(":8080")
}
```

```shell
$ curl http://localhost:8080/
Home

$ curl http://localhost:8080/?lang=zh-CN
主页

$ curl -H "Accept-Language: zh-CN" http://localhost:8080
主页
```
