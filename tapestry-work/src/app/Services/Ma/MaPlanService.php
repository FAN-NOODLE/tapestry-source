<?php

namespace App\Services\Ma;

use App\Library\WechatNotice\Chat;
use App\Sdk\MaPlan\Coach\Api\MomentComments;
use App\Sdk\MaPlan\Coach\Api\MomentSendResult;
use App\Sdk\MaPlan\Coach\Api\SendFinish;
use App\Sdk\MaPlan\Coach\Api\SendResult;
use GuzzleHttp\Exception\GuzzleException;

class MaPlanService
{
    protected $enterpriseId;
    protected $originType;

    public function __construct($enterpriseId,$originType)
    {
        $this->enterpriseId = $enterpriseId;
        $this->originType = $originType;
    }

    /**
     * 任务推送结果的发送
     * @param $params
     * @return array|mixed
     * @throws GuzzleException
     */
    public function sendResult($params)
    {
        $response = (new SendResult($this->enterpriseId,$this->originType))->push($params);
        if (!isset($response['success']) && !empty($response['error'])) {
            Chat::getInstance()->send('请求MA接口返回错误码', [
                'env' => config('app.env'),
                'api' => "send-result",
                'param' => $params,
                "error_code" => $response['error']['code'],
                'error_info' => $response['error']['message']
            ]);
        }
        return $response;
    }

    /**
     * 任务完成发送
     * @param $params
     * @return array|mixed
     * @throws GuzzleException
     */
    public function sendFinish($params)
    {
        $response = (new SendFinish($this->enterpriseId,$this->originType))->push($params);
        if (!isset($response['success']) && !empty($response['error'])) {
            Chat::getInstance()->send('请求MA接口返回错误码', [
                'env' => config('app.env'),
                'api' => "send-finish",
                'param' => $params,
                "error_code" => $response['error']['code'],
                'error_info' => $response['error']['message']
            ]);
        }
        return $response;
    }

    /**
     * 朋友圈任务群发助手导购发送后回执可见范围
     * @param $params
     * @return array|mixed
     * @throws GuzzleException
     */
    public function momentSendResult($params)
    {
        $response = (new MomentSendResult($this->enterpriseId,$this->originType))->push($params);
        if (!isset($response['success']) && !empty($response['error'])) {
            Chat::getInstance()->send('请求MA接口返回错误码', [
                'env' => config('app.env'),
                'api' => "MomentSendResult",
                'param' => $params,
                "error_code" => $response['error']['code'],
                'error_info' => $response['error']['message']
            ]);
        }
        return $response;
    }

    /**
     * 朋友圈任务回执评论和点赞
     * @param $params
     * @return array|mixed
     * @throws GuzzleException
     */
    public function momentComments($params)
    {
        $response = (new MomentComments($this->enterpriseId,$this->originType))->push($params);
        if (!isset($response['success']) && !empty($response['error'])) {
            Chat::getInstance()->send('请求MA接口返回错误码', [
                'env' => config('app.env'),
                'api' => "MomentComments",
                'param' => $params,
                "error_code" => $response['error']['code'],
                'error_info' => $response['error']['message']
            ]);
        }
        return $response;
    }
}