module go-blog

go 1.13

require (
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/elazarl/goproxy v0.0.0-20210110162100-a92cc753f88e // indirect
	github.com/gin-gonic/gin v1.4.0
	github.com/go-errors/errors v1.0.1
	github.com/go-redis/redis v6.15.2+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/ikeikeikeike/go-sitemap-generator/v2 v2.0.2
	github.com/microcosm-cc/bluemonday v1.0.2
	github.com/mojocn/base64Captcha v0.0.0-20190509095025-87c9c59224d8
	github.com/parnurzeal/gorequest v0.2.16
	github.com/pkg/errors v0.9.1 // indirect
	github.com/qiniu/go-sdk/v7 v7.9.7
	github.com/satori/go.uuid v1.2.0
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/speps/go-hashids v2.0.0+incompatible
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	gopkg.in/russross/blackfriday.v2 v2.0.0-00010101000000-000000000000
	gopkg.in/yaml.v2 v2.2.2
	moul.io/http2curl v1.0.0 // indirect
	xorm.io/xorm v1.1.2
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday v2.0.0+incompatible

replace gopkg.in/yaml.v2 => github.com/go-yaml/yaml v0.0.0-20181115110504-51d6538a90f8
