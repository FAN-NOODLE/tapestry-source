<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2019-04-25
 * Time: 18:12
 */

namespace App\Http\Controllers\App\V1\MiniProgram;

use App\Http\Controllers\App\V1\BaseController;
use App\Models\EnterpriseMiniProgramCodeLibrary;
use App\Services\OpenPlatform\MiniProgramService;
use Illuminate\Support\Facades\Validator;

class ReleaseController extends BaseController
{
    /**
     * 发布已通过审核的小程序
     *
     * @param MiniProgramService $miniProgramService 小程序相关服务
     *
     * @return array|void
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     */
    public function index(MiniProgramService $miniProgramService)
    {
        // 接收参数
        $params = request()->only('app_id', 'release_model', 'grp_gray_percentage', 'version');
        // 参数验证规则
        $rules = [
            'app_id'              => 'required',
            'release_model'       => 'required|in:1,2',
            'grp_gray_percentage' => 'nullable|required_if:release_model,2|integer|between:1,100'
        ];
        // 提示语
        $message = [
            'app_id.required'                 => '[参数]小程序app_id不能为空',
            'release_model.required'          => '[参数]发布模式不能为空',
            'release_model.in'                => '[参数]发布模式值不合法',
            'grp_gray_percentage.required_if' => '[参数]发布比例不能为空',
            'grp_gray_percentage.between'     => '[参数]发布比例值必须为1到100的整数'
        ];
        Validator::make($params, $rules, $message)->validate();

        // 小程序app_id
        $appId = $params['app_id'];
        if (isset($params['version']) && !empty($params['version'])) {
            $auditInfo = $miniProgramService->getAuditIdByVersionAppId($params['version'], $appId);
        } else {
            // 获取代码审核id
            $auditInfo = $miniProgramService->recentAudit($appId);
        }

        // 分阶段发布
        if (EnterpriseMiniProgramCodeLibrary::RELEASE_MODEL_GRAY === intval($params['release_model'])) {
            $result = $miniProgramService->grayRelease(
                $appId,
                $params['grp_gray_percentage'],
                $auditInfo['wx_audit_id']
            );
        } else { // 全量发布
            // 发布已通过审核的小程序
            $result = $miniProgramService->release($appId, $auditInfo['wx_audit_id']);
        }

        if (!$result) {
            return error(300108);
        }

        return success([
            'code'    => $result['errcode'],
            'message' => $result['errmsg']
        ]);
    }
}
