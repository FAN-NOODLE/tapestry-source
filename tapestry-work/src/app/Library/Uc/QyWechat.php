<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/10/12
 * Time: 10:29 AM
 */

namespace App\Library\Uc;

use App\Common\RedisKey;
use App\Exceptions\Error;
use App\Services\EnterpriseService;

class QyWechat extends Abstracts
{

    private $ucToken;

    /**
     * @throws \Exception
     */
    public function __construct()
    {
        parent::__construct();
        $this->ucToken = $this->getUcToken();
    }

    /**
     * @throws \Exception
     */
    private function getUcToken()
    {
        $lockKey = RedisKey::UC_API_ACCESS_TOKEN;
        // 查询缓存
        $token = getOrSetData($lockKey);
        if (!$token) {
            $url = $this->getUcApiUrl() . '/s/token';
            $res = $this->qyPost($url, [
                    'id' => env('UC_API_ID', 'cy-scrm'),
                    'secret' => env('UC_API_SECRET', 'zEgxo4Wfyc7K6UW3')
                ]);
            if (!isset($res['data']) || !isset($res['data']['access_token']) || empty($res['data']['access_token'])) {
                throw new Error(100000, '获取UC access_token失败');
            }
            $token = $res['data']['access_token'];
            getOrSetData($lockKey, $token, 3600);
        }

        return $token;

    }

    /**
     * @param array $data
     * @param array $otherParams = ['path' => 'cgi-bin/department/get',
     *                      'method' => 'GET',
     *                      'enterprise_id' => 3363603667487424,
     *                      'domain' => '65C03781AC13030F47308B810F788FB0']
     * $otherParams
     * @return mixed
     * @throws \Exception
     */
    public function getQyResult(array $data, array $otherParams = [])
    {
        if (empty($otherParams['path']) || (empty($otherParams['domain']) && empty($otherParams['enterprise_id']))) {
            throw new Error(100000, '参数错误');
        }
        $domain = $otherParams['domain'] ?? (new EnterpriseService)->detail($otherParams['enterprise_id'])['data']['domain'];
        $params = [
            'reqPL' => $otherParams['reqPL'] ?? 'ExternalContact',
            'reqMethod' => $otherParams['method'] ?? 'POST',
            'reqUrl' => $otherParams['path'],
            'token' => $this->ucToken,
        ];
        $ucUrl = sprintf('%s/b/%s/qy/gw', $this->getUcApiUrl(), $domain);
        return $this->qyPost($ucUrl, array_merge($params, $data));


    }

}
