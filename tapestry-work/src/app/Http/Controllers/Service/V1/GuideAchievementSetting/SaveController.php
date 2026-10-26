<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/26
 * Time: 10:04
 */

namespace App\Http\Controllers\Service\V1\GuideAchievementSetting;

use App\Common\Constant;
use App\Services\GuideAchievementSettingService;
use Illuminate\Http\Request;
use App\Http\Controllers\Service\V1\BaseController;

class SaveController extends BaseController
{

    /**
     * 设置保存
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request)
    {
        $params = $request->only(['enterprise_id', 'binding_type', 'duration_time', 'change_owner', 'protected_time', 'commission_type', 'fission_type']);
        $rules    = [
            'enterprise_id' => 'required',
            'binding_type' => 'required',
            'change_owner' => 'required',
            'commission_type' => 'required',
            'fission_type' => 'required',
        ];
        $messages = [
            'enterprise_id.required' => '企业ID不能为空',
            'binding_type.required' => '绑定时效不能为空',
            'change_owner.required' => '是否允许抢客不能为空',
            'commission_type.required' => '导购自己购买商品计算业绩不能为空',
            'fission_type.required' => '客户分享统计导购业绩不能为空',
        ];
        $this->validate($this->request, $rules, $messages);

        $params['type'] = Constant::GUIDE_TYPE;

        (new GuideAchievementSettingService())->save($params);

        return success([]);
    }
}
