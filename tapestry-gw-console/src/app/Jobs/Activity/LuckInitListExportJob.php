<?php


namespace App\Jobs\Activity;

use App\Exceptions\Error;
use App\Micros\Activity\LuckMicro;
use Illuminate\Support\Facades\Log;

class LuckInitListExportJob extends \App\Jobs\Export\ExportJob
{
    /**
     * MaterialExportJob constructor.
     * @param array $params 用于处理业务得参数
     * @param array $paramJson 前端原始请求参数-传数组-下载中心转json
     * @param string $name 任务名称或导出文件sheet名称
     */
    public function __construct(array $params, array $paramJson = [], string $name = '好友抽奖发起用户')
    {
        parent::__construct($params, $paramJson, $name);
    }

    public function handle()
    {
        $params = $this->params;

        $header = ['昵称', '手机号', '助力开始时间', '裂变好友数量', '状态'];

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
            $list = (new LuckMicro())->initRankList($this->params);
        } catch (\Exception $exception) {
            Log::info('获取活动订单列表失败' . $exception->getMessage(), $this->params);
            throw new Error(10000, '获取活动订单列表失败');
        }

        $data = [];
        if (isset($list['list']) && !empty($list['list'])) {
            foreach ($list['list'] as $item) {
                $data[] = [
                    'nick_name' => $item['nick_name'],
                    'mobile'    => $item['mobile'],
                    'init_at'   => date("Y-m-d H:i:s", $item['init_at']),
                    'cnt'       => $item['cnt'],
                    'status'    => $this->statusText($item['status']),
                ];

            }
        }

        // 插入excel
        $this->exportExcel($data)->exportClose();
    }

    /**
     * 1：助力中 2：助力成功 3：助力失败 4：未中奖 5：中奖
     * @param int $status
     * @return string
     */
    private function statusText(int $status): string
    {
        $text = "--";
        switch ($status) {
            case 1:
                $text = "助力中";
                break;
            case 2:
                $text = "助力成功";
                break;
            case 3:
                $text = "助力失败";
                break;
            case 4:
                $text = "未中奖";
                break;
            case 5:
                $text = "中奖";
                break;
        }
        return $text;
    }


}