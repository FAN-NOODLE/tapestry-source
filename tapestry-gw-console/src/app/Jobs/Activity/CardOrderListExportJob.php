<?php


namespace App\Jobs\Activity;


use App\Exceptions\Error;
use CrsSdk\MicroActivity\DiscountCard;
use Illuminate\Support\Facades\Log;

class CardOrderListExportJob extends \App\Jobs\Export\ExportJob
{
    /**
     * MaterialExportJob constructor.
     * @param array $params 用于处理业务得参数
     * @param array $paramJson 前端原始请求参数-传数组-下载中心转json
     * @param string $name 任务名称或导出文件sheet名称
     */
    public function __construct(array $params, array $paramJson = [], string $name = '先享卡-活动订单')
    {
        parent::__construct($params, $paramJson, $name);
    }

    /**
     * 导出Excel列表handle
     * @return bool
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Common\Exception\InvalidArgumentException
     * @throws \Box\Spout\Writer\Exception\InvalidSheetNameException
     * @throws \Box\Spout\Writer\Exception\WriterNotOpenedException
     */
    public function handle()
    {
        $params = $this->params;

        $header = ['先享卡ID', '先享卡模板ID', '用户ID', '用户昵称', '订单号', '收货人', '收货人手机号',
            '订单支付时间', '订单状态', '订单是否计入结算', '同步优惠金额', '扣款金额'];

        // 查询获取物料库存管理列表
        try {
            $this->excelExportInit($header)->exportData();
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

    private function exportData()
    {
        // 查询获取物料库存管理列表
        try {
            $openCardList = (new DiscountCard())->cardOrderList($this->params);
        } catch (\Exception $exception) {
            Log::info('获取活动订单列表失败' . $exception->getMessage(), $this->params);
            throw new Error(10000, '获取活动订单列表失败');
        }

        $data = [];
        if (isset($openCardList['list']) && !empty($openCardList['list'])) {
            foreach ($openCardList['list'] as $item) {
                $data[] = [
                    'out_order_no'     => $item['out_order_no'],
                    'discount_card_id' => $item['discount_card_id'],
                    'user_id'          => $item['user_id'],
                    'nick_name'        => $item['nick_name'],
                    'order_no'         => $item['order_no'],
                    'delivery_name'    => $item['delivery_name'],
                    'delivery_mobile'  => $item['delivery_mobile'],
                    'pay_at'           => $item['pay_at'] ?? '--',
                    'order_status'     => $item['order_status'],
                    'is_settle'        => $item['is_settle'],
                    'discount_amount'  => $item['discount_amount'],
                    'deduction_amount' => $item['deduction_amount'],
                ];

            }
        }

        // 插入excel
        $this->exportExcel($data)->exportClose();
    }
}