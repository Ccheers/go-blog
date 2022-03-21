module go-blog

go 1.17

require (
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.4.0
	github.com/go-errors/errors v1.0.1
	github.com/go-redis/redis v6.15.2+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/microcosm-cc/bluemonday v1.0.2
	github.com/mojocn/base64Captcha v0.0.0-20190509095025-87c9c59224d8
	github.com/parnurzeal/gorequest v0.2.16
	github.com/qiniu/go-sdk/v7 v7.9.7
	github.com/satori/go.uuid v1.2.0
	github.com/snabb/sitemap v1.0.0
	github.com/speps/go-hashids v2.0.0+incompatible
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	gopkg.in/russross/blackfriday.v2 v2.0.0-00010101000000-000000000000
	gopkg.in/yaml.v2 v2.2.2
	xorm.io/xorm v1.1.2
)

require (
	github.com/elazarl/goproxy v0.0.0-20210110162100-a92cc753f88e // indirect
	github.com/gin-contrib/sse v0.0.0-20190301062529-5545eab6dad3 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/protobuf v1.3.1 // indirect
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/snabb/diagio v1.0.0 // indirect
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/ugorji/go v1.1.4 // indirect
	golang.org/x/image v0.0.0-20190501045829-6d32002ffd75 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a // indirect
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
	moul.io/http2curl v1.0.0 // indirect
	xorm.io/builder v0.3.8 // indirect
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday v2.0.0+incompatible

replace gopkg.in/yaml.v2 => github.com/go-yaml/yaml v0.0.0-20181115110504-51d6538a90f8
