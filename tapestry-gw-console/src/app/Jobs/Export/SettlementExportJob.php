<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/10/16 0016
 * Time: 17:28
 */

namespace App\Jobs\Export;


use App\Common\Constant;
use CrsSdk\Work\DistributionSettlement;
use Illuminate\Support\Facades\Log;

class SettlementExportJob extends ExportJob
{

    private $settlementListService;

    protected $header
        = [
            '结算单号',
            '结算状态',
            '订单编号',
            '订单状态',
            '订单支付时间',
            '推广导购',
            '推广导购佣金',
            '推广导购所属门店',
            '服务导购',
            '服务导购佣金',
            '服务导购所属门店',
        ];

    public function __construct(array $params, array $paramsJson, string $name = '结算列表')
    {
        parent::__construct($params, $paramsJson, $name);
        $this->settlementListService = new DistributionSettlement();
    }

    /**
     * @return bool
     */
    public function handle()
    {
        $params = $this->params;
        $header = array_values($this->header);
        // 查询顾客列表

        try {
            $this->excelExportInit($header)->settlementList($params);
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
     * @param array $params
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Common\Exception\InvalidArgumentException
     * @throws \Box\Spout\Writer\Exception\WriterNotOpenedException
     */
    private function settlementList(array $params)
    {

        $settlementGen = $this->yieldSettlement($params);

        foreach ($settlementGen as $settlements) {

            $excelData = [];
            foreach ($settlements['list'] as $settlement) {

                $aSettlement = $this->getSettlement($settlement);
                $excelData[] = $aSettlement->formSingle();
            }
            $this->exportExcel($excelData);
        }

        $this->exportClose();
    }


    private function getSettlement(array $settlement)
    {
        return new class ($settlement)
        {
            public $settlement;

            public function __construct($settlement)
            {
                $this->settlement = $settlement;
            }

            /**
             * excel数据
             *
             * @return array
             */
            public function formSingle()
            {
                return [
                    "\t" . (string)$this->settlement['id'],
                    $this->getSettlementStatus(),
                    $this->settlement['order_no'],
                    $this->getOrderStatus(),
                    date("Y-m-d H:i:s", $this->settlement['order_pay_time']),
                    $this->settlement['promotion_name'],
                    $this->settlement['promotion_commission'] / 100,
                    $this->settlement['promotion_dept'],
                    $this->settlement['service_name'],
                    $this->settlement['service_shopping_guid_id'] ? $this->settlement['service_commission'] / 100 : "",
                    $this->settlement['service_dept']
                ];
            }

            private function getSettlementStatus()
            {
                /**
                 * 结算状态：0待结算1已结算2已取消
                 */
                switch ($this->settlement['settlement_status']) {
                    case 0:
                        return "待结算";
                    case 1:
                        return "已结算";
                    case 2 :
                        return "已取消";
                    default :
                        return "未知结算状态";
                }
            }

            private function getOrderStatus()
            {
                $orderStatusName = [
                    '20'  => '未支付',
                    '21'  => '已取消',
                    '23'  => '待支付',
                    '25'  => '支付中',
                    '30'  => '待确认',
                    '40'  => '待发货',
                    '41'  => '售后中',
                    '60'  => '待收货',
                    '61'  => '售后中',
                    '80'  => '已签收',
                    '81'  => '售后中',
                    '83'  => '部分退款',
                    '100' => '已完成',
                    '200' => '退款中',
                    '201' => '退款失败',
                    '205' => '已退款',
                ];
                return $orderStatusName[$this->settlement['order_status']] ?? "未知订单状态";
            }

        };
    }

    /**
     * @param array $params
     * @return \Generator
     * @throws \Exception
     */
    private function yieldSettlement(array $params)
    {
        $params['page'] = Constant::PAGE;
        $params['size'] = Constant::BATCH_UPDATE_LIMIT;

        while (true) {
            $data = $this->settlementListService->list($params);

            if (empty($data['list'])) {
                break;
            }
            yield $data;
            $params['page'] += 1;
        }
    }

}