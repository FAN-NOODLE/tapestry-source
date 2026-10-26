<?php
/**
 * Created by PhpStorm.
 * User: zp
 * Date: 19-2-27
 * Time: 下午3:50
 */

namespace App\Http\Controllers\Service\V1\EnterprisePlugin;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\EnterprisePlugin;
use App\Library\Uc\Member;
use Illuminate\Http\Request;

class ExtShopUpdateController extends BaseController
{
    /**
     * 扩展应用更新接口
     *
     * @param Request $request
     * @return array|\Psr\Http\Message\StreamInterface
     * @throws \Exception
     * @author   zhaoqi  zhaoqi@vchangyi.com
     * @date     2019/02/27 17:24 PM
     */
    public function index(Request $request)
    {
        $miniProgramSecret = $request->post('mini_program_secret');
        $qyDomain          = $request->post('qy_domain');
        $plPluginId        = $request->post('pl_plugin_id');
        $agentId           = $request->post('agent_id', 0);

        if (empty($qyDomain) || empty($miniProgramSecret) || empty($plPluginId)) {
            throw new Error(100000, '必填参数错误');
        }

        $enterprisePlugin = new EnterprisePlugin();

        // 安装小程序应用
        $miniInstallResult = $enterprisePlugin->appUpdate($qyDomain,
            [
                'plPluginid' => $plPluginId,
                'plSecret'   => $miniProgramSecret,
                'agentid'    => $agentId
            ]
        );

        //调UC同步通讯录
        (new Member())->syncContact($qyDomain);

        return success([
            'data' => [
                'mini_program_install_result' => $miniInstallResult
            ]
        ]);
    }
}