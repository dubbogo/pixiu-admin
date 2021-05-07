# pixiu-admin
pixiu admin


## prepare etcd config

run etcd local or in docker, then use etcdctl to set api config

- export ETCDCTL_API=3
- if use docker, run `docker cp api_config.yaml mycontainer:/path` to copy file to docker
- run `cat api_config.yaml | etcdctl put "/proxy/config/api"` to set api config

## Start admin

run cmd/admin/admin.go

config program arguments：
- -c /xx/pixiu-admin/configs/admin_config.yaml
