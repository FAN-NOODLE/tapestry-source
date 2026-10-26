<?php
/**
 * Created by PhpStorm.
 * User: zp
 * Date: 19-2-27
 * Time: 下午3:37
 */

namespace App\Http\Controllers\V1\Auth;

use CrsSdk\Boss\Auth;
use Exception;
use GuzzleHttp\Client;
use CrsSdk\Boss\Enterprise;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use CrsSdk\Work\EnterprisePlugin;
use Illuminate\Support\Facades\Log;
use App\Http\Controllers\V1\BaseController;

class ShopSaveController extends BaseController
{
    /**
     * ShopSaveController constructor.
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * shop授权信息编辑
     *
     * @param Request $request
     * @return JsonResponse|void
     * @throws Exception
     */
    public function index(Request $request)
    {
        $params = $request->only(['mini_program_secret', 'pl_plugin_id', 'app_secret', 'agent_id']);

        if (empty(trim($params['app_secret'])) || strlen(trim($params['app_secret'])) > 32) {
            notice(910002, [], $params);
        }

        // 去除参数首尾空格
        foreach ($params as $key => $value) {
            $params[$key] = trim($value);
        }

        // 获取企业app_id
        try {
            $appIdArr = (new Auth())->getShopAppid(['enterprise_id' => $this->enterpriseId]);
        } catch (Exception $exception) {
            $appIdArr['data']['shop_app_id'] = '';
        }

        // boss挂了
        if (empty($appIdArr['data']['shop_app_id'])) {
            notice(600005);
        }

        // 验证secret
        try {
            $appId    = $appIdArr['data']['shop_app_id'];
            $secret   = $params['app_secret'];
            $response = (new Client())
                ->request('GET', "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid={$appId}&secret={$secret}")
                ->getBody()
                ->getContents();
            $result   = json_decode($response, true);
        } catch (Exception $e) {
            notice(790002);
        }

        // secret无效
        if (isset($result['errcode']) && in_array((int)$result['errcode'], [40125, 40001])) {
            notice(790004);
        }

        $paramsData['enterprise_id'] = $this->enterpriseId;
        $paramsData['type']          = 1;
        $authDetail                  = (new Enterprise())->authDetail($paramsData);

        try {
            // shop小程序
            $params['qy_domain']      = $this->qyDomain;
            $params['corp_id']        = $this->enterprise['corp_id'] ?? '';
            $params['app_identifier'] = 'crs-shop';
            $params['enterprise_id']  = $this->enterpriseId;
            if (!empty($params['agent_id'])) {
                if (empty($authDetail['data']['agent_id'])) {
                    $result = (new EnterprisePlugin())->extShopInstall($params);
                } else {
                    $params['pl_plugin_id'] = $authDetail['data']['pl_plugin_id'];
                    $result                 = (new EnterprisePlugin())->extShopUpdate($params);
                }
                $params['pl_plugin_id'] = !empty($result['data']['mini_program_install_result']) ? $result['data']['mini_program_install_result']['plPluginid'] : '';
            }

            // 更新Boss库
            (new Enterprise())->shopSetting($params);

            // 日志埋点
            $this->writeOperatorLog("编辑商城小程序信息");

            return success([]);

        } catch (Exception $e) {
            Log::error($e->getMessage());
            notice(100000, [], $params);
        }
    }
}
