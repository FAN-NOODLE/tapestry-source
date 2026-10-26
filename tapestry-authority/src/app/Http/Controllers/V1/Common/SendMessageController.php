<?php
/**
 * Created By PhpStorm
 * @Date: 2021/4/15
 * @Time: 14:58
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: SendMessageController.php
 */

namespace App\Http\Controllers\V1\Common;

use App\Http\Controllers\V1\BaseController;
use App\Services\SendCodeService;
use App\Traits\SendMessageTrait;

class SendMessageController extends BaseController
{
    use SendMessageTrait;

    public function index()
    {
        $mobile    = $this->request->input('mobile');
        $content   = $this->request->input('content');
        $enterpriseId   = $this->request->input('enterpriseId');
        $template   = $this->request->input('template');
        $client    = $this->request->input('client', 'scrm');
        $rule      = $this->request->input('module_rule', []);
        $ip        = $this->request->getClientIp();

        if (empty($mobile)) {
            // 手机号不能为空
            notice(700031);
        }
        if (empty($client)) {
            // 客户端标识不能为空
            notice(700013);
        }

        if (empty($enterpriseId)) {
            // 企业id不能为空
            notice(700020);
        }

        if (empty($content)) {
            // 短信内容不能为空
            notice(700036);
        }

        $mobile = str_replace('+86-', '', $mobile);

        // 根据调用端传递的验证规则验证--自己写调用规则
        if (!empty($rule)) {
            (new SendCodeService())->sendCodeRule($rule, $mobile);
        }

        if ($mobile == '13543211234') {
            $fixed = 1;
        }

        $this->sendMessage($mobile, $content, $enterpriseId, $client, $template, $ip);

        return success();
    }
}