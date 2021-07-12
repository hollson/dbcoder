# dbcoder
Golang-数据库代码生成器


### 快速使用

将项目放入到GOPATH/src目录下，进入项目根目录下进行

```shell
$ go build
```

对于不同的操作系统

`windows`

```shell
$ dbcoder.exe -host=127.0.0.1 -port=5432 -user=postgres -pwd=postgres  -dbname=db_test -gorm=true -driver=pgsql
```

`linux`

```shell
$ ./dbcoder -host=127.0.0.1 -port=3306 -user=root -pwd=root -dbname=db_test -gorm=true -driver=mysql -package=hello
```

### 命令行提示

执行

```shell
$ ./dbcoder -help
```



```powershell
Usage of dbcoder.exe:
  -dbname string
        必填，数据库名称，否则会报错
  -driver string
        必填，需要连接的数据库，现在只支持mysql、pgsql 例如 -driver=mysql，-driver=pgsql
  -gorm
        选填，是否添加 gorm tag，true添加，false不添加，默认不添加
  -host string
        选填，数据库ip，默认为localhost (default "localhost")
  -outdir string
        选填，go 文件输出路径，不设置的话会输出到当前程序所在路径 (default "./go_output")
  -package string
        选填，go 文件中 package 的名字，默认为 package main (default "main")
  -port int
        必填，数据库端口
  -pwd string
        必填，数据库密码
  -table string
        选填，需要导出的数据库表名称，如果不设置的话会将该数据库所有的表导出
  -user string
        必填，数据库用户名
```