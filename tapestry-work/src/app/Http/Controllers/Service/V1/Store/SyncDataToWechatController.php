<?php

namespace App\Http\Controllers\Service\V1\Store;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Jobs\SyncDataToWechatJob;
use App\Models\StoreDepartment;
use Illuminate\Support\Facades\Redis;

class SyncDataToWechatController extends BaseController
{

    /**
     * @return array
     * User: zhaoyang@vchangyi.com
     * Date: 2022/2/7 15:45
     * @throws \Exception
     */
    public function index()
    {
        $params = $this->request->only(['enterprise_id', 'is_sync']);
        if (empty($params['is_sync'])) {
            return success(['last_sync_time' => StoreDepartment::query()->max('updated_at')]);
        }
        if (empty($params['enterprise_id'])) {
            throw new Error(600000, '企业id不能为空');
        } else {
            $key = redisKey('SYNC_STORE_INFO_BY_CHAT_TIME', ['enterprise_id' => $params['enterprise_id']]);
            $syncTime = Redis::get($key);
            if ($syncTime > 0) {
                $time = time() - $syncTime;
                if ($time < 300) {
                    throw new Error(600017, '请' . (300 - $time) . '秒后再试（数据同步最小间隔为5分钟）');
                }
            }

            dispatch(new SyncDataToWechatJob($params['enterprise_id']));
        }

        return success(['last_sync_time' => time()]);
    }
}
