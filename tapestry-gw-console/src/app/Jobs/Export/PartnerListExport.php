<?php


namespace App\Jobs\Export;


use App\Exceptions\Error;
use CrsSdk\Work\Partner;
use Illuminate\Support\Facades\Log;

class PartnerListExport extends ExportJob
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
    public function __construct(array $params, array $paramJson = [], string $name = '提现审批列表')
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
            '合伙人ID',
            '合伙人',
            '手机号',
            '所属导购',
            '待结算佣金（元）',
            '可提现佣金（元）',
            '提现中佣金（元）',
            '已提现佣金（元）',
            '注册时间',
            '状态',
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
            $result = (new Partner())->partnerList($params);
        } catch (\Exception $exception) {
            Log::info('获取合伙人列表失败' . $exception->getMessage(), $params);
            throw new Error(10000, '获取合伙人列表失败');
        }
        // 組裝数据
        $cellData = [];
        if (isset($result['list']) && !empty($result['list'])) {
            foreach ($result['list'] as $key => $value) {
                $cellInfo = [
                    $value['partner_sn'] ?? '',
                    $value['partner_name'],
                    $value['partner_mobile'],
                    $value['worker_name'],
                    $value['for_withdrawal'] / 100,
                    $value['can_withdrawal'] / 100,
                    $value['withdrawal_of'] / 100,
                    $value['have_withdrawal'] / 100,
                    date('Y-m-d H:i:s', $value['created_at']),
                    $value['partner_status'] == 1 ? '启用' : '禁用',
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