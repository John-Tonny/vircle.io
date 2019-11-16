module vircle.io/mnhost/web

go 1.13

require (
	github.com/astaxie/beego v1.12.0
	github.com/garyburd/redigo v1.6.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/protobuf v1.3.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/micro/go-micro v1.13.2
	github.com/shiena/ansicolor v0.0.0-20151119151921-a422bbe96644 // indirect
	vircle.io/mnhost/vps/model v0.0.0-00010101000000-000000000000
	vircle.io/mnhost/vps/utils v0.0.0-00010101000000-000000000000
)

replace github.com/John-Tonny/micro/vps/amazon => /root/mygo/src/github.com/John-Tonny/micro/vps/amazon

replace vircle.io/mnhost/vps/model => /root/mygo/src/vircle.io/mnhost/vps/model

replace vircle.io/mnhost/vps/utils => /root/mygo/src/vircle.io/mnhost/vps/utils

replace vircle.io/mnhost/vps/proto/vps => /root/mygo/src/vircle.io/mnhost/vps/proto/vps
