<?php


namespace App\Jobs\Export;


use App\Exceptions\Error;
use CrsSdk\Work\StaffShare;
use Illuminate\Support\Facades\Log;

class ExportBindingCustomerDetailJob extends ExportJob
{

    public $queue = 'export-binding-customer-detail';

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
    public function __construct(array $params, array $paramJson = [], string $name = '导购绑客明细列表')
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
        Log::info("ExportBindingCustomerDetailJob consume", $this->params);
        $_SERVER['enterprise_hash'] = $this->params['enterprise_hash'];
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
            'openid',
            '顾客昵称',
            '绑定导购',
            '绑定时间',
            '是否企微好友',
            '是否商城授权'
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
            $result = (new StaffShare())->bindingCustomerDetail($params);
        } catch (\Exception $exception) {
            Log::info('获取绑客明细列表失败' , [
                'code' => $exception->getCode(),
                'statusCode' => $exception instanceof Error ? $exception->getStatusCode() : 0,
                'message' => $exception->getMessage(),
                'exception' => $exception->getTraceAsString()
            ]);
            throw new Error(10000, '获取绑客明细列表失败');
        }
        // 組裝数据
        $cellData = [];
        if (isset($result['data']) && !empty($result['data'])) {
            foreach ($result['data'] as $key => $value) {
                $cellInfo = [
                    $value['open_id'] ?? '-',
                    $value['nick_name'] ?? '-',
                    ($value['name'] ?? '-').'('.($value['worker_no'] ?? '-').')',
                    $value['binding_time'] ?? '-',
                    !empty($value['whether']) ? ($value['whether'] == 1 ? '是' : '否') : '否',
                    !empty($value['union_id']) ? '是' : '否',
                ];
                $cellData[] = $cellInfo;
            }
            $this->exportExcel($cellData);

        }
        // 判断是否取完数据
        if ($params['page_total'] > $params['page']) {
            $params['page']++;
            $this->listExport($params);
        }
        $this->exportClose();
    }

}