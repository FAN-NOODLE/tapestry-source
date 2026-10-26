<?php
/**
 * 生成体验二维码
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/1/11
 * Time: 18:07
 */

namespace App\Http\Controllers\App\V1\MiniProgram;

use Illuminate\Http\Request;

class QrCodeController extends MiniBaseController
{
    /**
     * 小程序体验二维码
     *
     * @param Request $request
     *
     * @return array|void
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     */
    public function index(Request $request)
    {
        // 小程序app_id
        $appId = $request->post('app_id', '');
        if (!$appId) {
            return error(300100);
        }

        // 体验二维码
        $experienceQrcode = $this->miniProgramService->getExperienceQrcode($appId);

        return ['data' => ['experience_qrcode' => $experienceQrcode]];
    }
}
