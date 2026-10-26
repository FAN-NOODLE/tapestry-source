<?php
/**
 * Created by :PhpStorm
 * User: lvshuaishuai
 * Date: 2020/11/18
 * Time: 5:53 下午
 */

namespace App\Http\Controllers\App\V1\DistributionPartner;


use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\Work\PartnerMember;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;
use App\Exceptions\Error;

class DisableDetailController extends BaseController
{

    public function index(Request $request)
    {
        $params = $request->only([
            'partner_id'
        ]);

        // 验证规则
        $rules = [
            'partner_id' => 'integer'
        ];

        // 提示语
        $messages = [
            'partner_id.integer' => '合伙人ID类型为int'
        ];
        // 参数验证
        $validator = Validator::make($params, $rules, $messages);

        if ($validator->fails()) {
            throw new Error(500, $validator->errors()->first());
        }
        $params['enterprise_id'] = $this->enterpriseId;
        $result = (new PartnerMember())->partnerDisableDetail($params);
        return success($result ?? []);
    }
}