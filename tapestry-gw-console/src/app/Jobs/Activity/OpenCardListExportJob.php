<?php


namespace App\Jobs\Activity;


use App\Exceptions\Error;
use CrsSdk\MicroActivity\DiscountCard;
use Illuminate\Support\Facades\Log;

class OpenCardListExportJob extends \App\Jobs\Export\ExportJob
{
    /**
     * MaterialExportJob constructor.
     * @param array $params 用于处理业务得参数
     * @param array $paramJson 前端原始请求参数-传数组-下载中心转json
     * @param string $name 任务名称或导出文件sheet名称
     */
    public function __construct(array $params, array $paramJson = [], string $name = '开卡信息数据')
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

        $header = ['先享卡ID', '先享卡模板ID', '用户ID', '用户昵称', '分享者', '用户手机号', '开卡时间',
            '用户履约状态', '履约进度', '履约时间', '是否发生扣款', '总优惠金额', '总扣款金额'];

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
            $openCardList = (new DiscountCard())->openCardList($this->params);
        } catch (\Exception $exception) {
            Log::info('获取开卡信息列表失败' . $exception->getMessage(), $this->params);
            throw new Error(10000, '获取开卡信息列表失败');
        }

        $data = [];
        if (isset($openCardList['list']) && !empty($openCardList['list'])) {
            foreach ($openCardList['list'] as $item) {
                $data[] = [
                    'out_order_no'     => $item['out_order_no'],
                    'discount_card_id' => $item['discount_card_id'],
                    'user_id'          => $item['user_id'],
                    'nick_name'        => $item['nick_name'],
                    'share_name'       => $item['share_name'],
                    'mobile'           => $item['mobile'],
                    'open_at'          => empty($item['open_at']) ? '' : date('Y-m-d H:i:s', $item['open_at']),
                    'status'           => $item['status'],
                    'process'          => $item['process'],
                    'process_at'       => empty($item['process_at']) ? '' : date('Y-m-d H:i:s', $item['process_at']),
                    'deduction_status' => $item['deduction_status'],
                    'discount_all'     => $item['discount_all'],
                    'deduction_all'    => $item['deduction_all'],
                ];

            }
        }

        // 插入excel
        $this->exportExcel($data)->exportClose();
    }


}