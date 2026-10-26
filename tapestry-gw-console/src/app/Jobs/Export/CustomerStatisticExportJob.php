<?php


namespace App\Jobs\Export;

use App\Exceptions\Error;
use CrsSdk\Shop\Order;
use CrsSdk\Work\DataCompass;
use CrsSdk\Work\Member;
use Illuminate\Support\Facades\Log;

class CustomerStatisticExportJob extends ExportJob
{

    protected $version;
    // 中心商城无门店版
    const  CENTER_VERSION = 3;
    /**
     * @param array $params 用于处理业务得参数
     * @param array $paramJson 前端原始请求参数-传数组-下载中心转json
     * @param string $name 任务名称或导出文件sheet名称
     * @throws \Exception
     */
    public function __construct(array $params, array $paramJson = [], string $name = '顾客数据统计列表导出')
    {
        parent::__construct($params, $paramJson, $name);
        $this->version = $params['version'];

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
     * exportByDate
     * @desc 根据日期导出
     * @param $params
     *
     * @return bool
     * @date 2020-08-23
     */
    protected function exportByDate($params)
    {
        $header = [
            '日期',
            '顾客总数',
            '导购名片添加总数',
            '其他方式添加总数',
            '新增顾客数',
            '导购名片新加数',
            '其他方式新加数',
            '发起会话数',
            '咨询应答率',
            '绑定顾客总数',
            '新增绑客数',
        ];

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
     * exportByGuide
     * @desc
     * @param $params
     *
     * @return bool
     * @date 2020-08-23
     */
    protected function exportByGuide($params)
    {
        if($this->version==self::CENTER_VERSION){
            $header = [
                '导购ID',
                '导购',
                '导购所在部门',
                '顾客总数',
                '导购名片添加总数',
                '其他方式添加总数',
                '新增顾客数',
                '导购名片新添加数',
                '其他方式新添加数',
                '发起会话数',
                '咨询应答率',
                '绑定顾客总数',
                '新增绑客数',
            ];
        }else{
            $header = [
                '导购ID',
                '导购',
                '导购所在门店',
                '导购所在部门',
                '顾客总数',
                '导购名片添加总数',
                '其他方式添加总数',
                '新增顾客数',
                '导购名片新添加数',
                '其他方式新添加数',
                '发起会话数',
                '咨询应答率',
                '绑定顾客总数',
                '新增绑客数',
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
            $result = (new DataCompass())->customerStatisticListByDate($params);
        } catch (\Exception $exception) {
            Log::info('获取按时间统计顾客列表失败' . $exception->getMessage(), $params);
            throw new Error(10000, '获取按时间统计顾客列表失败');
        }
        // 組裝数据
        $data = [];
        if (isset($result['list']) && !empty($result['list'])) {
            foreach ($result['list'] as $key => $item) {
                $data[] = [
                    date('Y-m-d', $item['stat_time']),
                    $item['cumulative_customer_num'] ?? '-',
                    $item['cumulative_contact_card_way_customer_num'] ?? '-',
                    $item['cumulative_other_way_customer_num'] ?? '-',
                    $item['today_customer_num'] ?? '-',
                    $item['today_contact_card_way_customer_num'] ?? '-',
                    $item['today_other_way_customer_num'] ?? '-',
                    $item['chat_cnt'] ?? '-',
                    !empty($item['reply_percentage']) ? ($item['reply_percentage'] * 100) . '%' : '-',
                    $item['cumulative_binding_customer_num'] ?? '-',
                    $item['today_binding_customer_num'] ?? '-',
                ];
            }
        }

        // 插入excel
        $this->exportExcel($data);

        $this->exportClose();
    }

    /**
     * listExportGuide
     * @desc
     * @param $params
     * @date 2020-08-23
     *
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Common\Exception\InvalidArgumentException
     * @throws \Box\Spout\Writer\Exception\WriterNotOpenedException
     */
    private function listExportGuide($params)
    {
        // 查询获取按导购统计订单列表
        try {
            $result = (new DataCompass())->customerStatisticListByWorker($params);
        } catch (\Exception $exception) {
            Log::info('顾客数据列表查询失败' . $exception->getMessage(), $params);
            throw new Error(10000, '顾客数据列表查询失败');
        }
        // 組裝数据
        $data = [];

        if (isset($result['list']) && !empty($result['list'])) {
            foreach ($result['list'] as $key => $item) {
                if ($this->version == self::CENTER_VERSION) {
                    $data[] = [
                        $item['worker_no'] ?? '-',
                        $item['name'] ?? '-',
                        $item['dp_name'] ?? '-',
                        $item['cumulative_customer_num'] ?? '-',
                        $item['cumulative_contact_card_way_customer_num'] ?? '-',
                        $item['cumulative_other_way_customer_num'] ?? '-',
                        $item['today_customer_num'] ?? '-',
                        $item['today_contact_card_way_customer_num'] ?? '-',
                        $item['today_other_way_customer_num'] ?? '-',
                        $item['chat_cnt'] ?? '-',
                        !empty($item['reply_percentage']) ? ($item['reply_percentage'] * 100) . '%' : '-',
                        $item['cumulative_binding_customer_num'] ?? '-',
                        $item['today_binding_customer_num'] ?? '-',
                    ];
                } else {
                    $data[] = [
                        $item['worker_no'] ?? '-',
                        $item['name'] ?? '-',
                        $item['store_name'] ?? '-',
                        $item['dp_name'] ?? '-',
                        $item['cumulative_customer_num'] ?? '-',
                        $item['cumulative_contact_card_way_customer_num'] ?? '-',
                        $item['cumulative_other_way_customer_num'] ?? '-',
                        $item['today_customer_num'] ?? '-',
                        $item['today_contact_card_way_customer_num'] ?? '-',
                        $item['today_other_way_customer_num'] ?? '-',
                        $item['chat_cnt'] ?? '-',
                        !empty($item['reply_percentage']) ? ($item['reply_percentage'] * 100) . '%' : '-',
                        $item['cumulative_binding_customer_num'] ?? '-',
                        $item['today_binding_customer_num'] ?? '-',
                    ];
                }

            }
        }

        // 插入excel
        $this->exportExcel($data);

        $this->exportClose();
    }

}