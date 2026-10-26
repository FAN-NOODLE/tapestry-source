<?php


namespace App\Jobs\Export;

use App\Exceptions\Error;
use CrsSdk\Shop\Order;
use CrsSdk\Work\DistributionStatistics;
use Illuminate\Support\Facades\Log;

class DistributionStatisticalExportJob extends ExportJob
{

    protected $version;
    // 中心商城无门店版
    const  CENTER_VERSION = 3;
    /**
     * GoodsExportJob constructor.
     *
     * @param array  $params 用于处理业务得参数
     * @param array  $paramJson 前端原始请求参数-传数组-下载中心转json
     * @param string $name 任务名称或导出文件sheet名称
     *
     * @throws \Exception
     */
    public function __construct(array $params, array $paramJson = [], string $name = '分销统计列表导出')
    {
        parent::__construct($params, $paramJson, $name);
        $this->version = isset($params['version']) ? $params['version'] : '';

    }

    /**
     * 导出Excel列表handle
     *
     * @return bool
     */
    public function handle()
    {
        $params = $this->params;
        // 判断导出类型
        if ($params['export_type'] == 1) {
            // 按时间导出
            $this->exportByDate($params);
        } else {
            // 按导购导出
            $this->exportByGuide($params);
        }
    }

    /**
     * 根据日期导出
     *
     * @param $params
     *
     * @return bool
     * @author  baizhongyang <baizhongyang@vchangyi.com>
     * @date    2020/8/21 19:34
     */
    protected function exportByDate($params)
    {
        $header = [
            '日期',
            '分销结算笔数',
            '分销结算金额（元）'
        ];

        // 查询列表
        try {
            $this->excelExportInit($header)->listExportDate($params);
        } catch (\Exception $exception) {
            Log::info('分销统计按时间导出失败' . $exception->getMessage(), $params);
            // 修改状态为失败
            $this->jobStatus(false, $exception->getMessage());
            return true;
        }

        // 导出完成 修改状态为成功
        $this->jobStatus();
        return true;
    }

    /**
     * 根据导购导出
     *
     * @param $params
     *
     * @return bool
     * @author  baizhongyang <baizhongyang@vchangyi.com>
     * @date    2020/8/21 20:36
     */
    protected function exportByGuide($params)
    {
        $header = [
            '导购ID',
            '导购',
            '导购所在门店',
            '导购所在部门',
            '分销结算笔数',
            '分销结算金额（元）',
        ];

        if ($this->version == self::CENTER_VERSION) {
            $header = [
                '导购ID',
                '导购',
                '导购所在部门',
                '分销结算笔数',
                '分销结算金额（元）',
            ];
        }

        // 查询列表
        try {
            $this->excelExportInit($header)->listExportGuide($params);
        } catch (\Exception $exception) {
            Log::info('分销统计按导购导出失败' . $exception->getMessage(), $params);
            // 修改状态为失败
            $this->jobStatus(false, $exception->getMessage());
            return true;
        }

        // 导出完成 修改状态为成功
        $this->jobStatus();
        return true;
    }

    /**
     * 按时间循环获取数据并插入excel
     *
     * @param $params
     *
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Common\Exception\InvalidArgumentException
     * @throws \Box\Spout\Writer\Exception\WriterNotOpenedException
     */
    private function listExportDate($params)
    {
        // 查询获取按时间统计分销统计
        try {
            $result = (new DistributionStatistics())->customerStatisticDate($params);
        } catch (\Exception $exception) {
            Log::info('获取按时间统计分销列表失败' . $exception->getMessage(), $params);
            throw new Error(10000, '获取按时间统计分销列表失败');
        }
        Log::info('获取按时间统计分销列表返回结果' . json_encode($result), $params);
        // 組裝数据
        $data = [];
        foreach ($result['list'] as  $item) {

            $data[] = [
                'data'             => "\t" .$item['time'],
                'settlement_num'   => "\t" . $item['settlement_num'],
                'settlement_money' => "\t" . intval($item['settlement_money']) / 100
            ];
        }

        // 插入excel
        $this->exportExcel($data);

        $this->exportClose();
    }

    /**
     * 按导购循环获取数据并插入excel
     *
     * @param $params
     *
     * @author  baizhongyang <baizhongyang@vchangyi.com>
     * @date    2020/8/21 20:48
     */
    private function listExportGuide($params)
    {
        // 查询获取按导购统计订单列表
        try {
            $result = (new DistributionStatistics())->customerStatisticGuide($params);
        } catch (\Exception $exception) {
            Log::info('获取按导购统计分销列表失败' . $exception->getMessage(), $params);
            throw new Error(10000, '获取按导购统计分销列表失败');
        }
        Log::info('获取按导购统计分销列表返回结果' . json_encode($result), $params);

        // 組裝数据
        $data = [];
        foreach ($result['list'] as $item) {
            // 无门店版本
            if ($this->version == self::CENTER_VERSION) {
                $versionParams = [
                    'id'               => "\t" . $item['worker_no'] ?? '',
                    'name'             => $item['worker_name'],
                    'dp_name'          => $item['dp_name'],
                    'settlement_num'   => "\t" . $item['settlement_num'],
                    'settlement_money' => "\t" . intval($item['settlement_money']) / 100,
                ];
            } else {
                $versionParams = [
                    'id'               => "\t" . $item['worker_no'] ?? '',
                    'name'             => $item['worker_name'],
                    'store_name'       => $item['store_name'],
                    'dp_name'          => $item['dp_name'],
                    'settlement_num'   => "\t" . $item['settlement_num'],
                    'settlement_money' => "\t" . intval($item['settlement_money']) / 100,
                ];
            }

            $data[] = $versionParams;
        }

        // 插入excel
        $this->exportExcel($data);

        // 判断是否取完数据
        if (ceil($result['total'] / $params['size']) > $params['page']) {
            $params['page']++;
            $this->listExportGuide($params);
        }
        $this->exportClose();
    }

}