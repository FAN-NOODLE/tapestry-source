<?php
/**
 * 生成体验二维码
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/1/11
 * Time: 18:07
 */

namespace App\Http\Controllers\Service\V1\MiniProgram;

class QrCodeController extends MiniBaseController
{
    /**
     * 小程序体验二维码
     *
     * @return array|void
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     */
    public function index()
    {
        // 小程序app_id
        $appId = $this->appId;
        if (!$appId) {
            return error(300100);
        }

        // 体验二维码
        $experienceQrcode = $this->miniProgramService->getExperienceQrcode($appId);

        return ['data' => ['experience_qrcode' => $experienceQrcode]];
    }
}
