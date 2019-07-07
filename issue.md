@[TOC](TCP/IP如何确保网络通讯质量)

# 操作步骤

## geth

    - geth --datadir ./data --networkid 15 --port 30303 --rpc --rpcaddr 0.0.0.0 --rpcport 8545 --rpcvhosts "*" --rpcapi "db,net,eth,web3,personal" --rpccorsdomain "*" --ws --wsaddr "localhost" --wsport "8546" --wsorigins "*" --nat "any" --nodiscover --dev --dev.period 1 console 2> 1.log

## mysqld server

    mysqld >> src/mysql.log

## mysql client

    - mysql -u root -p
    - 需要重新构建mysql数据库表

## run go

    go run main.go -c etc/copyright.dev.toml

## bytes32

    0xee440896028860593a2659daddb5817701dd9e04cc1e0f190120b1fe1e44d79a

### 项目作业

    设计一个酒店管理系统?要做到数据保护和身份验证?比如酒店前台人员只需要关注入住人年龄是否已满18岁?有无犯罪记录即可?而不需要获取全部的个人信息。要求将个人数据信息上链(关键数据即可?合约一)?再提供给酒店一个查询是否可以入住的接口(合约二?返回true或者false)?使用go(js或其他语言皆可)调用查询接口?无需前端页面。
