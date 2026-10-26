<?php
/**
 * Created by PhpStorm.
 * User: zp
 * Date: 19-3-5
 * Time: 上午11:58
 */

namespace App\Http\Controllers\Service\V1\Auth;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Enterprise;
use App\Models\EnterpriseWorkPlugin;
use App\Models\EnterpriseWorkSetting;
use App\Models\EnterpriseThirdPartAccount;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class AuthDetailController extends BaseController
{
    public function __construct()
    {
        parent::__construct();
    }

    public function index(Request $request)
    {
        $data = [
            'id'                  => '',
            'domain'              => '',
            'corp_id'             => '',
            'created_at'          => '',
            'status'              => '',
            'external_secret'     => '',
            'mini_program_secret' => '',
            'pl_plugin_id'        => '',
            'agent_id'            => '',
            'app_id'              => '',
        ];

        $params = $request->only(['enterprise_id', 'type']);

        if (empty($params['enterprise_id']) || !is_int($params['enterprise_id']) || $params['enterprise_id'] < 1) {
            return notice(201104, [], $params);
        }

        if (empty($params['type']) || !is_int($params['type'])
            || $params['type'] < 1
            || !in_array($params['type'],
                [EnterpriseThirdPartAccount::TYPE_QY_WX, EnterpriseThirdPartAccount::TYPE_SHOP])) {
            return notice(910005, [], $params);
        }

        $enterprise = Enterprise::query()
            ->where('id', $params['enterprise_id'])
            ->first(['id', 'domain', 'corp_id']);

        if ($enterprise) {
            $enterprise      = $enterprise->toArray();
            $data['id']      = $enterprise['id'];
            $data['domain']  = $enterprise['domain'];
            $data['corp_id'] = $enterprise['corp_id'];
        }

        // 查询授权详情
        $authInfo = EnterpriseThirdPartAccount::query()
            ->where([
                'enterprise_id'  => $params['enterprise_id'],
                'authorize_type' => EnterpriseThirdPartAccount::AUTHORIZE_TYPE_THIRD, // 第三方授权
                'status'         => EnterpriseThirdPartAccount::STATUS_ACTIVE, // 有效状态
                'platform'       => EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT,// 小程序
                'type'           => $params['type'] // 小程序类型
            ])->first();

        if ($authInfo) {
            $authInfo           = $authInfo->toArray();
            $data['created_at'] = date('Y-m-d H:i:s', $authInfo['created_at']);
            $data['status']     = $authInfo['status'];
            $data['app_id']     = $authInfo['app_id'];
            $data['app_secret']     = $authInfo['app_secret'];
        }

        if (EnterpriseThirdPartAccount::TYPE_SHOP == $params['type']) {
            $shopSetting = EnterpriseWorkPlugin::query()
                ->where('enterprise_id', $params['enterprise_id'])
                ->where('identifier', 'crs-shop')
                ->first();

            if ($shopSetting) {
                $shopSetting                 = $shopSetting->toArray();
                $data['mini_program_secret'] = $shopSetting['secret'];
                $data['pl_plugin_id']        = $shopSetting['pl_plugin_id'];
                $data['corp_id']             = $shopSetting['corp_id'];
                $data['agent_id']            = $shopSetting['agent_id'];
            }
        }

        if (EnterpriseThirdPartAccount::TYPE_QY_WX == $params['type']) {
            $workSetting = EnterpriseWorkSetting::query()
                ->where('enterprise_id', $params['enterprise_id'])
                ->first(['external_secret', 'mini_program_secret', 'pl_plugin_id']);

            if ($workSetting) {
                $workSetting                 = $workSetting->toArray();
                $data['mini_program_secret'] = $workSetting['mini_program_secret'];
                $data['pl_plugin_id']        = $workSetting['pl_plugin_id'];
                $data['external_secret']     = $workSetting['external_secret'];
            }
        }
        Log::info("授权" . json_encode($data));
        return success(['data' => $data]);
    }
}
