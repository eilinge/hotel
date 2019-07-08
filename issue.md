@[TOC](TCP/IP如何确保网络通讯质量)

# 操作步骤

## geth

    geth --datadir ./data --networkid 15 --port 30303 --rpc --rpcaddr 0.0.0.0 --rpcport 8545 --rpcvhosts "*" --rpcapi "db,net,eth,web3,personal" --rpccorsdomain "*" --ws --wsaddr "localhost" --wsport "8546" --wsorigins "*" --nat "any" --nodiscover --dev --dev.period 1 console 2> 1.log

## mysqld server

    mysqld >> src/mysql.log

## mysql client

    - mysql -u root -p
    - 需要重新构建mysql数据库表

## run go

    go run main.go -c etc/copyright.dev.toml

## Data

    bytes32
        0xee440896028860593a2659daddb5817701dd9e04cc1e0f190120b1fe1e44d79a
    asset(信息登记)
        http://localhost:8086/mint?name=%22eilinge%22&&age=17&&idcard=431222199112133117&&reword=0
        http://localhost:8086/mint?name=%22lisi%22&&age=19&&idcard=431222199112133118&&reword=1
        http://localhost:8086/mint?name=%22lisi%22&&age=19&&idcard=431222199112133119&&reword=0
    check(查询)
        http://localhost:8086/check?idcard=431222199112133118

### 项目作业

    设计一个酒店管理系统?要做到数据保护和身份验证?比如酒店前台人员只需要关注入住人年龄是否已满18岁?有无犯罪记录即可?而不需要获取全部的个人信息。要求将个人数据信息上链(关键数据即可?合约一)?再提供给酒店一个查询是否可以入住的接口(合约二?返回true或者false)?使用go(js或其他语言皆可)调用查询接口?无需前端页面。

    1. 不需要登录注册, 也不需要数据库存储
    2. 通过/mint?idcard...进行信息上链(get)
    3. 通过/check?idcard...进行信息查询(get)
    4. 使用身份证号存储游客, 目前功能不完善, 可以使用相同身份证号进行存储, 有待改进

    返回数据
        {"errno":"0","errmsg":"成功","data":{"Age": 19,"CriminalRecord": false,"checkin": true}}

        "checkin"--表示能否成功入住 true:可以入住 false: 无法入住
