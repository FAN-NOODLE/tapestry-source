<?php
/**
 * Created by PhpStorm.
 * User: admin
 * Date: 2018/12/27
 * Time: 16:38
 */

namespace App\Http\Controllers\V1\Setting;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Enterprise as WorkEnterprise;

class CallbackController extends BaseController
{
    /**
     * @return \Illuminate\Http\JsonResponse
     */
    public function index()
    {
        $corpId = $this->enterprise['corp_id'];

        // 该企业未授权CRS小程序
        if (empty($this->qyDomain) || empty($corpId)) {
            return success([
                'uc_callback_url' => '',
                'token'           => '',
                'key'             => '',
                'message'         => '该企业未授权CRS小程序'
            ]);
        }

        $callbackUrl = config('app.uc_customer_callback_url');

        $workEnterprise = new WorkEnterprise();

        $result = $workEnterprise->postUrl(['qy_domain' => $this->qyDomain]);

        $callbackUrl = vsprintf($callbackUrl, [$corpId]);

        $info['key']             = $result['data']['key']['wxqyEncodingAeskey'];
        $info['token']           = $result['data']['token']['wxqyToken'];
        $info['uc_callback_url'] = $callbackUrl;

        return success($info);
    }
}
