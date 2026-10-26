<?php
/**
 * Created By PhpStorm
 * @Date: 2021/5/19
 * @Time: 9:24
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: ConfigSaveController.php
 */

namespace App\Http\Controllers\Service\V1\MissionCenter;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\MissionCenterAppletsConfig;
use App\Services\MissionCenterTaskService;
use Illuminate\Http\Request;
use Illuminate\Validation\Rule;

class ConfigSaveController extends BaseController
{
    public function index(Request $request, MissionCenterTaskService $missionCenterService)
    {
        $rules = [
            'name'            => 'required',
            'enterprise_id'   => 'required',
            'app_id'          => 'required',
            'app_original_id' => 'required',
            'applet_path'     => 'required',
            'type'            => [
                'required',
                Rule::in([
                    MissionCenterAppletsConfig::TYPE_TRAINING_TASK,
                    MissionCenterAppletsConfig::TYPE_INVITATION_TASK
                ])
            ],
        ];

        $message = [
            'name.*'            => '小程序名称不能为空',
            'enterprise_id.*'   => '企业ID不能为空',
            'app_id.*'          => 'APP ID不能为空',
            'app_original_id.*' => '小程序原始ID不能为空',
            'applet_path.*'     => '小程序跳转路径不能为空',
            'type.*'            => '配置类型必须为1或2'
        ];
        $this->validate($request, $rules, $message);
        $input = $request->all();

        $result = $missionCenterService->saveConfig($input);

        return success($result);
    }
}