<?php
/**
 * Created by PhpStorm.
 * User: zhaoyang
 * Date: 2020/10/29
 * Time: 12:42
 */

namespace App\Http\Controllers\Open\V1\Guide;

use App\Models\CrsWorkerThirdPartRelation;
use App\Services\GuideAuthService;
use Illuminate\Http\Request;
use App\Http\Controllers\Open\V1\BaseController;
use Illuminate\Support\Facades\Log;

class BindUserController extends BaseController
{

    /**
     * @param Request $request
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     * @throws \Exception
     * User: zhaoyang@vchangyi.com
     * Date: 2020/10/29 13:11
     */
    public function index(Request $request)
    {
        Log::info('员工商城认证异步通知数据BindUserController', [$request->input()]);

        $rules = [
            'mem_uid'     => 'required',
            'source_user_id' => 'required',
            'source_type' => 'required|in:' . implode(',', CrsWorkerThirdPartRelation::SOURCE_LIST),
        ];

        $message = [
            'mem_uid.required'     => 'mem_uid必传',
            'source_user_id.required' => 'source_user_id必传',
            'source_type.required' => 'source_type必传',
            'source_type.in' => 'source_type参数错误',
        ];

        $input = $this->validate($request, $rules, $message);

        switchAppHost(['enterprise_id' => $this->enterpriseInfo['id']]);

        (new GuideAuthService())->bindUser($this->enterpriseInfo['id'], $input['mem_uid'], $input['source_type'], $input['source_user_id']);
        return $this->success();
    }
}
