<?php
/**
 * Created by PhpStorm.
 * User: 王鹏飞
 * Date: 2019-07-11
 * Time: 10:20
 */

namespace App\Services;

use App\Exceptions\Error;
use App\Library\EasyWechatCacheRedis\MiniProgram;

class SubscribeMsgService
{

    /**
     * @param $enterpriseId
     * @return array
     * @throws \Exception
     */
    public function getMiniInfo($enterpriseId): array
    {
        $list = EnterpriseService::getEnterpriseIdentifierList($enterpriseId);

        if (is_array($list)) {
            foreach ($list as $val) {
                if ($val['identifier'] == 'crs_store_assistant_self') {
                    $appId = $val['mini_app_id'];
                    $appSecret = $val['mini_app_secret'];
                    break;
                }
            }
        }
        if (empty($appId) || empty($appSecret)) {
            throw new Error(90001, '小程序配置错误');
        }
        return [$appId, $appSecret];
    }


    public function send($appId, $appSecret, $params)
    {

        $data = [
            'touser' => $params['touser'],
            'template_id' => $params['template_id'],
            'page' => $params['page'],
            'data' => $params['data'],
        ];
        return  MiniProgram::init(['app_id' => $appId, 'secret' => $appSecret])
            ->subscribe_message->send($data);
    }


    /**
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function getTplInfo($appId, $appSecret, $tplId)
    {

        return  MiniProgram::init(['app_id' => $appId, 'secret' => $appSecret])
            ->subscribe_message->getTemplateKeywords($tplId);
    }


}
