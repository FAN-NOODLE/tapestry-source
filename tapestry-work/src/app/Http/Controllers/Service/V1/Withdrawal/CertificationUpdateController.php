<?php
/**
 * 提现认证 - 更新
 * Created by PhpStorm.
 * User: Yang Hui
 * Date: 2020年10月16日16:56:4
 */

namespace App\Http\Controllers\Service\V1\Withdrawal;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Enterprise;
use App\Models\DistributionGoods;
use App\Models\StoreWorkerCertification;
use Illuminate\Http\Request;
use Illuminate\Validation\Rule;

class CertificationUpdateController extends BaseController
{
    /**
     * 提现认证 - 更新
     * @param Request $request
     * @return mixed
     * @throws \Exception
     * @author   Yang Hui
     * @date     2020年10月16日17:1:28
     */
    public function index(Request $request)
    {
        $params = $request->only([
            'id',
            'enterprise_id',
            'worker_id',
            'real_name',
            'mobile',
            'id_type',
            'id_code',
            'bank_card_username',
            'bank_card_code',
            'bank_name',
            'agree_deal'
        ]);

        $rules = [
            'id' => 'required|integer',
            'enterprise_id' => 'required|integer',
            'worker_id' => 'required|integer',
            'real_name' => [
                'required',
                'regex:/^[\x{4e00}-\x{9fa5}A-Za-z]{1,20}/u'
            ],
            'mobile' => ['required','regex:/^1[3|4|5|6|7|8|9]\d{9}$/'],
            'id_type' => [
                'required',
                Rule::in([
                    StoreWorkerCertification::ID_TYPE_MAINLAND_RESIDENT,
                    StoreWorkerCertification::ID_TYPE_HK_MA_PASS,
                    StoreWorkerCertification::ID_TYPE_HK_MA_RESIDENT,
                    StoreWorkerCertification::ID_TYPE_TW_PASS,
                    StoreWorkerCertification::ID_TYPE_TW_RESIDENT
                ])
            ],
            'id_code' => 'required',
            'bank_card_username' => [
                'required',
                'regex:/^[\x{4e00}-\x{9fa5}A-Za-z]{1,20}/u'
            ],
            'bank_card_code' => 'required',
            "bank_name" => 'required',
            'agree_deal' => 'required|in:1'
        ];

        $message = [
            'id.required' => 'id参数错误',
            'id.integer' => 'id参数错误',
            'enterprise_id.required' => '企业ID不存在',
            'enterprise_id.integer'  => '企业ID数据类型有误',
            'worker_id.required'      => '导购ID不存在',
            'worker_id.integer' => '导购ID数据类型错误',
            'real_name.required' => '姓名不能为空',
            'real_name.regex' => '姓名只能为汉字或者字母',
            'mobile.required' => '手机号码不为空',
            'mobile.regex' => '手机号码格式错误',
            'id_type.required' => '请选择证件类型',
            'id_type.in' => '证件类型错误',
            'id_code.required' => '证件号码不能为空',
            'bank_card_username.required' => '户名不能为空',
            'bank_card_username.regex' => '户名只能为汉字或者字母',
            'bank_card_code.required' => '银行卡号不能为空',
            'agree_deal.required' => '请阅读并同意协议',
            'bank_name.required' => '请选择银行名称',
            'agree_deal.in' => '请阅读并同意协议'
        ];

        $this->validate(request(), $rules, $message);

        $modelObj = StoreWorkerCertification::query()
            ->where('id', $params['id'])
            ->where('worker_id', $params['worker_id'])
            ->where('enterprise_id', $params['enterprise_id'])
            ->first();
        if (empty($modelObj)) {
            return error(1100002);
        }

        $modelObj->enterprise_id = $params['enterprise_id'];
        $modelObj->worker_id = $params['worker_id'];
        $modelObj->real_name = $params['real_name'];
        $modelObj->mobile = $params['mobile'];
        $modelObj->id_type = $params['id_type'];
        $modelObj->id_code = $params['id_code'];
        $modelObj->bank_card_username = $params['bank_card_username'];
        $modelObj->bank_card_code = $params['bank_card_code'];
        $modelObj->bank_name = $params['bank_name'];
        $modelObj->agree_deal = $params['agree_deal'];
        //更新时是默认，如果不是在验证通过时，status一律为禁用
        $modelObj->auth_status = StoreWorkerCertification::AUTH_STATUS_DEFAULT;
        $modelObj->status = StoreWorkerCertification::STATUS_DISABLED;  //未认证通过时禁用

        $result = $modelObj->save();

        if ($result) {
            $returnData = [
                'data' => ['id' => $modelObj->id]
            ];
            return success($returnData);
        } else {
            return error(1100001);
        }
    }
}
