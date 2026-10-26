<?php


namespace App\Jobs\Export;


use App\Exceptions\Error;
use CrsSdk\Shop\Order;
use CrsSdk\Work\PartnerMember;
use CrsSdk\Work\StoreWorker;
use Illuminate\Support\Facades\Log;

class OrderStatisticalDetailExportJob extends ExportJob
{
    protected $version;
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
     */
    protected function export($params)
    {
        $header = [
            '订单编号',
            '订单类型',
            '订单状态',
            '支付时间',
            '所属导购',
            '下单客户',
            '收货人',
            '收货人手机号',
            '实付金额（元）',
            '实收金额（元）',
            '退款总额（元）',
            '退款到账时间',
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
        // 查询合伙人
        $partnerQueryParams = [
            'enterprise_id' => $params['enterprise_id'],
            'page' => 1,
            'size' => 100000
        ];
        if (!empty($param['worker_ids'])) {
            $partnerQueryParams['worker_ids'] = $params['worker_ids'];
        }
        $partnerList = (new PartnerMember())->partnerList($partnerQueryParams)['list'];
        $partnerMap = array_column($partnerList, null, 'id');
        if (!empty($partnerList) && !empty($params['worker_ids'])) {
            $params['worker_ids'] = array_merge($params['worker_ids'], array_column($partnerList, 'id'));
        }
        try {
            $result = (new Order())->exportOrderStatisticList($params);
        } catch (\Exception $exception) {
            Log::info('获取数据罗盘订单明细列表失败' . $exception->getMessage(), $params);
            throw new Error(10000, '获取订单明细列表失败');
        }

        // 组装数据
        $cellData = [];
        if (isset($result['list']) && !empty($result['list'])) {
            $workerList = (new StoreWorker())->getList([
                'enterprise_id' => $params['enterprise_id'],
                'is_page' => false
            ])['list'];
            $workerMap = array_column($workerList, null, 'id');
            foreach ($result['list'] as $key => $value) {
                $orderTypeArr = [];
                $belongGuideArr = [];
                if (!empty($value['store_staff_id'])) {
                    if (!empty($workerMap[$value['store_staff_id']])) {
                        $promoteWorker = $workerMap[$value['store_staff_id']];
                        if (empty($param['worker_ids']) || in_array($promoteWorker['id'], $params['worker_ids'])) {
                            $orderTypeArr[] = '推广订单';
                            $belongGuideArr[] = $promoteWorker['name'] . '(' .$promoteWorker['worker_no'] . ')';
                        }
                    } elseif (!empty($partnerMap[$value['store_staff_id']])) {
                        $promotePartner = $partnerMap[$value['store_staff_id']];
                        $promoteWorker = $workerMap[$promotePartner['worker_id']] ?? [];
                        if (empty($param['worker_ids']) || in_array($promoteWorker['id'], $params['worker_ids'])) {
                            $orderTypeArr[] = '推广订单（二级）';
                            if (!empty($promoteWorker)) {
                                $belongGuideArr[] = $promoteWorker['name'] . '(' .$promoteWorker['worker_no'] . ')';
                            } else {
                                $belongGuideArr[] = '';
                            }
                        }
                    }
                }
                if (!empty($value['service_staff_id'])) {

                    if (!empty($workerMap[$value['service_staff_id']])) {
                        $serviceWorker = $workerMap[$value['service_staff_id']];
                        if (empty($params['worker_ids']) || in_array($serviceWorker['id'], $params['worker_ids'])) {
                            $orderTypeArr[] = '服务订单';
                            $belongGuideArr[] = $serviceWorker['name'] . '(' .$serviceWorker['worker_no'] . ')';
                        }
                    }
                }
                if (count($belongGuideArr) == 2 && $belongGuideArr[0] == $belongGuideArr[1]) {
                    array_push($cellData, [
                        $value['order_no'],
                        implode('、', $orderTypeArr),
                        $value['status_desc'],
                        date('Y-m-d H:i:s', $value['pay_at']),
                        $belongGuideArr[0],
                        $value['nick_name'],
                        $value['delivery_name'],
                        $value['delivery_mobile'],
                        $value['cash_fee'],
                        $value['receipts_fee'],
                        $value['total_refund_fee'],
                        !empty($value['refund_at']) ? date('Y-m-d H:i:s', $value['refund_at']) : '-',
                    ]);
                } else {
                    foreach ($belongGuideArr as $k => $belongGuide) {
                        array_push($cellData, [
                            $value['order_no'],
                            $orderTypeArr[$k],
                            $value['status_desc'],
                            date('Y-m-d H:i:s', $value['pay_at']),
                            $belongGuide,
                            $value['nick_name'],
                            $value['delivery_name'],
                            $value['delivery_mobile'],
                            $value['cash_fee'],
                            $value['receipts_fee'],
                            $value['total_refund_fee'],
                            !empty($value['refund_at']) ? date('Y-m-d H:i:s', $value['refund_at']) : '-',
                        ]);
                    }
                }
            }
        }
        $this->exportExcel($cellData);
        $this->exportClose();
    }
}