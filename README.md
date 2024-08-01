# Token Payment

`token-payment` 是一个区块链收提币系统，旨在为各种应用提供的收款和转账功能。它支持所有 EVM 兼容的区块链，并计划在未来支持更多的链。该系统提供了到账检测、到账通知和转账等功能，通过 API 接口可以方便地集成到您的应用中。

## 特性

- **多链支持**：当前仅支持 EVM 兼容链。
- **到账检测**：实时检测和确认到账情况。
- **到账通知**：通过通知机制及时告知用户到账状态。
- **转账功能**：支持从一个地址向另一个地址进行代币转账。
- **API 接口**：提供简单易用的 API 接口，方便集成到各种应用中。

## 快速开始

### 安装

要开始使用 `token-payment`，您可以从 GitHub 克隆代码库并安装相关依赖。

```bash
git clone https://github.com/seanlan/token-payment.git
cd token-payment
go mod tidy
go run main.go cron # 启动定时任务
go run main.go web  # 启动 Web 、 API 服务
```

### 配置

导入数据库文件 `init/init.sql` 到您的数据库中，以创建必要的数据表。

需要配置conf.yaml, 修改conf.yaml.example为conf.yaml
可以在 `conf.yaml` 文件中设置这些变量：

```yaml
# 这里是mysql数据库连接地址
db:
  uri: mysql://xxxx:xxxx@tcp(127.0.0.1:3306)/tokenpay?parseTime=true&loc=Local&charset=utf8mb4&collation=utf8mb4_unicode_ci
# 这里是redis数据库连接地址
redis:
  hosts: "127.0.0.1:6379"
  prefix: "token-payment"
  username: ""
  password: ""
  db: 0
```
配置支持的链，在chain表中插入一条数据
```sql
INSERT INTO `chain` (`id`, `chain_symbol`, `name`, `chain_id`, `currency`, `chain_type`, `confirm`, `gas`, `gas_price`, `latest_block`, `rebase_block`, `has_branch`, `concurrent`, `address_pool`, `watch`)
VALUES
	(1, 'amoy', 'Polygon Amoy', 80002, 'matic', 'evm', 80, 500000, 50000000000, 10191842, 10191668, 1, 20, 100, 1);
```
配置链的rpc地址，在chain_rpc表中插入一条数据
```sql
INSERT INTO `chain_rpc` (`id`, `chain_symbol`, `rpc_url`, `disable`)
VALUES
	(1, 'amoy', 'https://polygon-amoy-bor-rpc.publicnode.com', 0);
```
这里使用的是公共的rpc地址，由于会限制请求频率和请求次数，建议多配置几个rpc地址。

这里提供几个有免费额度的rpc地址，欢迎补充：
- http://publicnode.com/
- https://www.ankr.com/
- https://blockpi.io/

系统使用了`internal/chain/equalizer.go` 实现了rpc地址的分配，可以根据实际情况修改`internal/chain/eqchain`中的`Get`方法，实现自己的分配策略。

### 使用示例

在使用 `token-payment` 之前，您需要创建一个application。

在application表中插入一条数据
```sql
INSERT INTO `application` (`id`, `app_key`, `app_secret`, `app_name`, `hook_url`, `create_at`)
VALUES
	(1, 'f4399f851e984405aa1eba51ecbce790', 'TQYAHn6A961zzGKjfb99', '测试', '', 1697176794);
```

设置这个application的不同链的配置`application_chain`，主要是配置
- `hot_wallet` 热钱包地址、用于提币
- `cold_wallet` 冷钱包地址、应用下所有收到的代币都会转到这个地址
- `fee_wallet` 手续费地址、整理零钱到冷钱包时，会给整理的地址转手续费


#### 收款

实现收款功能，需要先创建一个收款地址，然后将这个地址展示给用户。
这个地址是一个临时地址，收到代币到达一定阀值后会自动转到冷钱包（功能暂未实现）。
这个地址每次收到代币后，会调用`notify_url`通知应用。

生成收款地址，可以调用：
```http
POST /api/v1/payment/address
Content-Type: application/json

{
  "app_key": "Application app_key",
  "data": "Create Address PARAMS",
  "sign": "SIGN"
}
```

- `app_key`：应用的 app_key。
- `data`：创建地址的参数，这里是一个json字符串例如：
  ```json
  {
    "chain": "amoy",
    "notify_url": "https://your-notify-url.com"
  }
  ```
  -  `chain`：链的标识。
  -  `notify_url`：收到代币后通知的地址。
- `sign`：参数签名 MD5(data + app_secret).ToUpper()。

响应内容
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "address": "0x"
  }
}
```


#### 转账

这里是一个异步的操作，实际上是将转账请求放到队列中，由定时任务进行处理。
到账后会通知应用。 调用：

```http
POST /api/v1/payment/withdraw
Content-Type: application/json

{
  "app_key": "Application app_key",
  "data": "Create Address PARAMS",
  "sign": "SIGN"
}
```
- `app_key`：应用的 app_key。
- `data`：创建地址的参数，这里是一个json字符串例如：
  ```json
  {
    "chain": "amoy",
    "serial_no": "serial_no",
    "symbol": "symbol",
    "contract_address": "contract_address",
    "token_id": 0,
    "value": "0.123123",
    "to_address": "0x",
    "notify_url": "https://your-notify-url.com/notify"
    }
    ```
  - `chain`：链的标识。
  - `serial_no`：交易流水号 (唯一)，回调时会返回，用于标识交易。
  - `symbol`：代币符号 与`chain_token`表中的`symbol`对应
  - `contract_address`：代币合约地址 与`chain_token`表中的`contract_address`对应。非必要参数，防止symbol重复时使用的情况
  - `token_id`：NFT的token_id，非NFT时为0。（暂未实现NFT的转账功能）
  - `value`：转账金额。使用字符串表示，例如："0.123123"。
  - `to_address`：转账目标地址。
  - `notify_url`：转账完成后通知的地址。

响应内容
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "exist": false
  }
}
```

#### 到账通知
完成交易后，系统会根据应用提供的`notify_url`地址通知应用。
包括收款地址的到账通知，转账申请的到账通知。

```http
POST NOTIFY_URL
Content-Type: application/json

{
  "data": "Create Address PARAMS",
  "sign": "SIGN"
}
```
- `data`：通知的数据，这里是一个json字符串例如：
  ```json
  {
      "application_id": 12345,
      "chain_symbol": "ETH",
      "tx_hash": "0xabc123...",
      "from_address": "0x12345...",
      "to_address": "0x67890...",
      "contract_address": "0xcontract...",
      "symbol": "USDT",
      "decimals": 18,
      "token_id": 10001,
      "value": 100.0,
      "tx_index": 1,
      "batch_index": 0,
      "confirm": 3,
      "max_confirm": 12,
      "transfer_type": 1,
      "serial_no": "TX20230801",
      "create_at": 1690905600
  }
  ```
    - `application_id`：应用ID。
    - `chain_symbol`：链的标识。
    - `tx_hash`：交易哈希。
    - `from_address`：转出地址。
    - `to_address`：转入地址。
    - `contract_address`：代币合约地址。
    - `symbol`：代币符号。
    - `decimals`：代币精度。
    - `token_id`：NFT的token_id。
    - `value`：转账金额。
    - `tx_index`：交易索引。
    - `batch_index`：批次索引。
    - `confirm`：确认数。
    - `max_confirm`：最大确认数。
    - `transfer_type`：转账类型。 1: 收款 2: 转账
    - `serial_no`：交易流水号。
    - `create_at`：创建时间。
- `sign`：参数签名 MD5(data + app_secret).ToUpper()。

成功处理后返回 `success`。
如没有正确的返回 `success`，系统会在一段时间后重试通知。

通知的间隔时间为

| 失败次数 | 重试间隔 |
|------|------|
| 1    | 1min |
| 2    | 2min |
| 3    | 4min |
| 4    | 8min |

以此类推，最多重试 10 次。

## 开发计划

- [ ] 零钱整理。收款地址到账后自动转到冷钱包
- [ ] 管理后台
- [ ] 数据统计
- [ ] 其它非EVM链支持
- [ ] NFT转账功能

关于其它链的功能，可以通过实现`internal/chain/chain.go`接口来实现。
```go
type BaseChain interface {
	GetLatestBlockNumber(ctx context.Context) (int64, error)               // 获取最新区块
	GetBlock(ctx context.Context, number int64) (*Block, error)            // 获取区块
	GetTransaction(ctx context.Context, hash string) (*Transaction, error) // 获取交易
	GenerateAddress(ctx context.Context) (string, string, error)           // 生成地址
	GetNonce(ctx context.Context, address string) (uint64, error)          // 获取nonce
	GenerateTransaction(ctx context.Context, order *TransferOrder) error   // 生成交易订单
	Transfer(ctx context.Context, order *TransferOrder) (string, error)    // 转账
}
```

具体的实现可以参考`internal/chain/evmchain.go`。

## 问题反馈

如果您在使用过程中遇到问题，请在 [GitHub Issues](https://github.com/yourusername/token-payment/issues) 中报告。
