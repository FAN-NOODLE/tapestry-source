<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2019-04-25
 * Time: 18:11
 */

namespace App\Http\Controllers\Service\V1\MiniProgram;

use Illuminate\Support\Facades\Validator;

class SubmitAuditController extends MiniBaseController
{
    /**
     * 小程序提审
     * @return array|void
     */
    public function index()
    {
        // 接收参数
        $params = request()->only('version', 'online_type');
        $params['app_id'] = $this->appId;
        // 参数验证规则
        $rules = [
            'version'     => 'required',
            'online_type' => 'required|in:1,2'
        ];
        // 提示语
        $message = [
            'version.required'     => '[参数]提审版本号不能为空',
            'online_type.required' => '[参数]上线方式不能为空',
            'online_type.in'       => '[参数]上线方式值不合法'
        ];
        Validator::make($params, $rules, $message)->validate();

        try {
            // 提审
            $result = $this->miniProgramService->submitSpecifiedVersion(
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
