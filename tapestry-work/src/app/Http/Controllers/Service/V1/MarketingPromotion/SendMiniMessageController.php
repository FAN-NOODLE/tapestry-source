<?php
/*
 * @Author: your name
 * @Date: 2020-02-11 12:02:35
 * @LastEditTime: 2020-02-24 20:36:32
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \shop\src\app\Jobs\PasswordCouponJob.php
 */

namespace App\Http\Controllers\Service\V1\MarketingPromotion;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\EnterprisePlugin;
use Illuminate\Support\Facades\Validator;

class SendMiniMessageController extends BaseController
{
    /**
     * 发送小程序应用消息
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        $params = request()->only(['enterprise_id', 'domain', 'identifier', 'content', 'mini_program', 'mem_uid']);

        $rules = [
            'enterprise_id' => 'required',
            'mem_uid'       => 'required',
            'domain'        => 'required',
            'identifier'    => 'required',
            'content'       => 'required',
            'mini_program'  => 'required',
        ];

        $message = [
            'enterprise_id.required' => '企业id不能为空',
            'mem_uid.required'       => '导购mem_uid不能为空',
            'domain.required'        => 'domain不能为空',
            'identifier.required'    => 'identifier不能为空',
            'content.required'       => '发送内容不能为空',
            'mini_program.required'  => '小程序信息不能为空',
        ];

        Validator::make($params, $rules, $message)->validate();

        $ucParams = [
            'minicontent' => $params['content'],
            'miniprogram' => $params['mini_program'],
            'msgType' => 'miniprogram',
            'toUser' => implode('|', $params['mem_uid'])
        ];

        (new EnterprisePlugin())->sendMiniMessage($params['domain'], $params['identifier'], $ucParams);

        return success([]);
    }
}
