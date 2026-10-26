<?php


namespace App\Jobs\Activity;


use App\Exceptions\Error;
use App\Micros\Activity\LuckMicro;
use Illuminate\Support\Facades\Log;

class LuckJoinUserExport extends \App\Jobs\Export\ExportJob
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

        $header = ['昵称', '手机号', '参与时间'];

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
            $list = (new LuckMicro())->joinUserList($this->params);
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
                    'created_at'   => date("Y-m-d H:i:s", $item['created_at']),
                ];

            }
        }

        // 插入excel
        $this->exportExcel($data)->exportClose();
    }

}