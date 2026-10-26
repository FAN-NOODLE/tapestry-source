<?php

namespace App\Http\Controllers\Service\V1\Auth;

use App\Common\RedisKey;
use App\Exceptions\Error;
use App\Library\Uc\EnterprisePlugin;
use App\Models\Enterprise;
use App\Models\EnterpriseWorkPlugin;
use GuzzleHttp\Client;
use Illuminate\Http\Request;
use App\Http\Controllers\Controller;
use Illuminate\Support\Facades\Log;

class InstallAppController extends Controller
{
    /*
     * 调用地址
     */
    const INSTALL_APP_URL = "%s/s/enplugin/one-to-many-Install";

    const UPDATE_APP_URL = "%s/b/%s/enplugin/pluginUpdate";

    // 企业微信红包H5地址
    const WORK_WE_CHART_RED_PACKET_H5_UTL = '%s/%s/crs-manu-customer/red-envelope';

    // 安装企业微信小程序地址
    const WORK_WE_CHART_MINI_APP_INSTALL_URL = '%s/s/suite/install?suiteId=%s&callbackUrl=%s';

    /*
     * 必填参数
     */
    protected $requestParams = [
        'installType'       => 1,
        'enterpriseChannel' => 2,
        'plName'            => '客户管理',
        'plIdentifier'      => 'crs-manu-customer',
        'imgName'           => 'icon-app-customer',
        'token'             => 'REDACTED_FOR_SOURCE_SHARING',
        'encodingAeskey'    => 'REDACTED_FOR_SOURCE_SHARING'
    ];

    // 手动安装应用,目前仅用于侧边栏
    public function install(Request $request, Client $client)
    {
        $params = $request->only(['enterprise_id', 'agent_id', 'secret']);
        if (empty($params['enterprise_id']) || empty($params['agent_id']) || empty($params['secret'])) {
            throw new Error(1400, '参数不完整!');
        }
        array_walk($params, function (&$value) {
            $value = trim($value);
        });

        $enterpriseInfo = Enterprise::query()
            ->where('id', $params['enterprise_id'])
            ->first();
        if (!$enterpriseInfo) {
            throw new Error(1401, '企业信息不存在');
        }
        if (!$enterpriseInfo->corp_id) {
            throw new Error(1402, '企业微信未安装');
        }

        $params['plHost']  = config('frontend');
        $params['corpId']  = $enterpriseInfo->corp_id;
        $params['agentId'] = $params['agent_id'];
        $rowData           = array_merge($this->requestParams, $params);
        unset($rowData['enterprise_id']);
        unset($rowData['agent_id']);

        try {
            $host        = config('client.api_url.uc_api');
            $installHost = sprintf(self::INSTALL_APP_URL, $host);
            $response1   = $client->post($installHost, [
                'body'    => json_encode($rowData),
                'headers' => [
                    'Content-Type' => 'application/json;charset=UTF-8',
                ]
            ]);
            // $res = '{"code":"SUCCESS","msg":"成功","timestamp":1539945359709,"requestId":"8BE5688A7F00000134054043DA9C58F9","data":{"corpId":"wwe19af9e9307e1003","plIdentifier":"car-manu-customer","epId":"850711C27F00000118E6BD970F4B8758","epEnumber":"850711C27F00000118E6BD970F4B8758","plPluginid":"2DC8FAC17F0000014B7C213E29500FF3"}}';
            $response1Data = json_decode($response1->getBody()->getContents(), true);
            Log::info("【安装自建应用：返回结果001】", $response1Data);
            if ($response1Data['code'] != 'SUCCESS') {
                throw new Error(1403, $response1Data['msg']);
            }

            $where = [
                'enterprise_id' => $params['enterprise_id'],
                'identifier'    => $rowData['plIdentifier']
            ];

            $update = [
                'enterprise_id' => $params['enterprise_id'],
                'agent_id'      => $params['agent_id'],
                'secret'        => $params['secret'],
                'corp_id'       => $rowData['corpId'],
                'identifier'    => $rowData['plIdentifier'],
                'pl_plugin_id'  => $response1Data['data']['plPluginid'] ?? '',
                'check_file'    => $request->input('check_file', 2),  //默认未上传域名校验文件
                'file_url'      => $request->input('file_url', '')  //证书文件地址
            ];

            $query = EnterpriseWorkPlugin::query()->where($where)->first();

            if (!empty($query)) {
                $query->update($update);
                delRedis([RedisKey::ENTERPRISE_WORK_PLUGIN_BY_ENTERPRISE_ID . $params['enterprise_id']]);
            } else {
                EnterpriseWorkPlugin::query()->create($update);
            }

        } catch (\Exception $e) {
            Log::info("【安装自建应用：失败返回结果】", [$e->getMessage()]);
            throw new Error(1404, '提交失败');
        }

        return success([]);
    }

    /**
     *
     * 获取自建应用信息
     */
    public function getWorkPlugin()
    {
        $params = request()->only(['enterprise_id']);

        $data = EnterpriseWorkPlugin::query()
            ->where([
                'enterprise_id' => $params['enterprise_id'],
                'identifier'    => 'crs-manu-customer'
            ])->first();

        if (!$data) {
            $data   = [];
            $corpId = Enterprise::query()->where([
                'id' => $params['enterprise_id'],
            ])->value('corp_id');
        } else {
            $data   = $data->toArray();
            $corpId = $data['corp_id'];
        }
        $data['domain'] = config('frontend');
        $data['url']    = config('frontend') . '/' . $corpId . '/crs-manu-customer/index';

        // 小红包
        $data['red_packets_url'] = config('frontend') . '/' . $corpId . '/crs-manu-customer/red-envelope';
        // 客户资料
        $data['customer_information_url'] = config('frontend') . '/' . $corpId . '/crs-manu-customer/index';
        // 口令券
        $data['password_coupon_url'] = '';
        // 素材库
        $data['material_library_url'] = config('frontend') . '/' . $corpId . '/crs-manu-customer/material-library';
        // 群分析
        $data['group_analysis_url'] = config('frontend') . '/' . $corpId . '/crs-manu-customer/group-analysis';
        // 营销推广
        $data['marketing_promotion_url'] = config('frontend') . '/' . $corpId . '/crs-manu-customer/marketing-promotion';

        $data['domain']         = config('frontend');
        $data['url']            = config('frontend') . '/' . $corpId . '/crs-manu-customer/index';
        $data['red_packet_url'] = vsprintf(self::WORK_WE_CHART_RED_PACKET_H5_UTL, [config('frontend'), $corpId]);

        return success($data);
    }

    // 手动更新应用,目前仅用于侧边栏
    public function update(Request $request, Client $client)
    {
        $params = $request->only(['enterprise_id', 'agent_id', 'secret']);
        if (empty($params['enterprise_id']) || empty($params['agent_id']) || empty($params['secret'])) {
            throw new Error(1400, '参数不完整!');
        }
        array_walk($params, function (&$value) {
            $value = trim($value);
        });

        $enterpriseInfo = Enterprise::query()
            ->where('id', $params['enterprise_id'])
            ->first();
        if (!$enterpriseInfo) {
            throw new Error(1401, '企业信息不存在');
        }
        if (!$enterpriseInfo->corp_id) {
            throw new Error(1402, '企业微信未安装');
        }

        $pluginInfo = EnterpriseWorkPlugin::query()
            ->where([
                'enterprise_id' => $params['enterprise_id'],
                'identifier'    => $this->requestParams['plIdentifier']
            ])
            ->first();

        if (!$pluginInfo) {
            throw new Error(1403, '企业微信未安装');
        }

        $rowData['agentId']      = $params['agent_id'];
        $rowData['corpId']       = $pluginInfo->corp_id ?? '';
        $rowData['plIdentifier'] = $pluginInfo->identifier ?? '';
        $rowData['plPluginid']   = $pluginInfo->pl_plugin_id;
        $rowData['secret']       = $params['secret'];

        try {
            $host        = config('client.api_url.uc_api');
            $installHost = sprintf(self::UPDATE_APP_URL, $host, $enterpriseInfo->domain);
            $response1   = $client->post($installHost, [
                'body'    => json_encode($rowData),
                'headers' => [
                    'Content-Type' => 'application/json;charset=UTF-8',
                ]
            ]);
            // $res = '{"code":"SUCCESS","msg":"成功","timestamp":1539945359709,"requestId":"8BE5688A7F00000134054043DA9C58F9","data":{"corpId":"wwe19af9e9307e1003","plIdentifier":"car-manu-customer","epId":"850711C27F00000118E6BD970F4B8758","epEnumber":"850711C27F00000118E6BD970F4B8758","plPluginid":"2DC8FAC17F0000014B7C213E29500FF3"}}';
            $response1Data = json_decode($response1->getBody()->getContents(), true);
            Log::info("【更新自建应用：返回结果001】", $response1Data);
            if ($response1Data['code'] != 'SUCCESS') {
                throw new Error(1404, $response1Data['msg']);
            }

            EnterpriseWorkPlugin::query()
                ->where([
                    'enterprise_id' => $params['enterprise_id'],
                    'identifier'    => $this->requestParams['plIdentifier']
                ])
                ->update([
                    'enterprise_id' => $params['enterprise_id'],
                    'agent_id'      => $params['agent_id'],
                    'secret'        => $params['secret'],
                    'corp_id'       => $rowData['corpId'],
                    'identifier'    => $rowData['plIdentifier'],
                    'pl_plugin_id'  => $response1Data['data']['plPluginid'] ?? '',
                    'check_file'    => $request->input('check_file', 2),  //默认未上传域名校验文件
                    'file_url'      => $request->input('file_url', '')  //默认未上传域名校验文件
                ]);
            delRedis([RedisKey::ENTERPRISE_WORK_PLUGIN_BY_ENTERPRISE_ID . $params['enterprise_id']]);
        } catch (\Exception $e) {
            Log::info("【更新自建应用：失败返回结果】", [$e->getMessage()]);
            throw new Error(1404, '提交失败');
        }

        return success([]);
    }

    /**
     * 获取智慧导购标识
     * @param Request $request
     * @return array|void
     */
    public function daoGouPlugin(Request $request)
    {
        $corpId = $request->input('corp_id');

        if (empty($corpId)) {
            return error(201105);
        }

        $plugins = EnterpriseWorkPlugin::query()
            ->select(['enterprise_id', 'identifier'])
            ->where('corp_id', '=', $corpId)
            ->whereIn('identifier', [EnterpriseWorkPlugin::CRS_ASSISTANT, EnterpriseWorkPlugin::CRS_APP_ASSISTANT])
            ->first();

        if (empty($plugins)) {
            return error(201106);
        }

        $enterprise = Enterprise::query()
            ->where('id', '=', $plugins->enterprise_id)
            ->first();

        if (empty($enterprise)) {
            return error(201108);
        }

        return success([
            'data' => [
                'hash'          => $enterprise->hash,
                'domain'        => $enterprise->domain,
                'identifier'    => $plugins->identifier,
                'enterprise_id' => $plugins->enterprise_id,
            ]
        ]);
    }

    /**
     * 安装自建导购
     * @param Request $request
     * @return array|void
     */
    public function daoGouInstall(Request $request)
    {
        $enterpriseId = $request->input('enterprise_id');
        if (empty($enterpriseId)) {
            return error(201101);
        }

        $enterprise = Enterprise::query()
            ->where('id', '=', $enterpriseId)
            ->first();
        if (empty($enterprise)) {
            return error(201108);
        }
        if (empty($enterprise->corp_id)) {
            return error(201109);
        }

        $plugins = EnterpriseWorkPlugin::query()
            ->where('corp_id', '=', $enterprise->corp_id)
            ->whereIn('identifier', [EnterpriseWorkPlugin::CRS_ASSISTANT, EnterpriseWorkPlugin::CRS_APP_ASSISTANT])
            ->first();
        if ($plugins) {
            $plugins->update([
                'agent_id'   => $request->input('agent_id', ''),
                'secret'     => $request->input('mini_program_secret', ''),
                'identifier' => EnterpriseWorkPlugin::CRS_ASSISTANT,
            ]);
            delRedis([RedisKey::ENTERPRISE_WORK_PLUGIN_BY_ENTERPRISE_ID . $enterpriseId]);
        } else {
            EnterpriseWorkPlugin::query()->create([
                'enterprise_id' => $enterpriseId,
                'corp_id'       => $enterprise->corp_id,
                'agent_id'      => $request->input('agent_id'),
                'secret'        => $request->input('mini_program_secret'),
                'identifier'    => EnterpriseWorkPlugin::CRS_ASSISTANT,
                'pl_plugin_id'  => $request->input('pl_plugin_id'),
            ]);
        }

        return success();
    }

    /**
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function thirdInstall(Request $request)
    {
        $identifier   = $request->input('identifier', EnterpriseWorkPlugin::CRS_APP_ASSISTANT);
        $enterpriseId = $request->input('enterprise_id', 0);
        $wisdomSet    = $request->input('wisdom_set', 0);

        $env        = config('app.env');
        $workPlugin = config('work_plugin.' . $identifier);
        $ucApi      = config('client.api_url.uc_api');
        $workPlugin = $workPlugin[$env];

        // default 返回
        $data = [
            'can_install' => 1, // 默认可以安装
            'logo'        => '',
            'app_name'    => '',
            'auth_url'    => '',
            'installed'   => 2,
        ];

        if (empty($workPlugin)) {
            $data['can_install'] = 4; // 无此应用
        }

        $enterprise = Enterprise::query()->where('id', $enterpriseId)->first();

        if (empty($enterprise)) {
            $data['can_install'] = 2; // 无此企业
        }

        if ($enterprise) {
            $availableSystem = json_decode($enterprise['available_system'], 1);

            // 业务来源; 0:-,1:企业微信,10:畅移（直签）,11:畅移（代理）,12:畅移（试用）20:畅鲟（直签）21:畅鲟（代理）22:畅鲟（试用）
            if (1 <> $enterprise['business_source'] && (!in_array(2, $availableSystem ?? [])) && $wisdomSet == 0) {
                $data['can_install'] = 3; // 未勾选智慧导购
            }

            $plugins = EnterpriseWorkPlugin::query()
                ->where('enterprise_id', '=', $enterpriseId)
                ->whereIn('identifier', [EnterpriseWorkPlugin::CRS_ASSISTANT, EnterpriseWorkPlugin::CRS_APP_ASSISTANT])
                ->first();

            // 如果已经有自建智慧导购安装, 则不可安装
            if ($plugins && $plugins->identifier === EnterpriseWorkPlugin::CRS_ASSISTANT) {
                $data['can_install'] = 5; // 已安装自建智慧导购，不可安装第三方智慧导购
            }
        }

        if (1 == $data['can_install']) {
            $backUrl = $request->input('back_url', '');

            $args = [
                $ucApi,
                $workPlugin['suite_id'],
                $backUrl
            ];

            $installUrl = vsprintf(self::WORK_WE_CHART_MINI_APP_INSTALL_URL, $args);

            $data['logo']     = $workPlugin['app_logo'];
            $data['app_name'] = $workPlugin['app_name'];
            $data['auth_url'] = $installUrl;


            try {
                // channelType 1=99HR 2=员圈 3=车店 5=CRS
                $enterprisePluginList = (new EnterprisePlugin())->list($enterprise['domain'], [
                    'channelType'  => 5,
                    'plIdentifier' => $identifier
                ]);
            } catch (\Exception $exception) {
                $enterprisePluginList = [];
            }

            // 判断是否可以安装
            $data['installed'] = 2;

            $idArr = [$identifier];
            if ($identifier === EnterpriseWorkPlugin::CRS_APP_ASSISTANT) {
                $idArr[] = EnterpriseWorkPlugin::I_CLICK_CRS_APP_ASSISTANT;
            }

            foreach ($enterprisePluginList['data'] as $datum) {
                // uc这里坑available可能返回null
                if (in_array($datum['plIdentifier'], $idArr) && in_array($datum['available'], [0, 1], true)) {
                    $data['installed'] = 1;
                    break;
                }
            }

            $enterpriseCurrentPlugin = collect($enterprisePluginList['data'] ?? [])
                ->where('plIdentifier', $identifier)
                ->whereIn('available', [0, 1]) // available 0,1都是正常，2是停用，null是未安装
                ->all();
            if ($enterpriseCurrentPlugin) {
                $workPlugin = EnterpriseWorkPlugin::query()->where(['enterprise_id' => $enterpriseId, 'identifier' => $identifier])->first();
                if (empty($workPlugin)) {
                    // 应用商店智慧导购
                    EnterpriseWorkPlugin::query()->create([
                        'enterprise_id' => $enterpriseId,
                        'corp_id'       => $enterprise['corp_id'] ?? '',
                        'identifier'    => $identifier
                    ]);
                }
            }
        }


        return $data;
    }


    /**
     * 获取导购配置
     * @return array|void
     */
    public function getDaoGouConfig()
    {
        $enterpriseId = $this->request->input('enterprise_id');
        if (empty($enterpriseId)) {
            return error(201101);
        }

        $plugins = EnterpriseWorkPlugin::query()
            ->where('enterprise_id', '=', $enterpriseId)
            ->whereIn('identifier', [EnterpriseWorkPlugin::CRS_ASSISTANT, EnterpriseWorkPlugin::CRS_APP_ASSISTANT])
            ->first();

        $canConfig = 1;
        // 如果是应用商店安装, 则不可安装
        if ($plugins && $plugins->identifier === EnterpriseWorkPlugin::CRS_APP_ASSISTANT) {
            $canConfig = 2;
        }

        $plugins = empty($plugins) ? [] : $plugins->toArray();

        return success([
            'data' => array_merge([
                'can_config' => $canConfig
            ], $plugins)
        ]);
    }
}
