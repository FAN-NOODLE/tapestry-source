<?php
/**
 * Created by :PhpStorm
 * User: lvshuaishuai
 * Date: 2020/11/18
 * Time: 10:31 上午
 */

namespace App\Http\Controllers\App\V1\DistributionPartner;


use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;

class ListController extends BaseController
{
    public function index(Request $request)
    {
        // 接收参数
        $params = $request->only([
            'store_id',
            'enterprise_id',
            'worker_id',
            'partner_status',
            'worker_name',
            'parnter_name',
            'parnter_mobile',
            'page',
            'size'
        ]);

        // 验证规则
        $rules = [
            'enterprise_id'  => 'integer',
            'parnter_name'   => 'string',
            'parnter_mobile' => 'integer',
            'worker_id'      => 'integer'
        ];

        // 提示语
        $messages = [
            'parnter_mobile.integer' => '手机号需为int类型',
            'worker_id.integer'      => '导购id需int类型'
        ];
        // 参数验证
        $validator = Validator::make($params, $rules, $messages);

        if ($validator->fails()) {
            throw new Error(500, $validator->errors()->first());
        }
        $params['size'] = $params['size'] ?? 10;
        $params['page'] = $params['page'] ?? 1;
        $this->input    = $params;


    }
}