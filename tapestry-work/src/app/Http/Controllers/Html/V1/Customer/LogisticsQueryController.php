<?php

namespace App\Http\Controllers\Html\V1\Customer;

use App\Http\Controllers\Html\V1\BaseController;
use App\Models\RefundApply;
use CrsSdk\Open\Logistics;
use Illuminate\Http\Request;

class LogisticsQueryController extends BaseController
{
    /**
     * 批量查询物流接口
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $input = $request->input();

        $this->validate($request, [
            'data'                => 'required|array',
            'data.*.logistic_key' => 'required',
            'data.*.logistic_no'  => 'required',
        ], [
            'data.required'                => 'data必填',
            'data.array'                   => 'data需为array',
            'data.*.logistic_key.required' => 'logistic_key必填',
            'data.*.logistic_no.required'  => 'logistic_no必填',
        ]);

        $params = [
            'enterprise_hash' => $this->enterprise['hash'],
            'data'            => []
        ];
        foreach ($input['data'] as $key => $value) {
            $params['data'][] = [
                'logistic_key'    => $this->matchKey($value['logistic_key']),
                'logistic_no'     => $value['logistic_no'],
                'sender_mobile'   => $value['sender_mobile'] ?? '',
                'receiver_mobile' => $value['receiver_mobile'] ?? ''
            ];
        }

        $data = (new Logistics())->query($params);

        return success($data);
    }

    /**
     * 匹配新的物流公司标识
     *
     * @param $key string 物流公司标识
     *
     * @return mixed
     */
    private function matchKey($key)
    {
        $arr = [
            'SF'   => 'sf',
            'STO'  => 'sto',
            'YTO'  => 'yt',
            'YD'   => 'yd',
            'HHTT' => 'tt',
            'EMS'  => 'ems',
            'ZTO'  => 'zto',
            'HTKY' => 'ht',
            'DBL'  => 'db',
            'JD'   => 'jd',
            'YZPY' => 'youzheng',
            'BTWL' => 'bsky',
            'SNWL' => 'suning'
        ];

        return $arr[$key] ?? $key;
    }
}
