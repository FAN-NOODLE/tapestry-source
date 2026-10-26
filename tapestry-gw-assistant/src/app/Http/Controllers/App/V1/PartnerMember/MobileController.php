<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/11/19 10:06
 */

namespace App\Http\Controllers\App\V1\PartnerMember;


use App\Http\Controllers\App\V1\BaseController;
use App\Services\Decrypt;
use App\Services\Wechat\WxDecrypt;
use Illuminate\Http\Request;


class MobileController extends BaseController
{
    protected $isRequireLogin = false;

    public function index(Request $request)
    {
        $this->validate($request, [
            'encryptedData' => 'required',
            'iv'            => 'required',
            'cache_key' => 'required'//获得的微信用户的数据信息
        ], [
            'encryptedData.required' => '加密数据不能为空',
            'iv.required'            => 'iv不能为空',
            'cache_key.*' => 'cache_key必须存在',
        ]);

        $input = $this->request->all();
        //解参
        $cache = (new Decrypt())->index($input['cache_key']);
        //微信解密获得手机号
        $data = (new WxDecrypt($cache['sessionKey'], $this->identifier))->index($input);
        $data = json_decode($data,true);
        $mobile['phoneNumber'] = $data['phoneNumber'];
        return success(['data'=>$mobile]);
    }

}