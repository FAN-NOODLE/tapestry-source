<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/12/3
 * Time: 7:15 PM
 */

namespace App\Http\Controllers\Service\V2\Auth;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\EnterprisePlugin;
use App\Models\Enterprise;
use Illuminate\Http\Request;

class WorkSettingController extends BaseController
{
    /**
     * 保存企业填写的外部联系人的接口调用凭证信息
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse|void
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $enterpriseId = $request->post('enterprise_id');
        $externalSecret = $request->post('external_secret');

        $enterprise = Enterprise::query()
            ->where('id', $enterpriseId)
            ->first();

        if (empty($enterprise)) {
            // 企业信息不存在
            return error(600002);
        }

        $enterprise->update(['external_secret' => $externalSecret]);

        $this->extInstall($enterprise,$externalSecret);


        return successV2();
    }


    /**
     * @param $enterprise
     * @param $externalSecret
     * @throws \Exception
     */
    protected function extInstall($enterprise, $externalSecret)
    {
        if (empty($enterprise['domain'])||empty($enterprise['corp_id'])) {
            // 企业domain或corpId不存在
            return error(600007);
        }
        $enterprisePlugin = new EnterprisePlugin();

        $ucEnterprise = new \App\Library\Uc\Enterprise();
        // 企业微信传送给UC的回调地址

        $callback = env('WORK_EXTERNAL_CALLBACK', 'https://work-api-crs.vchangyi.com/callback/v1/customer/add?enterprise_id=') . $enterprise['id'];

        $ucEnterprise->entSettingModify($enterprise['domain'], $callback);

        // 写入外部联系人secret
        $enterprisePlugin->modify($enterprise['domain'],
            [
                'corpId' => $enterprise['corp_id'],
                'wxExternalSecret' => $externalSecret
            ]
        );
    }
}
