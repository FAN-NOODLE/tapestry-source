<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/4/9 0009
 * Time: 10:22
 */

namespace App\Library\WechatNotice;


class Request
{
    public static function requestPostJson($host, $uri, $body, $options = [])
    {
        $client = ClientFactory::create($host);
        $options['headers'] = ["Content-Type" => "application/json"];
        $options['body'] = json_encode($body);
        return $client->post($uri, $options)->getBody()->getContents();
    }
}