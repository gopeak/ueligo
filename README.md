
# morego 实时通信服务器


##安装：

    1.安装go环境,要求go1.6以上
    2.安装git
    3.执行 go get -u github.com/kardianos/govendor ，进入根目录后执行
           govendor init 
           govendor add +external
           govendor sync 
    4.在mysql中创建webim数据库并导入webim.sql  ,修改worker/golang/db.toml连接配置
    5.执行 go run manager.go
    6. demo webim   http://localhost:9898/im



##todo：
    1.vendor 依赖包引入              √
    2.go php python lua workerserver服务  
    3.workerserver 状态维护 
    4.websocket 支持                 √
    5.redis 数据对象持久化           
    6.websocket                      √
    7.错误日志写入mongodb            √
    8.单机版                         √
    9.demo                          √
    10.文档
