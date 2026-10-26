<?php
/**
 * Created by :PhpStorm
 * User: lvshuaishuai
 * Date: 2020/11/24
 * Time: 2:03 下午
 */

namespace App\Http\Controllers\App\V1\DistributionPartner;


use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\Work\PartnerMember;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;

class MyController extends BaseController
{
    public function index(Request $request)
    {
        // 接收参数
        $params = $request->only([
            'partner_id',
        ]);
        $params['enterprise_id'] = $this->enterpriseId;
        // 验证规则
        $rules = [
            'partner_id' => 'required|integer',
        ];
        // 提示语
        $messages = [
            'partner_id.required' => '合伙人ID不能为空',
            'partner_id.integer' => '合伙人ID错误',
        ];
        // 参数验证
        $validator = Validator::make($params, $rules, $messages);

        if ($validator->fails()) {
            throw new Error(500, $validator->errors()->first());
        }


        $partnerMyInfo = (new PartnerMember())->partnerMy($params);

        return success($partnerMyInfo ?? []);
    }
}