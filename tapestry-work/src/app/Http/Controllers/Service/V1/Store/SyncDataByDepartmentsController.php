<?php

namespace App\Http\Controllers\Service\V1\Store;

use App\Console\Commands\SyncStoreDepartmentByCoach;
use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\StoreDepartment;

class SyncDataByDepartmentsController extends BaseController
{

    /**
     * @return array
     * User: zhaoyang@vchangyi.com
     * Date: 2022/2/7 15:45
     * @throws \Exception
     */
    public function index(): array
    {
        $params = $this->request->only(['enterprise_id', 'is_sync']);
        if (empty($params['enterprise_id'])) {
            throw new Error(600000, '企业id不能为空');
        }
        $key = redisKey('SYNC_STORE_INFO_BY_DEPARTMENT_TIME', ['enterprise_id' => $params['enterprise_id']]);
        $syncTime = getOrSetData($key);
        if (empty($syncTime)) {
            $syncTime = StoreDepartment::query()
                ->where('enterprise_id', $params['enterprise_id'])
                ->max('updated_at');
        }
        if (empty($params['is_sync'])) {
            return success(['last_sync_time' => $syncTime]);
        }

        $diffTime = time() - $syncTime;
        if ($diffTime < 300) {
            throw new Error(600017, '请' . (300 - $diffTime) . '秒后再试（数据同步最小间隔为5分钟）');
        }
        (new SyncStoreDepartmentByCoach())->handle();

        return success(['last_sync_time' => getOrSetData($key)]);
    }
}
