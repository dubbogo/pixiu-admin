# pixiu-admin
pixiu admin


## prepare etcd config

run etcd local or in docker

## Start admin

run cmd/admin/admin.go

config program arguments：
- -c /xx/pixiu-admin/configs/admin_config.yaml

## operation

Response format:
```
{
"code":
"data": 
}

code：
10001 success
10002 not found data
10003 concurrent operate, fresh and try again
```


#### query resource list
```
curl -X GET \
http://127.0.0.1:8080/config/api/resource/list \
-H 'Postman-Token: 98f218e4-456e-4e60-abe8-a6554fadc57b' \
-H 'cache-control: no-cache'
```
Result:
- resource config list
- or error when not set any resource config

#### set resource config

```
curl -X POST \
  http://127.0.0.1:8080/config/api/resource/ \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: c9ebd52f-0ef0-4365-8383-64320f9786fa' \
  -H 'cache-control: no-cache' \
  -H 'content-type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW' \
  -F 'content=    path: '\''/api/v1/test-dubbo/friend2'\''
    type: restful
    description: user
    timeout: 100ms
    plugins:
      pre:
        pluginNames:
          - rate limit
          - access
      post:
        groupNames:
          - group2
    methods:
      - httpVerb: GET
        resourcePath: '\''/api/v1/test-dubbo/friend2'\''
        onAir: true
        timeout: 1000ms
        inboundRequest:
          requestType: http
          queryStrings:
            - name: name
              required: true
        integrationRequest:
          requestType: http
          host: 127.0.0.1:8889
          path: /UserProvider/GetUserByName
          mappingParams:
            - name: queryStrings.name
              mapTo: queryStrings.name
          group: "test"
          version: 1.0.0'
```

#### query resource detail by id

```
curl -X GET \
  'http://122.51.143.73:8080/config/api/resource/detail?resourceId=1' \
  -H 'Postman-Token: 3a7612e4-7837-4bf0-afd5-8429d3ee5408' \
  -H 'cache-control: no-cache'
```

#### modify resource config

```
curl -X PUT \
  http://127.0.0.1:8080/config/api/resource \
  -H 'Postman-Token: b92d07d3-f2e5-43d3-a3b1-b108a6de20b2' \
  -H 'cache-control: no-cache' \
  -H 'content-type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW' \
  -F 'content=    id: 1
    path: '\''/api/v1/test-dubbo/friend'\''
    type: restful
    description: update
    timeout: 1000ms
    plugins:
      pre:
        pluginNames:
          - rate limit
          - access
      post:
        groupNames:
          - group2
    methods:
      - httpVerb: GET
        onAir: true
        timeout: 1000ms
        inboundRequest:
          requestType: http
          queryStrings:
            - name: name
              required: true
        integrationRequest:
          requestType: http
          host: 127.0.0.1:8889
          path: /UserProvider/GetUserByName
          mappingParams:
            - name: queryStrings.name
              mapTo: queryStrings.name
          group: "test"
          version: 1.0.0'
```
#### delete resource config

```
curl -X DELETE \
  'http://122.51.143.73:8080/config/api/resource/?resourceId=1' \
  -H 'Postman-Token: dd32793a-844c-4edf-b9b1-26b4d2b082f8' \
  -H 'cache-control: no-cache'
```

#### query method list below one resource

```
curl -X GET \
  'http://127.0.0.1:8080/config/api/resource/method/list?resourceId=1' \
  -H 'Postman-Token: 77bb4f7d-b21e-4ef4-a019-64732d9464a2' \
  -H 'cache-control: no-cache'
```

#### create method config below one resource

```
curl -X POST \
  'http://127.0.0.1:8080/config/api/resource/method/?resourceId=1' \
  -H 'Postman-Token: f0ab0783-2090-428e-98d3-01c97308fa3e' \
  -H 'cache-control: no-cache' \
  -H 'content-type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW' \
  -F 'content=httpVerb: PUT
resourcePath: '\''/api/v1/test-dubbo/friend'\''
onAir: true
timeout: 1000ms
inboundRequest:
    requestType: http
    queryStrings:
    - name: name
      required: true
integrationRequest:
    requestType: http
    host: 127.0.0.1:8889
    path: /UserProvider/GetUserByName
    mappingParams:
    - name: queryStrings.name
      mapTo: queryStrings.name
    group: "test"
    version: 1.0.0'
```

#### query method detail 

```
curl -X GET \
  'http://122.51.143.73:8080/config/api/resource/method/detail?resourceId=1&methodId=2' \
  -H 'Postman-Token: 740d91fb-2d06-49f9-b900-948235f40a86' \
  -H 'cache-control: no-cache'
```

#### delete method 
```
curl -X DELETE \
  'http://127.0.0.1:8080/config/api/resource/method/?resourceId=1&methodId=2' \
  -H 'Postman-Token: 51f3975c-4e12-434a-abfc-7be1e6ec82b9' \
  -H 'cache-control: no-cache'
```


#### other
also support below operation:
- modify method config
- manage plugin ratelimit config
- manage plugin group config