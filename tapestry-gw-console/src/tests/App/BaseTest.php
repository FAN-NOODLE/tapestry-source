<?php
/*
 * Created by PhpStorm.
 * User: wangce
 * Date Time: 2019/3/6 10:23
 */
namespace Tests\App\V1;


class BaseTest extends \TestCase
{

    /**
     * 基础url
     */
    const BASE_URL = 'http://console.crs.work/app/v1';

    /**
     * Token key
     */
    const TOKEN_KEY = 'REDACTED_FOR_SOURCE_SHARING';

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
            $response = $this->call('POST', self::BASE_URL.'/admin/login', [
                'account' => '18000000000',
                'password' => 'REDACTED_FOR_SOURCE_SHARING'
            ]);
            $body = json_decode($response->content(), true);

            if (!isset($body['token'])) {
                return false;
            }
            $token = $body['token'];
            $this->token = $token;
        }
    }

    /**
     * 构造请求头
     *
     * @param $params array
     * @return array
     */
    public function buildHeader($params = [])
    {
        $params = array_merge($params, [self::TOKEN_KEY => $this->token]);
        return $params;
    }

}