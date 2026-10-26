<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/26
 * Time: 10:04
 */

namespace App\Http\Controllers\Service\V1\GuideAchievementSetting;

use App\Services\GuideAchievementSettingService;
use Illuminate\Http\Request;
use App\Http\Controllers\Service\V1\BaseController;

class DetailController extends BaseController
{

    /**
     * 设置详情
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     */
    public function index(Request $request)
    {
        $params = $request->only(['enterprise_id']);
        $detail = (new GuideAchievementSettingService())->detail($params);

        return success(['data' => $detail]);
    }
}
