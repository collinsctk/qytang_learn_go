### 先get依赖包
```shell
[root@jiaozhu-rocky snmp_3_getall]# go get -d ./...
go: downloading github.com/influxdata/influxdb1-client v0.0.0-20220302092344-a9ab5670611c
```
### 编译
```shell
[root@jiaozhu-rocky snmp_3_getall]# env GOOS=linux GOARCH=amd64 go build -o snmp_3_getall_main snmp_3_getall_main.go
```
### 运行
```shell
[root@jiaozhu-rocky snmp_3_getall]# ls -an
total 7720
drwxr-xr-x 6 0 0    4096 Jun  1 08:19 .
drwxr-xr-x 5 0 0      67 May 31 20:52 ..
-rwxr-xr-x 1 0 0 7867992 Jun  1 08:12 snmp_3_getall_main

[root@jiaozhu-rocky snmp_3_getall]# ./snmp_3_getall_main
```
