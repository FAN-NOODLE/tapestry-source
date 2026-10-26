<?php

namespace App\Http\Controllers\Service\V2\Auth;

use App\Common\RedisKey;
use App\Exceptions\Error;
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

        return successV2();
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

        return successV2($data);
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

        return successV2();
    }
}
