<?php

use Laravel\Lumen\Testing\DatabaseTransactions;

/**
 * Class BaseTest
 * @test
 */
class BaseTest extends TestCase
{
    use DatabaseTransactions;

    /**
     * 基础url
     */
    const BASE_URL = 'http://www.work.com/app/v1';

    const MOCK_HASH = '054edce1d4cec91c76001d968841c08f';
    /**
     * 当前登陆的token
     *
     * @var null
     */
    private $token = null;

    /**
     * 登陆接口
     */
    public function login()
    {
        if (!$this->token) {
            $response = $this->call('POST', self::BASE_URL.'/member/mock', [
                'memUid' => 'F89CC9427F00000136E287B2D710455F',
                'corpId' => 'ww881249f03e5b024f',
                'qyDomain' => '72D5430F7F000001381C7A561A35C79D',
                'userId' => '2'
            ]);
            $body = json_decode($response->content(), true);
            if (!isset($body['data']['token'])) {
                return false;
            }
            $token = $body['data']['token'];
            $this->token = $token;
        }

        return true;
    }
    /**
     * 封装接口调用
     *
     * @param $method
     * @param $uri
     * @param array $params
     * @return BaseTest
     */
    public function callApi($method, $uri, $params = [])
    {
        $this->login();
        if ($method == 'GET') {
            return $this->get($uri, $this->buildHeader());
        }else {
            return $this->post($uri, $params, $this->buildHeader());
        }
    }

    /**
     * 构造请求头
     *
     * @param array $params
     * @return array
     */
    public function buildHeader($params = [])
    {
        $params = array_merge($params, ['Enterprise-Hash'=> self::MOCK_HASH,'Token'=> $this->token]);
        return $params;
    }

    public function testData()
    {
        $this->assertTrue(true);
    }
}
