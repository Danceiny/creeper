# leads-factory
创造leads

## creeper
crawler for `leads-factory`  

## Run
修改`run_client`和`run_worker`两个脚本中的环境变量: `REDIS_HOST`。

更多的环境变量设置（如redis的端口号），请参考`config.go`

```bash
# chmod +x ./build
./build

# session 1
# chmod +x ./run_worker
./run_worker --proxy

# session 2
# chmod +x ./run_client
./run_client
```
