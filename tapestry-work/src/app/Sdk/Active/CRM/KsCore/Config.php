<?php

namespace App\Sdk\Active\CRM\KsCore;


final class Config
{

    const  LOG_NAME = 'ks';
    const  NONCE = 'jsonstring';

    /**
     * storeType
     * 1 奥莱  用户绑定券---活动代码为bindWeCom-AL
     * 2 零售  用户绑定券---活动代码为bindWeCom
     */
    const CampaignCodeAL = "bindWeCom-AL";

    // 奥莱和零售不区分优惠券的领取状态，一个用户只能领取一次优惠券
    const CampaignCode = "bindWeCom";

    // uat prod
    // appid：WeCom
    // secret：u3p9g1x6uyuetd8

    // 测试账号
    static $test = [
        'appid'   => 'WeCom',
        'app_key' => 'REDACTED_FOR_SOURCE_SHARING',
    ];
    // 生产账号
    static $prod = [
        'appid'   => 'WeCom',
        'app_key' => 'REDACTED_FOR_SOURCE_SHARING',
    ];

    /**
     * ks企业
     * 测试appId：wx65d16d0e8772bac6
     * 线上appid：wx5d831b3fbcb54f75
     * 已注册用户的path路径：
     * pages/linkMemberCenter/linkMemberCenter?redirectTo=/MiniMemberCenter/Coupon
     * 新用户path路径：
     * pages/index/index?shopCode={0}&clerkCode={1}&source=WeCom
     * 错误代码errcode
     * 0表示成功
     *
     * @param bool $isNew
     *
     * @return array
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/29
     */
    static function getPathConfig(bool $isNew = false): array
    {
        if (config('app.env') == 'prod') {
            $config['app_id'] = 'wx5d831b3fbcb54f75';
        } else {
            $config['app_id'] = 'wx65d16d0e8772bac6';
        }
        if ($isNew == false) {
            // 新用户path路径：门店编号，店员编号
            // pages/index/index?shopCode={0}&clerkCode={1}&source=WeCom
            $config['path']   = 'pages/index/index?shopCode={shopCode}&clerkCode={clerkCode}&source=WeCom';
            $config['source'] = 'WeCom';

        } else {
            // 已注册用户的path路径：
            $config['path']   = 'pages/linkMemberCenter/linkMemberCenter?redirectTo=/MiniMemberCenter/Coupon';
            $config['source'] = '';

        }
        return $config;
    }


    static function getUrl(): string
    {
        if (config('app.env') == 'prod') {
//            $baseUrl = "https://wechat.katespade.aowenmarketing.com";
            $baseUrl = 'https://member-external.katespade.cn';
        } else {
//            $baseUrl = "http://twechat.katespade.wedochina.cn";
            $baseUrl = 'https://dev-member-external.katespade.cn';
        }
        return $baseUrl;
    }

    static function getAppId(): string
    {
        return Config::isProdEnv() ? Config::$prod['appid'] : Config::$test['appid'];
    }

    static function getAppKey(): string
    {
        return Config::isProdEnv() ? Config::$prod['app_key'] : Config::$test['app_key'];
    }


    static function isProdEnv(): bool
    {
        return config("app.env") == 'prod';
    }

}