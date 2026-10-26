<?php


namespace App\Http\Controllers\Service\V1\Contact;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\ChannelService;
use Illuminate\Support\Facades\DB;


class StatisticsDetailController extends BaseController
{

    /**
     * 联系人二维码明细统计
     * @param ChannelService $channelService
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     * @author: DLong
     */
    public function index(ChannelService $channelService)
    {
        $this->validate($this->request, [
            'id'            => 'required|integer',
            'start_time'    => 'integer',
            'end_time'      => 'integer',
            'type'          => 'in:1,2,3',
            'enterprise_id' => 'required|integer',
        ], [
            'id.required'            => '[参数]联系人二维码id必填',
            'id.integer'             => '[参数]联系人二维码id需为integer',
            'start_time.integer'     => '[参数]start_time需为integer',
            'end_time.integer'       => '[参数]end_time需为integer',
            'type.in'                => '[参数]type值不合法',
            'enterprise_id.required' => '[参数]enterprise_id必填',
            'enterprise_id.integer'  => '[参数]enterprise_id需为integer',
        ]);

        $input = $this->request->only(['id', 'start_time', 'end_time', 'type', 'enterprise_id']);

        // 获取查询的起止时间
        [$start, $end] = $channelService->getTimeRang($input);
        // 查询数据
        $list  = DB::select("SELECT t.date,
count(*) AS count FROM (SELECT FROM_UNIXTIME(created_at,'%Y-%m-%d') AS date
FROM member_user_relation WHERE enterprise_id= " . $input['enterprise_id'] . " AND deleted_at = 0 AND
 source_type=" . $input['id'] . " AND created_at BETWEEN " . $start . " AND " . $end . ") t GROUP BY t.date ORDER BY t.date");
        $chart = [
            'key'   => array_column($list, 'date'),
            'value' => array_column($list, 'count')
        ];
        // 填充没有数据的日期
        $dateArr     = $channelService->getDateRange($start, $end);
        $newChart    = array_combine($chart['key'], $chart['value']);
        $formatChart = [];
        foreach ($dateArr as $item) {
            if (!in_array($item, $chart['key'])) {
                $formatChart['key'][]   = $item;
                $formatChart['value'][] = 0;

                $list[] = [
                    'date'  => $item,
                    'count' => 0,
                ];
            } else {
                $formatChart['key'][]   = $item;
                $formatChart['value'][] = $newChart[$item];
            }
        }

        $dateKeys = array_column($list, 'date');
        array_multisort($dateKeys, SORT_DESC, $list);

        return success([
            'data' => [
                'chart' => $formatChart,
                'table' => $list
            ]
        ]);
    }

}