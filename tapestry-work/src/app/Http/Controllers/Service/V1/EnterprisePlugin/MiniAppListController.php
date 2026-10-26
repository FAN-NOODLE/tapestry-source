<?php
/*
* +------------------------------------+
* | Created: phpstorm
* +------------------------------------+
* | Author : wyg
* +------------------------------------+
* | Date   : 2019/11/20
* +------------------------------------+
* | Time   : 14:15
* +------------------------------------+
* | Email  : <wangyongguo@vchangyi.com>
* +------------------------------------+
*/


namespace App\Http\Controllers\Service\V1\EnterprisePlugin;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\EnterprisePlugin;
use CrsSdk\Boss\Enterprise;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;


class MiniAppListController extends BaseController
{
    // 安装企业微信小程序地址
    const WORK_WE_CHART_MINI_APP_INSTALL_URL = '%s/s/suite/install?suiteId=%s&callbackUrl=%s';

    // 禁用企业微信小程序
    const WORK_MINI_APP_STATUS = 2;

    /**
     * 企业微信安装小程序列表
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $domain = $request->input('qy_domain', '');

        $params = $request->input('params', []);

        $rules = [
            'enterprise_id' => 'required',
            'channel_type'  => 'required',
        ];

        $message = [
            'enterprise_id.required' => '[enterpriseId]不能为空',
            'channel_type.required'  => '[channel_type]不能为空',
        ];

        Validator::make($params, $rules, $message)->validate();

        $env        = env('APP_ENV');
        $workPlugin = config('work_plugin.crs-guide');
        $ucApi      = config('client.api_url.uc_api');
        $workPlugin = $workPlugin[$env];

        $enterprise = (new Enterprise())->detail(['enterprise_id' => $params['enterprise_id']]);

        $enterprise = $enterprise['data'] ?? [];

        $isOpen = $enterprise['crs_status'] ?? 1;

        $callbackUrl = $workPlugin['callback_url'] . '?hash=' . $enterprise['hash'];

        $args = [$ucApi, $workPlugin['suite_id'], $callbackUrl];

        $installUrl = vsprintf(self::WORK_WE_CHART_MINI_APP_INSTALL_URL, $args);

        $default = [
            'logo'      => $workPlugin['app_logo'],
            'app_name'  => $workPlugin['app_name'],
            'installed' => false,
            'auth_url'  => $installUrl,
            'is_open'   => $isOpen
        ];

        $list = [];
        if (empty(trim($domain)) || self::WORK_MINI_APP_STATUS == $isOpen) {
            return success(['data' => [$default]]);
        }

        try {
            $result = (new EnterprisePlugin())->list($domain, ['channelType' => $params['channel_type']]);
            unset($result['requestId']);
            foreach ($result as $miniApp) {
                // 通讯录暂时不需要返回
                if ('通讯录' == $miniApp['plName']) {
                    continue;
                }

                // dev排除test的
                if ('dev' == $env && '测试-CRS小程序平台版' != trim($miniApp['plName'])) {
                    $list = self::getList($miniApp, $installUrl, $isOpen);
                }

                // test排除dev的
                if ('test' == $env && 'CRS小程序平台版' != trim($miniApp['plName'])) {
                    $list = self::getList($miniApp, $installUrl, $isOpen);
                }

                // 生产环境不变
                if ('prod' == $env) {
                    $list = self::getList($miniApp, $installUrl, $isOpen);
                }
            }
        } catch (\Exception $e) {
            return success(['data' => [$default]]);
        }

        return success(['data' => $list]);
    }


    /**
     * 格式化列表
     *
     * @param $miniApp
     * @param $installUrl
     * @param $isOpen
     *
     * @return array
     */
    private function getList($miniApp, $installUrl, $isOpen)
    {
        $list = [];
        // installed 是否安装企业微信小程序:false 未安装, true 已安装
        if ($miniApp['available'] == 1) {
            $list[] = [
                'logo'      => $miniApp['plIcon'],
                'app_name'  => $miniApp['plName'],
                'installed' => true,
                'auth_url'  => '',
                'is_open'   => $isOpen
            ];
        } else {
            $list[] = [
                'logo'      => $miniApp['plIcon'],
                'app_name'  => $miniApp['plName'],
                'installed' => false,
                'auth_url'  => $installUrl,
                'is_open'   => $isOpen
            ];
        }
        return $list;
    }
}
