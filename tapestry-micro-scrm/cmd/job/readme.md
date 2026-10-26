## 定时任务

### 统计重复客户(全量)
    第一次上线需要手动执行一次，此后每天凌晨3点执行一次.本地执行：
    go run cmd\job\main.go -config app.yaml -job.cmd=RepeatCustomerCount -job.spec=企业ID

### 给relation_callback_records 的存量删除关系的数据添加关系建立时间
    上线手动执行一次即可：
    go run cmd\job\main.go -config app.yaml -job.cmd=LossCustomerAddRelationTimeJob -job.spec=企业ID

### 更改欢迎语content的结构
    上线手动执行一次即可：
    go run cmd\job\main.go -config app.yaml -job.cmd=FixDrainageWelcomeContentJob -job.spec=企业ID
