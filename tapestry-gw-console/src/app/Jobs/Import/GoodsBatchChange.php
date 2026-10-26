<?php


namespace App\Jobs\Import;


use App\Jobs\Job;
use App\Services\GoodsBatchChangeService;
use CrsSdk\MicroGoodsV1\Goods;
use Illuminate\Support\Facades\Redis;
use Throwable;

class GoodsBatchChange extends Job
{
    /**
     * @var string
     */
    public $queue = 'console-goods-batch-change';

    /**
     * @var int 重试1此
     */
    public $tries = 1;

    /**
     * @var int $enterpriseId
     */
    protected $enterpriseId;

    /**
     * @var int $type
     */
    protected $type;

    protected $redis;

    protected $hash;

    public function __construct(int $enterpriseId, int $type, string $hash)
    {
        $this->enterpriseId = $enterpriseId;
        $this->type = $type;
        $this->hash = $hash;
    }

    public function handle(): void
    {
        $_SERVER['enterprise_hash'] = $this->hash;
        $this->redis = Redis::connection();
        $sdk = new Goods();
        $service = new GoodsBatchChangeService($this->enterpriseId, $this->type);
        $data = json_decode($this->redis->get($service->getSaveListKey()), true);
        $page = 1;
        $size = 20;
        $current = 0;
        $count = count($data);
        $failed = [];
        startGet: $offset = $page === 1 ? 0 : ($page - 1) * $size;
        $pageData = array_slice($data, $offset, $size);
        if ($pageData) {
            try {
                $response = $sdk->goodsBatchChange([
                    'enterprise_id' => $this->enterpriseId,
                    'type'          => $this->type,
                    'data'          => $pageData
                ]);
                $failed = array_merge($failed, $response['errors']);
                $current += count($pageData);
                $this->redis->set($service->getProgressKey(), json_encode($service->makeProgressData($count, $current)));
                // 进度1800秒(半小时)后失效,避免极端情况进度无法释放不能再次导入
                $this->redis->setex($service->getProgressSystemKey(), 1800, json_encode($service->makeProgressData($count, $current)));
                ++ $page;
                goto startGet;
            } catch (Throwable $exception) {
                $this->redis->del($service->getSaveListKey());
                $this->redis->del($service->getProgressSystemKey());
                $this->redis->set($service->getProgressKey(), json_encode($service->makeProgressData($count, $count)));
                $this->redis->set($service->getResultKey(), json_encode($service->makeResultData($service::STATUS_COMPLETE,
                    [
                        '错误原因', '错误详情'
                    ], [
                        ['系统错误', $exception->getMessage()]
                    ], $count)));
                return ;
            }
        }
        $this->redis->set($service->getResultKey(), json_encode($service->makeResultData($service::STATUS_COMPLETE,
            $this->getErrorTitle($service),
            $failed,
            $count)));
        $this->redis->del($service->getProgressSystemKey());
        $this->redis->del($service->getSaveListKey());
    }


    private function getErrorTitle(GoodsBatchChangeService $service): array
    {
        $title = [];
        switch ($this->type) {
            case $service::TYPE_WAIT_SALES_PRICE:
                $title = ['skuId', '价格', '失败原因'];
                break;
            case $service::TYPE_GOODS_NAME_SPU:
                $title = ['spuId', '名称', '失败原因'];
                break;
            case $service::TYPE_GOODS_NAME_SKU:
                $title = ['skuId', '名称', '失败原因'];
                break;
            case $service::TYPE_WAIT_SALES_STOCK_EXPRESS:
                $title = ['skuId', '线上库存', '失败原因'];
                break;
            case $service::TYPE_WAIT_SALES_STOCK_EXTRACT:
                $title = ['skuId', '线下库存', '失败原因'];
                break;
            case $service::TYPE_ALREADY_SALES_STOCK_EXPRESS:
                $title = ['skuId', '线上库存', '失败原因'];
                break;
            case $service::TYPE_ALREADY_SALES_STOCK_EXTRACT:
                $title = ['skuId', '线下库存', '失败原因'];
                break;
        }
        return $title;
    }
}