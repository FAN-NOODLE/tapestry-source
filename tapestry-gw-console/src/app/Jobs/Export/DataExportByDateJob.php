<?php


namespace App\Jobs\Export;


use App\Exceptions\Error;
use CrsSdk\Work\DataCompass;
use CrsSdk\Work\Partner;
use Illuminate\Support\Facades\Log;

class DataExportByDateJob extends ExportJob
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
    public function __construct(array $params, array $paramJson = [], string $name = '顾客数据按日期列表')
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
     */
    protected function export($params)
    {
        $header = [
            '日期',
            '顾客总数',
            '导购名片添加总数',
            '其它方式添加总数',
            '新增顾客数',
            '导购名片新加数',
            '其它方式新加数',
            '发起会话数',
            '咨询应答率',
            '绑定顾客总数',
            '新增帮客数',
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
     */
    private function listExport($params)
    {
        try {
            $result = (new DataCompass())->customerStatisticListByDate($params);
        } catch (\Exception $exception) {
            Log::info('获取数据罗盘顾客数据按日期列表失败' . $exception->getMessage(), $params);
            throw new Error(10000, '获取数据罗盘顾客数据按日期列表失败');
        }
        // 組裝数据
        $cellData = [];
        if (isset($result['list']) && !empty($result['list'])) {
            foreach ($result['list'] as $key => $value) {
                $cellInfo = [
                    date('Y-m-d H:i:s', $value['stat_time']),
                    $value['cumulative_customer_num'] ?? '-',
                    $value['cumulative_contact_card_way_customer_num'] ?? '-',
                    $value['cumulative_other_way_customer_num'] ?? '-',
                    $value['today_contact_card_way_customer_num'] ?? '-',
                    $value['today_contact_card_way_customer_num'] ?? '-',
                    $value['today_other_way_customer_num'] ?? '-',
                    $value['chat_cnt'] ?? '-',
                    $value['reply_percentage'] ?? '-',
                    $value['cumulative_binding_customer_num'] ?? '-',
                    $value['today_binding_customer_num'] ?? '-',
                ];
                $cellData[] = $cellInfo;
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