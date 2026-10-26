<?php


namespace App\Http\Controllers\V1\Message;


use App\Common\RedisKey;
use App\Http\Controllers\V1\BaseController;
use App\Models\EnterpriseActivityMessage;
use App\Models\EnterpriseMessageScene;
use Illuminate\Validation\Rule;


class ActivityMessageController extends BaseController
{
    public function lists()
    {
        $enterpriseId = $this->request->input('enterprise_id');

        if (empty($enterpriseId)) {
            error(700020);
        }

        $page = $this->request->input('page', 1);
        $size = $this->request->input('size', 8);
        $list = EnterpriseActivityMessage::query()
            ->where('enterprise_id', $enterpriseId)
            ->orderByDesc('send_time')
            ->paginate($size, ['id', 'send_time', 'title', 'page_type', 'activity_start_time', 'real_num', 'expected_num', 'send_status'], '', $page)
            ->toArray();

        $result = [
            'data' => $list['data'],
            'meta' => [
                'page'  => $page,
                'total' => $list['total'],
                'size'  => $size
            ],
        ];
        return success($result);

    }

    public function add()
    {
        $params = $this->request->only([
            'enterprise_id',
            'message_scene_id',
            'send_time',
            'title',
            'page_type',
            'activity_start_time',
            'expected_num',
            'page_query',
            'activity_condition',
            'page_url'
        ]);

        $this->validate($this->request, [
            'enterprise_id'       => 'required',
            'message_scene_id'    => 'required',
            'send_time'           => 'required',
            'title'               => 'required | string | between:0,20',
            'page_type'           => ['required', Rule::in([EnterpriseActivityMessage::PAGE_TYPE_GOODS, EnterpriseActivityMessage::PAGE_TYPE_ACTIVITY_PAGE, EnterpriseActivityMessage::PAGE_TYPE_LIVE_LIST, EnterpriseActivityMessage::PAGE_TYPE_LIVE_DETAIL, EnterpriseActivityMessage::PAGE_TYPE_ACTIVITY, EnterpriseActivityMessage::PAGE_TYPE_URL])],
            'activity_start_time' => 'required',
            'expected_num'        => 'required',
        ], [
            'enterprise_id'       => '企业id有误',
            'message_scene_id'    => '订阅消息策略id',
            'send_time'           => '发送时间有误',
            'title.*'             => '活动名称有误',
            'page_type.*'         => '跳转类型有误',
            'activity_start_time' => '活动开始时间有误',
            'expected_num'        => '预计发送数量有误',
        ]);

        // 验证预计发送数量
        $info = EnterpriseMessageScene::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->where('message_scene_id', $params['message_scene_id'])
            ->with('activityMessage')
            ->first();
        // 预计发送数量 不能大于 订阅数 - 发送数 - 占用数
        $diffNum = $info->subscribe_num - $info->send_num - $info->activityMessage
                ->where('send_status', EnterpriseActivityMessage::SEND_STATUS_NO)
                ->where('enable_status', EnterpriseActivityMessage::ENABLE_STATUS_ON)
                ->where('enterprise_id', $params['enterprise_id'])->sum('expected_num');
        if ($diffNum < $params['expected_num']) {
            notice(703002);
        }

        if ($params['send_time'] <= time()) {
            notice(703003);
        }

        $pageUrl = EnterpriseActivityMessage::PAGE_URL[$params['page_type']];
        // 活动 做一口价区分
        switch ($params['page_type']) {
            case EnterpriseActivityMessage::PAGE_TYPE_ACTIVITY:
                if (isset($params['activity_condition']) && $params['activity_condition'] == 1) {
                    $pageUrl = $pageUrl['one-price'];
                } else {
                    $pageUrl = $pageUrl['default'];
                }
                break;
            case EnterpriseActivityMessage::PAGE_TYPE_URL: // 自定义
                if (empty($params['page_url'])) {
                    notice(703004);
                }
                $pageUrl['url'] = $params['page_url'];
                break;
        }
        $pageQuery = $pageUrl['query'];
        foreach ($pageUrl['replace'] as $item) {
            if ($item == 'sku_id' && !isset($params['page_query'][$item])) {
                $params['page_query'][$item] = '';
            }
            if (!isset($params['page_query'][$item])) {
                notice(703007, ['message'], ['message' => 'page_query中缺少'. $item]);
            }
            $requestQuery = $params['page_query'][$item];
            // 活动处理spu_ids
            if ($params['page_type'] == EnterpriseActivityMessage::PAGE_TYPE_ACTIVITY && $item == 'spu_ids') {
                $requestQuery = $params['page_query'][$item][0];
            }
            $pageQuery = str_replace('%' . $item, $requestQuery, $pageQuery);
        }
        $data = [
            'page_url'            => $pageUrl['url'],
            'page_query'          => $pageQuery,
            'page_type'           => $params['page_type'],
            'message_scene_id'    => $params['message_scene_id'],
            'send_time'           => $params['send_time'],
            'title'               => $params['title'],
            'activity_start_time' => $params['activity_start_time'],
            'expected_num'        => $params['expected_num'],
            'enterprise_id'       => $params['enterprise_id']
        ];
        try {
            EnterpriseActivityMessage::query()->create($data);
        } catch (\Exception $e) {
            notice(700000, [], ['message' => $e->getMessage()]);
        }
        return success($data);

    }

    public function delete()
    {
        $params = $this->request->only([
            'enterprise_id',
            'send_id',
        ]);

        $this->validate($this->request, [
            'enterprise_id' => 'required',
            'send_id'       => 'required',
        ], [
            'enterprise_id' => '企业id有误',
            'send_id'       => '发送id有误',
        ]);

        $info = EnterpriseActivityMessage::query()->where('id', $params['send_id'])->first();
        if (empty($info)) {
            notice(703005);
        }
        // 不是未发送不可删除
        if ($info->send_status != 0) {
            notice(703006);
        }
        $info->delete();
        return success();

    }
}