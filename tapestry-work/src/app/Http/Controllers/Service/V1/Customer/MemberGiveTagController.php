<?php

namespace App\Http\Controllers\Service\V1\Customer;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Cmq\CmqQueue;
use Illuminate\Support\Facades\App;
use Illuminate\Support\Facades\Log;

class MemberGiveTagController extends BaseController
{
    const QUEUE_NAME = '-user-package';

    /**
     * 给客户打标签
     */
    public function memberGiveTag()
    {
        $env = App::environment();
        if ($env == 'dev' || $env == 'test') {
            $env = 'dev';
        }

        $params         = $this->request->all();
        $data           = [];
        $data['source'] = 'crs';
        $data['action'] = 'set_customer_tags';

        // 企业标识
        $data['query']['enterprise_id'] = $params['enterprise_id'];

        // 单独打
        if (!empty($params['user_id'])) {
            $data['query']['customer_ids'] = $params['user_id'];
        }

        // 要打的标签
        $data['query']['set_tags'] = $params['tag_ids'];

        // 标签条件
        if (!empty($params['tag_id'])) {
            $data['query']['tag_conditions']['tag_ids'][] = $params['tag_id'];
        }

        // 渠道条件
        if (!empty($params['channel_id'])) {
            $data['query']['tag_conditions']['channel_ids'][] = $params['channel_id'];
        }

        // 性别
        if (!empty($params['gender'])) {
            $data['query']['tag_conditions']['gender'] = $params['gender'];
        }

        // 分组条件
        if (!empty($params['group_id'])) {
            $data['query']['tag_conditions']['group_ids'][] = $params['group_id'];
        }

        // 员工
        if (!empty($params['mem_uid'])) {
            $data['query']['tag_conditions']['member_ids'][] = $params['mem_uid'];
        }

        // 部门
        if (!empty($params['dep_uid'])) {
            $data['query']['tag_conditions']['department_ids'][] = $params['dep_uid'];
        }

        // 开始时间
        if (!empty($params['start_time'])) {
            $data['query']['tag_conditions']['start_time'] = $params['start_time'];
        }

        // 结束时间
        if (!empty($params['end_time'])) {
            $data['query']['tag_conditions']['end_time'] = $params['end_time'];
        }

        // 是否被客户删除
        if (!empty($params['type']) && in_array($params['type'], [1, 2])) {
            $data['query']['tag_conditions']['is_deleted'] = $params['type'];
        }

        // 默认不覆盖
        $data['query']['override'] = false;
        if (!empty($params['override'])) {
            $data['query']['override'] = true;
        }

        $queue = (new CmqQueue($env . self::QUEUE_NAME))->sendMessage(json_encode($data));

        if ($env == 'dev' || $env == 'test') {
            Log::info('发送队列打标签', [$queue]);
        }

        return success(['data' => []]);
    }
}