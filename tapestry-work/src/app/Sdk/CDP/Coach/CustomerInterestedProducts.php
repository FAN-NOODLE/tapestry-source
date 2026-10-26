<?php

namespace App\Sdk\CDP\Coach;

use App\Library\WechatNotice\Chat;
use Exception;

class CustomerInterestedProducts
{

    /**
     * 获取EAS推荐商品
     *
     * @param $params
     * @return array
     * @throws Exception
     */
    public function get($params) :array
    {
        $config = Config::getEasConfig();
        $apiUrl = $config['url'];
        $payload = RecommendationFeed::buildPayload($params);
        $res = Request::Instance()->sendEasPost($apiUrl, $payload, $config);
        if (!in_array($res['code'] ?? 0, [200, 299], true)) {
            Chat::getInstance()->send('EAS推荐接口返回异常', [
                'env' => config('app.env'),
                'api' => $apiUrl,
                'param' => $payload,
                'error_info' => json_encode($res),
            ]);
        }
        return RecommendationFeed::normalizeResponse($res);
    }
}
