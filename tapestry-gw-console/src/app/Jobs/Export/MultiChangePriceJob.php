<?php
/**
 * Created by PhpStorm.
 * User: Leon Liu
 * Date: 2020/7/1
 * Time: 21:05
 */


namespace App\Jobs\Export;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Jobs\Job;
use App\Services\GoodsService;
use CrsSdk\MicroGoodsV1\MaterialChannel;
use Illuminate\Support\Facades\Log;

class MultiChangePriceJob extends Job
{
    /**
     * 标记有改价任务在进行中的锁
     */
    const LOCK_MARK = '_MultiPriceChangeLock';

    /**
     * 标记目前总的任务数
     */
    const SUCCESS_NUM = '_MultiPriceChangeSuccessTotal';

    /**
     * 标记目前总的任务数
     */
    const FAIL_NUM = '_MultiPriceChangeFailTotal';

    /**
     * 失败任务的链表
     */
    const FAIL_LIST = '_FailTaskList';

    /**
     * @var array 改价数据
     */
    protected $priceList;

    /**
     * @var string 上传的文件名
     */
    protected $fileName;

    /**
     * @var int 企业id
     */
    protected $enterpriseId;

    /**
     * @var int 操作账号
     */
    protected $enterpriseAdminId;

    /**
     * @var GoodsService
     */
    protected $goodsService;

    /**
     * @var string 企业hash
     */
    protected $hash;

    /**
     * MultiChangePriceJob constructor. 初始化数据
     *
     * @param array $priceList 改价数据
     * @param int $enterpriseId 企业id
     * @param int $enterpriseAdminId 操作账号
     * @param string $hash 企业hash
     * @param string|null $fileName 上传的改价文件名
     */
    public function __construct(
        array $priceList,
        $enterpriseId,
        $enterpriseAdminId,
        $hash,
        string $fileName = null
    ) {
        // 数据初始化：队列名
        $this->queue = Constant::MULTI_CHANGE_PRICE_QUEUE_NAME;

        // 数据初始化：改价数据
        $this->priceList = $priceList;

        // 数据初始化：上传改价模板的文件名
        $this->fileName = empty($fileName)?'待上架商品改价模板':$fileName;

        // 数据初始化：企业id
        $this->enterpriseId = $enterpriseId;

        // 数据初始化：操作人id
        $this->enterpriseAdminId = $enterpriseAdminId;

        // 企业hash
        $this->hash = $hash;

        // 数据初始化：GoodsService
        $this->goodsService = new GoodsService();
    }

    /**
     * 进行改价的业务
     *
     * @return bool
     */
    public function handle()
    {
        // 无修改内容
        if (empty($this->priceList)) {
            return true;
        }

        // 循环改价列表进行改价业务
        $_SERVER['enterprise_hash'] = $this->hash;
        foreach ($this->priceList as $sku => $price) {
            // 循环调取微服务进行改价操作
            $params = [
                'enterprise_id' => $this->enterpriseId,
                'channel_type' => 1,
                'operate_id' => $this->enterpriseAdminId,
                'sku_id' => $sku,
                'on_line_price' => $price,
                'off_line_price' => $price
            ];

            $res = true;
            $errorMsg = '';

            // 调取微服务
            try {
                (new MaterialChannel())->editPrice($params);
            } catch (Error $error) {
                $res = false;
                $errorMsg = $error->getMessage();
                Log::error('修改价格失败' . ' total: ' . count($this->priceList), $error->data);
            } catch (\Exception $exception) {
                $res = false;
                $errorMsg = $exception->getMessage();
                Log::error('修改价格失败' . $exception->getMessage() . ' total: ' . count($this->priceList), $params);
            }

            // 根据返回值进行失败数或成功数 +1
            if ($res) {
                $this->goodsService->incrSucc($this->enterpriseId);
            } else {
                $this->goodsService->incrFail($this->enterpriseId);

                // 插入失败任务列表
                $this->goodsService->pushFailList($this->enterpriseId, $sku, $price, $errorMsg);
            }
        }

        // 当所有改价都完成，解除批量改价的锁定
        $this->goodsService->unTaskLock($this->enterpriseId);
        return true;
    }
}
