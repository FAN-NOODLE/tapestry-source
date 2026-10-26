<?php
/**
 * Created by PhpStorm.
 * User: zhaoqi
 * Date: 2019-02-27
 * Time: 11:42
 */

namespace App\Http\Controllers\V1\Auth;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Boss\Enterprise;
use CrsSdk\Work\EnterprisePlugin;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class AuthCompleteController extends BaseController
{
    /**
     * AuthCompleteController constructor.
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * 授权完善信息
     *
     * @param Request $request
     * @return
     */
    public function index(Request $request)
    {
        $params = $request->only(['corp_id', 'work_program_secret', 'external_secret', 'shop_program_secret', 'crs_status', 'store_status', 'crs_display', 'store_display','agent_id']);
        if ($params['crs_status'] == 1 || $params['store_status'] == 1) {
            if (empty(trim($params['corp_id']))) {
                notice(910000, [], $params);
            }
        }

        if ($params['crs_status'] == 1 && $params['crs_display'] == true) {
            if (empty(trim($params['work_program_secret']))) {
                notice(910002, [], $params);
            }

            if (empty(trim($params['external_secret']))) {
                notice(910001, [], $params);
            }
        }

        if ($params['store_status'] == 1 && $params['store_display'] == true) {
            if (empty(trim($params['shop_program_secret']))) {
                notice(910003, [], $params);
            }
        }


        try {
            $qyDomain = $this->qyDomain;
            $enterpriseId = $this->enterpriseId;

            if ($params['store_status'] == 1 && $params['store_display'] == true) {
                // shop小程序
                $shopData['qy_domain'] = $qyDomain;
                $shopData['app_identifier'] = 'crs-shop';
                $shopData['corp_id'] = trim($params['corp_id']);
                $shopData['mini_program_secret'] = trim($params['shop_program_secret']);
                $shopData['enterprise_id'] = $enterpriseId;
                $shopData['agent_id'] = trim($params['agent_id']);

                // 小程序secret、外部联系人secret需要单独提交给UC
                $shopResult = (new EnterprisePlugin())->extShopInstall($shopData);
                $shopData['ep_enumber'] = $shopResult['data']['mini_program_install_result']['epEnumber'];
                $shopData['pl_plugin_id'] = $shopResult['data']['mini_program_install_result']['plPluginid'];

                (new Enterprise())->shopSetting($shopData);
            }

            if ($params['crs_status'] == 1 && $params['crs_display'] == true) {
                // work小程序
                $workData['qy_domain'] = $qyDomain;
                $workData['app_identifier'] = 'crs-mini-program';
                $workData['enterprise_id'] = $enterpriseId;
                $workData['corp_id'] = trim($params['corp_id']);
                $workData['mini_program_secret'] = trim($params['work_program_secret']);
                $workData['external_secret'] = trim($params['external_secret']);

                // 小程序secret、外部联系人secret需要单独提交给UC
                $workResult = (new EnterprisePlugin())->extInstall($workData);
                $workData['ep_enumber'] = $workResult['data']['mini_program_install_result']['epEnumber'];
                $workData['pl_plugin_id'] = $workResult['data']['mini_program_install_result']['plPluginid'];

                (new Enterprise())->workSetting($workData);
            }
        } catch (\Exception $e) {
            Log::error($e->getMessage());
            notice(910004, [$e->getMessage()], $params);
        }

        // 日志埋点
        $this->writeOperatorLog("完善授权信息");

        return success();
    }
}
