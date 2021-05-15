module github.com/dubbogo/pixiu-admin

go 1.14

require (
	github.com/apache/dubbo-getty v1.4.3
	github.com/dubbogo/dubbo-go-pixiu-filter v0.1.3
	github.com/dubbogo/gost v1.11.8
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/gin-gonic/gin v1.7.1
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	github.com/urfave/cli v1.22.4
	go.etcd.io/etcd v0.0.0-20200402134248-51bdeb39e698
	go.uber.org/zap v1.16.0
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/dubbogo/dubbo-go-pixiu-filter => ../dubbo-go-pixiu-filter
