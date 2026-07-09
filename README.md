# driver
Device driver for Self-Service Terminal project


Build executables
```shell
go build -o manager ./cmd/manager
go build -o equipment ./cmd/equipment
```

Run manager application
```shell
./manager
```

Run equipment application
```shell
./equipment -dev_name device -log_path ./log -cfg_path ./cfg

```
