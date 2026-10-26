<?php
/**
 * Created By PhpStorm
 * @Date: 2021/6/30
 * @Time: 15:14
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: AuthController.php
 */

namespace App\Http\Controllers\Brand\BrandAuth;


use App\Http\Controllers\V1\BaseController;
use App\Models\BrandAuth;
use Illuminate\Support\Facades\Validator;

class AuthController extends BaseController
{
    public function index()
    {
        Validator::make($this->request->input(), [
            'enterprise_id' => 'required',
            'business'      => 'required',
        ], [
            'enterprise_id.required' => '企业id必填',
            'business.required'      => '功能标识必填',

        ])->validate();

        $enterpriseId = $this->request->input('enterprise_id');
        $business      = $this->request->input('business');
        $method       = $this->request->input('method');
        $type         = $this->request->input('type');

        $params = [
            'method' => $method,
            'type'   => $type,
        ];

        $auth = (new BrandAuth())->getBrandAuth(
            $enterpriseId,
            $business,
            $params,
            [
                'enterprise_id',
                'business',
                'method',
                'type',
                'name'
            ]
        );

        return success([
            'is_can' => $auth ? 1 : 0
        ]);
    }
}