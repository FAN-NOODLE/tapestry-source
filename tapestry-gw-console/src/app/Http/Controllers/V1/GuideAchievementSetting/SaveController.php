<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/26
 * Time: 10:04
 */

namespace App\Http\Controllers\V1\GuideAchievementSetting;

use CrsSdk\Work\GuideAchievementSetting;
use App\Http\Controllers\V1\BaseController;

class SaveController extends BaseController
{

    /**
     * 设置保存
     */
    public function index()
    {
        $params = request()->only(['binding_type', 'duration_time', 'change_owner', 'protected_time', 'commission_type', 'fission_type']);

        $params['enterprise_id'] = $this->enterpriseId;
        (new GuideAchievementSetting())->save($params);

        return success([]);
    }
}
