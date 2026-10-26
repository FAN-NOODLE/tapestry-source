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

class DetailController extends BaseController
{

    /**
     * 设置详情
     */
    public function index()
    {
        $detail = (new GuideAchievementSetting())->detail(['enterprise_id' => $this->enterpriseId]);

        return success($detail['data']);
    }
}
