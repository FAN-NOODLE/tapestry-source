<?php

namespace App\Http\Controllers\Service\V1\FastRegisterWeApp;

use App\Http\Controllers\Service\V1\BaseController;
use App\Http\Requests\V1\FastRegisterWeAppRequest;
use App\Models\EnterpriseOpenPlatform;
use App\Models\RegisterProgramAgent;
use App\Services\EnterpriseLogService;
use App\Services\EnterpriseService;
use GuzzleHttp\Client;
use GuzzleHttp\Exception\GuzzleException;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use Overtrue\LaravelWeChat\Facade as LaravelWechat;

class ApplyController extends BaseController
{
    const API_URL = 'https://api.weixin.qq.com/cgi-bin/component/fastregisterweapp?action=create&component_access_token=%s';

    public function __construct()
    {
        parent::__construct();
    }

    public function index(Request $request, RegisterProgramAgent $model)
    {

        $params = $request->only([
            'id',
            'enterprise_id',
            'legal_persona_name',
            'legal_persona_wechat',
            'name',
            'code_type',
            'code',
        ]);
        FastRegisterWeAppRequest::apply($params);

        $formatData = [
            'legal_persona_name'   => $params['legal_persona_name'],
            'legal_persona_wechat' => $params['legal_persona_wechat'],
            'name'                 => $params['name'],
            'code_type'            => $params['code_type'],
            'code'                 => $params['code'],
        ];

        if (isset($params['id']) && $params['id'] != "") {
            $info = $model->getInfoById($params['id']);
            if (empty($info)) {
                notice(100000, ['ID有误']);
            }

            if ($info['status'] == $model::STATUS_APPLY) {
                notice(100000, ['正在审核中，请勿重复申请']);
            } elseif ($info['status'] == $model::STATUS_APPLY_SUCCESS) {
                notice(100000, ['审核成功，请勿重复申请']);
            }
        }

        $enterprise   = (new EnterpriseService())->enterprise($params['enterprise_id']);
        $openPlatform = EnterpriseOpenPlatform::query()->where('id', '=', $enterprise->open_platform_id)->first();
        if (empty($openPlatform)) {
            error(620004);
        }
        // token的处理
        $open                    = LaravelWechat::openPlatform('default');
        $open->config['token']   = $openPlatform->token;
        $open->config['app_id']  = $openPlatform->app_id;
        $open->config['secret']  = $openPlatform->secret;
        $open->config['aes_key'] = $openPlatform->aes_key;
        $token                   = $open->access_token->getToken()['component_access_token'] ?? '';

        $url       = sprintf(self::API_URL, $token);
        $openParam = ["code" => str_replace('-', '', $params['code']), "component_phone" => ""];

        $result = $this->httpPost($url, $openParam + $formatData);

        if ($result['errcode'] != 0) {
            $mess = RegisterProgramAgent::REGISTER_CODE_MAP[$result['errcode']] ?? $result['errmsg'];
            Log::info('开放平台响应数据', ['url' => $url, 'result' => $result]);
            notice(100000, [$mess]);
        }

        try {
            DB::beginTransaction();
            if (!empty($info)) {
                $formatData['status'] = $model::STATUS_APPLY;
                $model->updateData($info['id'], $formatData);
            } else {
                $formatData['enterprise_id'] = $params['enterprise_id'];
                $model->insertData($formatData);
            }
            DB::commit();
        } catch (\Exception $e) {
            Log::info($params['enterprise_id'] . ' 更新小程序信息失败 ' . $e->getMessage());
            DB::rollBack();
            notice(100000, ['更新失败']);
        }

        return success();
    }

    /**
     * @param $url
     * @param array $data
     * @return mixed|string
     * @throws GuzzleException
     */
    private function httpPost($url, $data = [])
    {
        $method = 'POST';
        $apiUrl = $url;
        $client = new Client();
        try {
            $response = $client->request($method, $apiUrl, [
                'headers'     => [
                    'Accept'       => 'application/json',
                    'Content-Type' => 'application/json'
                ],
                'body'        => json_encode($data, JSON_UNESCAPED_UNICODE),
                'http_errors' => false
            ]);

            $result = $response->getBody();
            $result = json_decode($result, true);

        } catch (\Exception $e) {
            Log::debug('openMiniPostRequest:Error ' . $e->getMessage());
            EnterpriseLogService::error('OpenPlatformError', [
                'url'     => $url,
                'data'    => $data,
                'message' => $e->getMessage(),
                'file'    => $e->getFile(),
                'line'    => $e->getLine()
            ]);
            return $e->getMessage();
        }
        return $result;
    }
}
