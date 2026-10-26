<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2019-04-25
 * Time: 18:11
 */

namespace App\Http\Controllers\App\V1\MiniProgram;

use App\Http\Controllers\App\V1\BaseController;
use App\Services\OpenPlatform\MiniProgramService;
use Illuminate\Support\Facades\Validator;

class SubmitAuditController extends BaseController
{
    /**
     * 小程序提审
     *
     * @param MiniProgramService $miniProgramService
     *
     * @return array|void
     */
    public function index(MiniProgramService $miniProgramService)
    {
        // 接收参数
        $params = request()->only('app_id', 'version', 'online_type');
        // 参数验证规则
        $rules = [
            'app_id'      => 'required',
            'version'     => 'required',
            'online_type' => 'required|in:1,2'
        ];
        // 提示语
        $message = [
            'app_id.required'      => '[参数]小程序app_id不能为空',
            'version.required'     => '[参数]提审版本号不能为空',
            'online_type.required' => '[参数]上线方式不能为空',
            'online_type.in'       => '[参数]上线方式值不合法'
        ];
        Validator::make($params, $rules, $message)->validate();

        try {
            // 提审
            $result = $miniProgramService->submitSpecifiedVersion(
                $params['app_id'],
                $params['version'],
                true,
                $params['online_type']
            );
        } catch (\Exception $e) {
            app('log')->error($e);
            $result = false;
        }

        if (!$result) {
            return error(300106);
        }

        return success([
            'code'    => 0,
            'message' => 'ok'
        ]);
    }
}
