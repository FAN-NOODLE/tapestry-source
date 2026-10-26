<?php

namespace App\Http\Controllers\Service\V1\Channel;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\ContactQrCode;
use App\Services\ChannelService;
use Illuminate\Support\Facades\DB;

class StatisticsDetailController extends BaseController
{
    /**
     * 渠道统计明细
     * {domain}/service/v1/channel/statistics-detail
     *
     * @param ChannelService $channelService
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
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
            'id.required'            => '[参数]渠道id必填',
            'id.integer'             => '[参数]渠道id需为integer',
            'start_time.integer'     => '[参数]start_time需为integer',
            'end_time.integer'       => '[参数]end_time需为integer',
            'type.in'                => '[参数]type值不合法',
            'enterprise_id.required' => '[参数]enterprise_id必填',
            'enterprise_id.integer'  => '[参数]enterprise_id需为integer',
        ]);

        $input = $this->request->only(['id', 'start_time', 'end_time', 'type', 'enterprise_id']);

        //当前渠道下的所有二维码id集合
        $qrCodeIds    = ContactQrCode::query()->where('channel_id', $input['id'])
            ->where('enterprise_id', $input['enterprise_id'])
            ->pluck('id')->toArray();
        $qrCodeIdsStr = implode(',', $qrCodeIds);

        // 获取查询的起止时间
        [$start, $end] = $channelService->getTimeRang($input);

        $list = [];
        if ($qrCodeIdsStr) {
            // 查询数据
            $list  = DB::select("SELECT t.date,
(select count(id) from member_user_relation as a where a.deleted_at=0 and a.enterprise_id=" . $input['enterprise_id'] . "
         and a.delete_type<>0 and FROM_UNIXTIME( a.cancel_relation_at, '%Y-%m-%d' )=t.date and a.source_type 
         in(" . $qrCodeIdsStr . ") ) as cancel_count,
count(*) AS count FROM (SELECT FROM_UNIXTIME(relation_at,'%Y-%m-%d') AS date 
FROM member_user_relation WHERE enterprise_id= " . $input['enterprise_id'] . " AND deleted_at = 0 AND source_type
 in(" . $qrCodeIdsStr . ") AND relation_at BETWEEN " . $start . " AND " . $end . ") t GROUP BY t.date ORDER BY t.date");
        }

        $chart = [
            'key'          => array_column($list, 'date'),
            'value'        => array_column($list, 'count'),
            'value_cancel' => array_column($list, 'cancel_count')
        ];
        // 填充没有数据的日期
        $dateArr        = $channelService->getDateRange($start, $end);
        $newChart       = array_combine($chart['key'], $chart['value']);
        $newChartCancel = array_combine($chart['key'], $chart['value_cancel']);
        $formatChart    = [];
        foreach ($dateArr as $item) {
            if (!in_array($item, $chart['key'])) {
                $formatChart['key'][]   = $item;
                $formatChart['value'][] = [0, 0];

                $list[] = [
                    'date'         => $item,
                    'cancel_count' => 0,
                    'count'        => 0,
                ];
            } else {
                $formatChart['key'][]   = $item;
                $formatChart['value'][] = [$newChart[$item], $newChartCancel[$item]];
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
