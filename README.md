# Mock Test

本项目旨在一次性生成多个已签名的数据，此生成过程可大致分为两步：

1. 初始化一个大账户（mock faucet account）用于给账户转账
2. 批量生成账户并生成已签名的交易数据

## Init mock faucet account

**Command**

```bash
mock faucet-init --seed="recycle light kid spider fire disorder relax end stool hip child leaf wild next veteran start theory pretty salt rich avocado card enact april"
```

**Parameters**

- `seed`：大账户的助记词

## Gen signed tx data

**Command**

```bash
mock gen-signed-tx --num 20 --receiver faa1t5wlur60xzzcxpgjn0d5y8ge7fsdmp7jejl7am --faucet faa1jyj90se9mel2smn3vr4u9gzg03acwuy8h44q3m --chain-id=rainbow-dev --node http://localhost:1317
```

**Parameters**

- `num`：需要生成已签名交易的数量
- `receiver`：交易接收方 address
- `faucet`： 大账户地址（faucet address），`faucet-init` 命令输出的结果
- `chain-id`：chain id
- `node`：lcd 接口地址