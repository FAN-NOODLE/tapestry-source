<?php


namespace App\Jobs\Export;

use App\Exceptions\Error;
use CrsSdk\Shop\Order;
use Illuminate\Support\Facades\Log;

class OrderStatisticalExportJob   extends ExportJob
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
            '推广订单数',
            '推广订单实收金额（元）'
        ];
        if($this->version!=self::CENTER_VERSION){
            $header = array_merge($header,[
                '服务订单数',
                '服务订单实收金额（元）',
            ]);
        }

        // 查询列表
        try {
            $this->excelExportInit($header)->listExportDate($params);
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
            '推广订单数',
            '推广订单实收金额（元）',
            '服务订单数',
            '服务订单实收金额（元）',
        ];
        if($this->version==self::CENTER_VERSION){
            $header = [
                '导购ID',
                '导购',
                '导购所在部门',
                '推广订单数',
                '推广订单实收金额（元）'
            ];
        }
        // 查询列表
        try {
            $this->excelExportInit($header)->listExportGuide($params);
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
     * 按时间循环获取数据并插入excel
     * @param $params
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Common\Exception\InvalidArgumentException
     * @throws \Box\Spout\Writer\Exception\WriterNotOpenedException
     */
    private function listExportDate($params)
    {
        // 查询获取按时间统计订单列表
        try {
            $result = (new Order())->wisdomOrderStatistical($params);
        } catch (\Exception $exception) {
            Log::info('获取按时间统计订单列表失败' . $exception->getMessage(), $params);
            throw new Error(10000, '获取按时间统计订单列表失败');
        }
        // 組裝数据
        $data = [];
        if (isset($result['detail_data']) && !empty($result['detail_data'])) {
            // 数据排序
            if (!empty($param['sort_type']) && !empty($param['statistical_sort'])) {
                switch ($param['sort_type']) {
                    case 'extension_num':
                        $sortData = array_column($result['detail_data'], 'extension_num');
                        break;
                    case 'extension_amount':
                        $sortData = array_column($result['detail_data'], 'extension_amount');
                        break;
                    case 'service_num':
                        $sortData = array_column($result['detail_data'], 'service_num');
                        break;
                    case 'service_amount':
                        $sortData = array_column($result['detail_data'], 'service_amount');
                        break;
                }
                if ($param['statistical_sort'] == 'desc') {
                    array_multisort($sortData, SORT_DESC, $result['detail_data']);
                } else {
                    array_multisort($sortData, SORT_ASC, $result['detail_data']);
                }
            }
            foreach ($result['detail_data'] as $key => $item) {
                if($this->version!=self::CENTER_VERSION){
                    $data[] = [
                        'data'             => $key,
                        'extension_num'    => "\t" . $item['extension_num'],
                        'extension_amount' => "\t" . intval($item['extension_amount']) / 100,
                        'service_num'      => "\t" . $item['service_num'],
                        'service_amount'   => "\t" . intval($item['service_amount']) / 100,
                    ];
                }else{
                    $data[] = [
                        'data'             => $key,
                        'extension_num'    => "\t" . $item['extension_num'],
                        'extension_amount' => "\t" . intval($item['extension_amount']) / 100
                    ];
                }
            }
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
            $result = (new Order())->wisdomOrderStatisticalByGuide($params);
        } catch (\Exception $exception) {
            Log::info('获取按时间统计订单列表失败' . $exception->getMessage(), $params);
            throw new Error(10000, '获取按时间统计订单列表失败');
        }
        // 組裝数据
        $data = [];
        if (isset($result['list']) && !empty($result['list'])) {
            foreach ($result['list'] as $key => $item) {
                if($this->version==self::CENTER_VERSION){
                    $data[] = [
                        'id'               => "\t" . $item['worker_no'] ?? '',
                        'name'             => $item['worker_name'],
                        'dp_name'          => $item['dp_name'],
                        'extension_num'    => "\t" . $item['extension_num'],
                        'extension_amount' => "\t" . intval($item['extension_amount']) / 100,
                    ];
                }else{
                    $data[] = [
                        'id'               => "\t" . $item['worker_no'] ?? '',
                        'name'             => $item['worker_name'],
                        'store_name'       => $item['store_name'],
                        'dp_name'          => $item['dp_name'],
                        'extension_num'    => "\t" . $item['extension_num'],
                        'extension_amount' => "\t" . intval($item['extension_amount']) / 100,
                        'service_num'      => "\t" . $item['service_num'],
                        'service_amount'   => "\t" . intval($item['service_amount']) / 100,
                    ];
                }

            }
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