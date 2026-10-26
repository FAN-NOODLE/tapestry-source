<?php

namespace App\Http\Controllers\Html\V1\Customer;

use App\Http\Controllers\Html\V1\BaseController;
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

        $params         = request()->only(['tag_ids', 'user_id', 'override']);
        $data           = [];
        $data['source'] = 'crs';
        $data['action'] = 'set_customer_tags';

        // 企业标识
        $data['query']['enterprise_id'] = $this->enterpriseId;

        // 单独打
        if (!empty($params['user_id'])) {
            $data['query']['customer_ids'] = [
                [
                    'ec_userid' => $params['user_id'],
                    'mem_uid' => $this->memUid
                ]
            ];
        }

        // 要打的标签
        $data['query']['set_tags'] = $params['tag_ids'];

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