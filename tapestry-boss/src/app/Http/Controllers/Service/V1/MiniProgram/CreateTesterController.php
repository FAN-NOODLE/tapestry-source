<?php
/**
 * 添加体验者
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/1/11
 * Time: 15:52
 */

namespace App\Http\Controllers\Service\V1\MiniProgram;

use Illuminate\Http\Request;

class CreateTesterController extends MiniBaseController
{
    /***
     * @param Request $request
     *
     * @return array
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     */
    public function index(Request $request)
    {
        $appId = $this->appId;
        $weChatId = $request->input('wechat_id');
        $remark = $request->input('remark');
        $name = $request->input('name');

        if (!$appId) {
            return error(201101);
        }

        if (!$weChatId) {
            return error(300003);
        }

        if (!$remark) {
            return error(300004);
        }

        if (!$name) {
            return error(300005);
        }

        $this->miniProgramService->addTester($weChatId, $appId, ['name' => $name, 'remark' => $remark]);

        return [
            'data' => [
                'message' => 'success',
            ],
        ];
    }
}
