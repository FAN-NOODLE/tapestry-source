package job

import "time"

const (
	// TASK_LISTEN_INTERVAL 任务监测间隔时间
	TASK_LISTEN_INTERVAL = 5 * time.Second

	// PAGE_SIZE 每次获取的基础数据的量
	PAGE_SIZE = 1000

	// EXTRACT_CUSTOMER_NUM 在职继承是每次提取的客户数量
	EXTRACT_CUSTOMER_NUM
)
