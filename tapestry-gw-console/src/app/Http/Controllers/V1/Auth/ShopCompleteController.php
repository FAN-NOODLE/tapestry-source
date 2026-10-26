<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/12/3
 * Time: 4:34 PM
 */

namespace App\Http\Controllers\V1\Auth;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Boss\Enterprise;
use CrsSdk\Work\EnterprisePlugin;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class ShopCompleteController extends BaseController
{
    /**
     * 商城小程序完善接口
     *
     * @param Request $request
     * @return JsonResponse
     * @throws Exception
     */
    public function index(Request $request)
    {
        $params = $request->all();

        // 去除参数首尾空格
        foreach ($params as $key => $value) {
            $params[$key] = trim($value);
        }

        if (empty($params['app_secret']) || strlen($params['app_secret']) > 32) {
            throw new Error(100000, '小程序secret有误');
        }

        $params['qy_domain'] = $this->qyDomain;
        $params['corp_id']   = $this->enterprise['corp_id'] ?? '';
        // 应用标识，用于区分是哪个小程序
        $params['app_identifier'] = 'crs-shop';

        $result = [];
        // 自建应用安装
        // 小程序secret、外部联系人secret需要单独提交给UC
        if (!empty($params['agent_id']) && !empty($params['mini_program_secret'])) {
            $result                 = (new EnterprisePlugin())->extShopInstall($params);
            $params['pl_plugin_id'] = !empty($result['data']['mini_program_install_result']) ? $result['data']['mini_program_install_result']['plPluginid'] : '';
        }

        // 关键参数在企业信息副表中存储一下
        $params['enterprise_id'] = $this->enterpriseId;
        (new Enterprise())->shopSetting($params);

        // 日志埋点
        $this->writeOperatorLog("完善商城小程序信息");

        return success($result);
    }
}
