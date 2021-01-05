# 分布式系统HDFS
### 启动服务
* 以introducer加入

运行命令 `go run *.go -intro`
* 以一般成员加入

运行命令 `go run *.go -introIp=123.123.10.1`

* 以gossip心跳机制加入（默认不带flag为all-to-all心跳机制）

运行命令 `go run *.go -introIp=123.123.10.1 -gossip`


### 命令
* 从HDFS获取文件

运行命令 `get [local file name] [hdfs file name]`。其中，local file name为下载到本地的文件名，hdfs file name为文件存储在hdfs中的名字。

* 上传文件到HDFS

运行命令 `put [local file name] [hdfs file name]`。其中，local file name为待上传的本地文件名，hdfs file name为文件存储在hdfs中的名字。

* 查找在哪些VM上有文件的replicas

运行命令 `ls [hdfs file name]`。 其中，hdfs file name为文件存储在hdfs中的名字。显示结果为存储VM的ID。

* 显示本机上存有哪些HDFS文件

运行命令 `store`。 
