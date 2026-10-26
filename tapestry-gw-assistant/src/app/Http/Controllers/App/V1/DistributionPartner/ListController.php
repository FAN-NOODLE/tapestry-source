<?php
/**
 * Created by :PhpStorm
 * User: lvshuaishuai
 * Date: 2020/11/18
 * Time: 5:53 下午
 */

namespace App\Http\Controllers\App\V1\DistributionPartner;


use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\Work\PartnerMember;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;

class ListController extends BaseController
{
    private $input = [];

    public function index(Request $request)
    {
        // 接收参数
        $params = $request->only([
            'store_id',
            'worker_id',
            'partner_status',
            'worker_name',
            'partner_name',
            'partner_mobile',
            'partner_ids',  //合伙人id数组
            'page',
            'size'
        ]);

        // 验证规则
        $rules = [
            'store_id'       => 'integer',
            'partner_name'   => 'string',
            'partner_mobile' => 'integer',
            'worker_id'      => 'integer'
        ];

        // 提示语
        $messages = [
            'store_id.integer'       => '门店id需为int类型',
            'partner_mobile.integer' => '手机号需为int类型',
            'worker_id.integer'      => '导购id需为int类型'
        ];
        // 参数验证
        $validator = Validator::make($params, $rules, $messages);

        if ($validator->fails()) {
            throw new Error(500, $validator->errors()->first());
        }
        $params['size']          = $params['size'] ?? 10;
        $params['page']          = $params['page'] ?? 1;
        $this->input             = $params;
        $params['enterprise_id'] = $this->enterpriseId;
        if(!isset($params['worker_id']) || empty($params['worker_id'])){
            $params['worker_id']     = $this->storeWorkerDetail($this->enterpriseId, $params['store_id'] ?? 0)['id'];
        }
        $result = (new PartnerMember())->partnerList($params);
        return success($result ?? []);
    }
}