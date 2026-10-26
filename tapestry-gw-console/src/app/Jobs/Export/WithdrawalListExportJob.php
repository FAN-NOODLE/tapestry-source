<?php


namespace App\Jobs\Export;

use App\Exceptions\Error;
use CrsSdk\Work\Withdrawal;
use Illuminate\Support\Facades\Log;

class WithdrawalListExportJob extends ExportJob
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
            '提现申请单号',
            '申请金额（元）',
            '导购ID',
            '申请导购',
            '导购手机号',
            '导购所在门店/部门',
            '申请发起时间',
        ];

        switch ($params['status']) {
            case 2:
                array_push($header,'申请处理时间');
                break;
            case 3:
                array_push($header,'申请处理时间');
                array_push($header,'提现失败原因');
                break;
            case 4:
                array_push($header,'申请处理时间');
                array_push($header,'驳回原因');
                break;
            default:
                break;
        }
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
        // 查询获取按导购统计订单列表
        try {
            $result = (new Withdrawal())->ApprovalList($params);
        } catch (\Exception $exception) {
            Log::info('获取提现审批列表失败' . $exception->getMessage(), $params);
            throw new Error(10000, '获取提现审批列表失败');
        }
        // 組裝数据
        $cellData = [];
        if (isset($result['list']) && !empty($result['list'])) {
            foreach ($result['list'] as $key => $value) {
                $cellInfo = [
                    "\t" . (string)$value['id'],
                    $value['application_amount']/100,
                    $value['worker_id'],
                    $value['name'],
                    $value['mobile'],
                    empty($value['dp_store']) ? '-' : $value['dp_store'],
                    date('Y-m-d H:i:s', $value['request_initiated_at']),
                ];
                switch ($params['status']) {
                    case 2:
                        array_push($cellInfo,date('Y-m-d H:i:s', $value['request_handle_at']));
                        break;
                    case 3:
                    case 4:
                        array_push($cellInfo,date('Y-m-d H:i:s', $value['request_handle_at']));
                        array_push($cellInfo,$value['reason']);
                        break;
                    default:
                        break;
                }
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