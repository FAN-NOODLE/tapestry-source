<?php

namespace App\Http\Controllers\Service\V1\Channel;

use App\Common\Constant;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Channel;
use App\Models\ContactQrCode;
use App\Models\MemberUserRelation;
use App\Services\ChannelService;
use Illuminate\Support\Facades\DB;

class StatisticsListController extends BaseController
{
    /**
     * 渠道统计列表
     * {domain}/service/v1/channel/statistics-list
     *
     * @param ChannelService $channelService
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(ChannelService $channelService)
    {
        $this->validate($this->request, [
            'start_time'    => 'integer',
            'end_time'      => 'integer',
            'type'          => 'in:1,2,3',
            'page'          => 'integer',
            'size'          => 'integer',
            'enterprise_id' => 'required|integer',
        ], [
            'start_time.integer'     => '[参数]start_time需为integer',
            'end_time.integer'       => '[参数]end_time需为integer',
            'type.in'                => '[参数]type值不合法',
            'page.integer'           => '[参数]page需为integer',
            'size.integer'           => '[参数]size需为integer',
            'enterprise_id.required' => '[参数]enterprise_id必填',
            'enterprise_id.integer'  => '[参数]enterprise_id需为integer',
        ]);

        $input = $this->request->only(['start_time', 'end_time', 'type', 'page', 'size', 'enterprise_id']);

        $page = intval($input['page'] ?? Constant::PAGE);
        $size = intval($input['size'] ?? Constant::SIZE);
        // 企业id
        $enterpriseId = intval($input['enterprise_id']);

        $data = [];

        // 总数据量
        $total = Channel::query()->where('channel.enterprise_id', $enterpriseId)->count();
        if (!empty($total)) {
            $data = Channel::query()
                ->where('enterprise_id', $enterpriseId)
                ->orderBy('id', 'DESC')
                ->forpage($page, $size)
                ->get(['id', 'channel'])
                ->toArray();

            // 获取查询的起止时间
            [$start, $end] = $channelService->getTimeRang($input);

            foreach ($data as &$value) {
                $totalQrCode = ContactQrCode::query()->where('channel_id', $value['id'])->pluck('id')->toArray();

                //历史总纪录
                $totalRelation         = MemberUserRelation::query()
                    ->whereIn('source_type', $totalQrCode)
                    ->count();
                $totalCancelRelation   = MemberUserRelation::query()
                    ->whereIn('source_type', $totalQrCode)
                    ->where('delete_type', '!=', 0)
                    ->count();
                $value['total']        = $totalRelation;
                $value['cancel_total'] = $totalCancelRelation;
                //指定时间段内
                $totalRelationByDate       = MemberUserRelation::query()
                    ->whereIn('source_type', $totalQrCode)
                    ->whereBetween('relation_at', [$start, $end])
                    ->count();
                $totalCancelRelationByDate = MemberUserRelation::query()
                    ->whereIn('source_type', $totalQrCode)
                    ->where('delete_type', '!=', 0)
                    ->whereBetween('cancel_relation_at', [$start, $end])
                    ->count();
                $value['count']            = $totalRelationByDate;
                $value['cancel_count']     = $totalCancelRelationByDate;
            }
        }

        return success([
            'data' => $data,
            'meta' => [
                'page'  => $page,
                'size'  => $size,
                'total' => $total
            ]
        ]);
    }
}
