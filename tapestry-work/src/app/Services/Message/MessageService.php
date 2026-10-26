<?php


namespace App\Services\Message;


use App\Library\Uc\EnterprisePlugin;
use App\Library\WechatNotice\Chat;
use App\Services\EnterpriseService;
use CrsSdk\Boss\Enterprise;
use CrsSdk\Boss\EnterpriseWorkPlugin;
use Illuminate\Support\Facades\Log;

class MessageService
{


    /**
     * 发送小程序文本消息
     * @param $enterpriseId
     * @param $content
     * @param $page
     * @param $title
     * @param $memUids
     * @param string $domain
     * @return bool
     * @throws \Exception
     */
    public function sendMiniMessage($enterpriseId, $content, $page, $title, $memUids, $domain = '')
    {
        Log::info("MessageService.sendMiniMessage params", func_get_args());
        $enterprise = (new Enterprise())->detail(['enterprise_id' => $enterpriseId])['data'];
        //单小程序给多个企业使用需要在pages参数中增加corp_id
        $pageParseArr = parse_url($page);
        parse_str($pageParseArr['query'] ?? '', $pageParams);
        if (!in_array('corp_id' , array_keys($pageParams))) {
            $pageParams['corp_id'] = $enterprise['corp_id'];
            $page = $pageParseArr['path'] . '?' .urldecode(http_build_query($pageParams));
        }
        if (empty($domain)) {
            $domain = $enterprise['domain'];
        }

        $sendMiniMessage = function($appId, $identifier) use ($domain, $page, $title, $content, $memUids) {
            if (empty($appId) || empty($identifier)) {
                Log::error('MessageService.sendMiniMessage enterprise app identifier error', [
                    'appId' => $appId,
                    'identifier' => $identifier
                ]);
                return false;
            }
            $ucParams = [
                'minicontent' => $content,
                'miniprogram' => [
                    'appid' => $appId,
                    'page'  => $page,
                    'title' => $title,
                ],
                'msgType' => 'miniprogram',
                'toUser' => implode('|', $memUids)
            ];
            Log::info("EnterprisePlugin.sendMiniMessage params", [
                'domain' => $domain,
                'identifier' => $identifier,
                'ucParams' => $ucParams
            ]);
            $result = (new EnterprisePlugin())->sendMiniMessage($domain, $identifier, $ucParams);
            Log::info("EnterprisePlugin.sendMiniMessage result", ['result' => $result]);
            return true;
        };

        $list = EnterpriseService::getEnterpriseIdentifierList($enterpriseId);

        $assistantAppId = '';
        $identifier = 'crs_store_assistant_self';
        //应用配置中获取智慧导购自建应用相关的小程序app_id
        foreach ($list as  $val) {
            if ($val['identifier'] == $identifier) {
                $assistantAppId = $val['mini_app_id'];
                break;
            }
        }
        if ($assistantAppId == '') {
            Chat::getInstance()->send('发送智慧导购应用消息,智慧导购自建应用配置有误,未配置mini_app_id');
            return false;
        } else {
            return $sendMiniMessage($assistantAppId, $identifier);
        }
    }
}