<?php
/**
 * Created by :PhpStorm
 * User: lvshuaishuai
 * Date: 2020/11/24
 * Time: 5:48 下午
 */

namespace App\Http\Controllers\App\V1\Common;


use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use App\Services\Decrypt;
use CrsSdk\Work\PartnerMember;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;

class CreateController extends BaseController
{
    public function index(Request $request)
    {
        // 接收参数
        $params = $request->only([
            'enterprise_id',
            'path',
            'params',
        ]);

        // 验证规则
        $rules = [
            'path'   => 'required|string',
            'params' => 'array'
        ];

        [$appid, $secret] = Decrypt::getEnvMiniProgram($this->identifier);

        $params['params']['app_id'] = $appid;
        $params['params']['secret'] = $secret;

        // 提示语
        $messages = [
            'path.string'      => '海报跳转路径必传',
            'params.array'     => '参数为数组格式',
            'params.mem_uid'   => '参数下的导购memuid必传',
            'params.store_id'  => '参数下的导购导购门店ID必传',
            'params.worker_id' => '参数下的导购导购ID必传'
        ];
        // 参数验证
        $validator = Validator::make($params, $rules, $messages);

        if ($validator->fails()) {
            throw new Error(500, $validator->errors()->first());
        }
        $res = (new PartnerMember())->partnerSunCode($params);

        return success($res);
    }
}