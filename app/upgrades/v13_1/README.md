### 前提

* 将验证人的私钥导入到keyring-test目录：`./evmosd keys unsafe-import-eth-key ${NAME} ${PRIVATEKEY} --home ./nodes/node0/evmosd --keyring-backend test` 请替换${NAME} ${PRIVATEKEY}，下面的例子${NAME}为node0

### 增加质押操作

* 领取质押奖励以及税费：`./evmosd tx distribution withdraw-rewards evmosvaloper1hajh6rhhkjqkwet6wqld3lgx8ur4y3khljfx82 --commission --gas="600000" --gas-prices="10000000000aevmos" --from=node0 --home=./nodes/node0/evmosd/ --keyring-backend=test --broadcast-mode sync -y` 请替换验证者地址以及其他相关参数。

* 质押：`./evmosd tx staking delegate evmosvaloper1hajh6rhhkjqkwet6wqld3lgx8ur4y3khljfx82 100000000000000000000aevmos --gas="600000" --gas-prices="10000000000aevmos" --from=node0 --home=./nodes/node0/evmosd/ --keyring-backend=test --broadcast-mode sync -y` 请替换验证者地址、质押数量以及其他相关参数。

### 修改投票周期

* 查看投票周期：`curl http://127.0.0.1:1317/cosmos/gov/v1beta1/params/voting` 返回的 `voting_params.voting_period` 就是投票周期

* 新建一个修改参数提案文件proposal_change_params.json，注意时间用的是纳秒表示。下面的提案的投票时间设为10分钟，那么填的数据为 10 * 60 * 1000000000 = 600000000000。提案内容如下。
```json
{
  "messages": [
    {
      "@type": "/cosmos.gov.v1.MsgExecLegacyContent",
      "content": {
        "@type": "/cosmos.params.v1beta1.ParameterChangeProposal",
        "title": "Parameter change: voting period.",
        "description": "decrease voting period time to 10 minutes",
        "changes": [
          {
            "subspace": "gov",
            "key": "votingparams",
            "value": "{\"voting_period\":\"600000000000\"}"
          }
        ]
      },
      "authority": "evmos10d07y265gmmuvt4z0w9aw880jnsr700jcrztvm"
    }
  ],
  "metadata": "AQ==",
  "deposit": "10000000aevmos"
}
```

* 发起提案：`./evmosd tx gov submit-proposal ./proposal_change_params.json --gas="600000" --gas-prices="10000000000aevmos" --from=node0 --home=./nodes/node0/evmosd/ --keyring-backend=test --broadcast-mode sync -y`

* 给提案投票：`./evmosd tx gov vote 1 yes --gas="600000" --gas-prices="10000000000aevmos" --from=node0 --home=./nodes/node0/evmosd/ --keyring-backend=test --broadcast-mode sync -y` 请替换vote后面的提案ID。

### 发起升级提案

* 新建一个升级提案文件proposal_upgrade.json，注意需要根据实际情况修改height的内容，简单计算一下在提案结束后的高度。举个简单的例子：假设现在每1s出一个区块，发起提案的高度为100，提案的时间为10分钟，则设置的height > 100 + 10 * 60，当然最好设置的稍微再大点，比如设置为800。如果提案未结束就达到了该高度，则升级会失败。提案内容如下：
```json
{
  "messages": [
    {
      "@type": "/cosmos.gov.v1.MsgExecLegacyContent",
      "authority": "evmos10d07y265gmmuvt4z0w9aw880jnsr700jcrztvm",
      "content": {
        "@type": "/cosmos.upgrade.v1beta1.SoftwareUpgradeProposal",
        "title": "Software upgrades: reset evm chain id",
        "description": "reset evm chain id from 9000 to 168168",
        "plan": {
          "name": "v13.1.0",
          "info": "{\"binaries\":{\"darwin/arm64\":\"https://github.com/evmos/evmos/releases/download/v16.0.0/evmos_16.0.0_Darwin_arm64.tar.gz?checksum=db85391340c5d6eaec460d355b8ca000dd0ed06b6f2338d24b20eb3ab32d745a\",\"darwin/amd64\":\"https://github.com/evmos/evmos/releases/download/v16.0.0/evmos_16.0.0_Darwin_amd64.tar.gz?checksum=98c9074e81998cb939c192a9982de93119657e7f7020982eb7629e45bd951a04\",\"linux/arm64\":\"https://github.com/evmos/evmos/releases/download/v16.0.0/evmos_16.0.0_Linux_arm64.tar.gz?checksum=ecfffdaadd0e05058b78a2c24bdcb763975e95236eb4a63e72c0be7a0ad30b6d\",\"linux/amd64\":\"https://github.com/evmos/evmos/releases/download/v16.0.0/evmos_16.0.0_Linux_amd64.tar.gz?checksum=dc55e04c7f12768fb32beb1d06f3d113e76059b76d1ac5f13657c8fccf5fc309\"}}",
          "height": "120",
          "time": "0001-01-01T00:00:00Z",
          "upgraded_client_state": null
        }
      }
    }
  ],
  "metadata": "AQ==",
  "deposit": "10000000aevmos"
}
```

* 发起提案：`./evmosd tx gov submit-proposal ./proposal_upgrade.json --gas="600000" --gas-prices="10000000000aevmos" --from=node0 --home=./nodes/node0/evmosd/ --keyring-backend=test --broadcast-mode sync -y`

* 给提案投票：`./evmosd tx gov vote 1 yes --gas="600000" --gas-prices="10000000000aevmos" --from=node0 --home=./nodes/node0/evmosd/ --keyring-backend=test --broadcast-mode sync -y` 请替换vote后面的提案ID。

* 等待提案通过。提案通过之后，到了指定高度所有节点都会停下来等待客户端的替换。此时表现为验证节点不再出块。使用分支为`depass/reset-evm-chain-id`的代码重新编译一个客户端，将所有的节点（验证节点以及非验证节点）旧的客户端程序停止，替换成新的客户端，验证节点就会恢复出块。

