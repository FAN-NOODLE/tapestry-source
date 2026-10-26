<?php
/**
 * Created by PhpStorm.
 * User: zp
 * Date: 19-2-27
 * Time: 下午3:49
 */

namespace App\Http\Controllers\Service\V1\EnterprisePlugin;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\EnterprisePlugin;
use Illuminate\Http\Request;
use App\Library\Uc\Enterprise;

class ExtUpdateController extends BaseController
{
    /**
     * 扩展应用更新接口
     *
     * @param Request $request
     * @return array|\Psr\Http\Message\StreamInterface
     * @throws \Exception
     * @author   zhaoqi zhaoqi@vchangyi.com
     * @date     2018/12/3 4:54 PM
     */
    public function index(Request $request)
    {
        $miniProgramSecret = $request->post('mini_program_secret', '');
        $externalSecret    = $request->post('external_secret');
        $corpId            = $request->post('corp_id');
        $qyDomain          = $request->post('qy_domain');
        $plPluginId        = $request->post('pl_plugin_id');
        $enterpriseId      = $request->post('enterprise_id');

        if (empty($qyDomain) || empty($externalSecret) || empty($corpId) || empty($plPluginId)) {
            throw new Error(100000, '必填参数错误');
        }

        $enterprisePlugin = new EnterprisePlugin();

        if ($miniProgramSecret) {
            // 安装小程序应用
            $miniInstallResult = $enterprisePlugin->appUpdate($qyDomain,
                [
                    'plSecret'            => $miniProgramSecret,
                    'plPluginid'            => $plPluginId,
                ]
            );
        }

        $enterprise = new Enterprise();
        // 企业微信传送给UC的回调地址

        $callback = "https://" . $_SERVER['HTTP_HOST'] . '/callback/v1/customer/add?enterprise_id=' . $enterpriseId;
        $enterprise->entSettingModify($qyDomain, $callback);

        // 写入外部联系人secret
        $externalResult = $enterprisePlugin->modify($qyDomain,
            [
                'corpId'           => $corpId,
                'wxExternalSecret' => $externalSecret
            ]
        );

        return success([
            'data' => [
                'mini_program_install_result' => $miniInstallResult ?? '',
                'external_install_result'     => $externalResult
            ]
        ]);
    }
}