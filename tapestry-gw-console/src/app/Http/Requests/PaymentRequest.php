<?php
namespace App\Http\Requests;

use Illuminate\Support\Facades\Validator;

class PaymentRequest
{
    /**
     * 企业支付数据提交
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/12/06 11:16
     * @param array $data 验证的数据
     * @return \Illuminate\Contracts\Validation\Validator
     */
    public function common($data = [])
    {
        return Validator::make($data,[
            'merchant_id' => 'required',
            'merchant_key' => 'required',
            'private_key' => 'required|file',
            'public_key' => 'required|file',
        ],[
            'merchant_id.required' => '商户ID不能为空',
            'merchant_key.required' => '微信支付API key不能为空',
            'private_key.required' => '微信支付证书私钥不能为空',
            'private_key.file' => '微信支付证书私钥必须为文件',
            'public_key.required' => '微信支付证书不能为空',
            'public_key.file' => '微信支付证书必须为文件',
        ]);
    }
}