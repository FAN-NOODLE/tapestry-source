<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/12/10
 * Time: 16:24
 */

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Console\Commands\UpdateStoreManagerByPosition;
use App\Http\Controllers\Service\V1\BaseController;

use Illuminate\Http\Request;


class UpdateStoreManagerByPositionController extends BaseController
{
    /**
     * 根据职位更新管理范围
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request): array
    {
        $enterpriseId = $request->post('enterprise_id', 0);
        $isUpdate = $request->post('is_update', 0);
        if ($isUpdate && $enterpriseId == env('COACH_ENTERPRISE_ID')) {
            (new UpdateStoreManagerByPosition())->handle();
        }
        $guideSettingKey = redisKey('UPDATE_STORE_MANAGER_BY_POSITION_TIME', ['enterprise_id' => $enterpriseId]);
        return success(['last_update_time' => getOrSetData($guideSettingKey)]);
    }
}
