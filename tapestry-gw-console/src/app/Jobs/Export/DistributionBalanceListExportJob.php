<?php


namespace App\Jobs\Export;

use App\Exceptions\Error;
use CrsSdk\Work\DistributionBalance;
use Illuminate\Support\Facades\Log;

class DistributionBalanceListExportJob extends ExportJob
{

    protected $version;
    // 中心商城无门店版
    const  CENTER_VERSION = 3;
    /**
     * GoodsExportJob constructor.
     * @param array $params 用于处理业务得参数
     * @param array $paramJson 前端原始请求参数-传数组-下载中心转json
     * @param string $name 任务名称或导出文件sheet名称
     * @throws \Exception
     */
    public function __construct(array $params, array $paramJson = [], string $name = '订单统计列表导出')
    {
        parent::__construct($params, $paramJson, $name);
        $this->version = isset($params['version']) ? $params['version'] : '';

    }

    /**
     * 导出Excel列表handle
     * @return bool
     */
    public function handle()
    {
        $this->export($this->params);

    }

    /**
     * 导出
     *
     * @param $params
     * @return bool
     * User: zhaoyang@vchangyi.com
     * Date: 2020/10/17 10:02
     */
    protected function export($params)
    {
        $header = [
            '导购ID',
            '导购',
            '导购手机号',
            '导购所在' . ($this->version == self::CENTER_VERSION ? '部门' : '门店/部门'),
            '待结算金额（元）',
            '可提现金额（元）',
            '已提现金额（元）',
        ];
        // 查询列表
        try {
            $this->excelExportInit($header)->listExport($params);
        } catch (\Exception $exception) {
            Log::info('导出失败' . $exception->getMessage(), $params);
            // 修改状态为失败
            $this->jobStatus(false, $exception->getMessage());
            return true;
        }

        // 导出完成 修改状态为成功
        $this->jobStatus();
        return true;
    }

    /**
     * 按循环获取数据并插入excel
     *
     * @param $params
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Common\Exception\InvalidArgumentException
     * @throws \Box\Spout\Writer\Exception\WriterNotOpenedException
     * User: zhaoyang@vchangyi.com
     * Date: 2020/10/16 20:13
     */
    private function listExport($params)
    {
        // 查询获取按导购统计订单列表
        try {
            $result = (new DistributionBalance())->balanceList($params);
        } catch (\Exception $exception) {
            Log::info('获取全员分销余额列表失败' . $exception->getMessage(), $params);
            throw new Error(10000, '获取全员分销余额列表失败');
        }
        // 組裝数据
        $cellData = [];
        if (isset($result['list']) && !empty($result['list'])) {
            foreach ($result['list'] as $key => $value) {
                $storeName  = $value['short_name'] ?: $value['store_name'];
                $cellData[] = [
                    "\t" . (string)$value['worker_no'],
                    $value['name'],
                    $value['mobile'],
                    $this->version == self::CENTER_VERSION
                        ? $storeName
                        : ($storeName . '/' . join(',', $value['dpName'])),
                    $value['expect_cash_amount'] / 100,
                    $value['available_cash_amount'] / 100,
                    $value['cash_withdrawal'] / 100,
                ];
            }
            $this->exportExcel($cellData);

        }
        // 判断是否取完数据
        if (ceil($result['total'] / $params['size']) > $params['page']) {
            $params['page']++;
            $this->listExport($params);
        }
        $this->exportClose();
    }

}