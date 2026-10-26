<?php
/**
 * Created by Phpstorm.
 * User : zoulongbo
 * Date : 2019-08-01
 * Tome : 10:38
 */

namespace App\Http\Controllers\V1\Template\AutoRelease;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Boss\MiniProgram;
use Illuminate\Support\Facades\Log;

class QrCodeController extends BaseController
{
    /**
     * 自动化发布： 预览二维码
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function index()
    {
        try {
            $result = (new MiniProgram())->experienceQrCode(['enterprise_id' => $this->enterpriseId]);
        } catch (\Exception $e) {
            Log::info('自动化发布，二维码获取失败:' . $e->getMessage());
            error(920303, [$e->getMessage()]);
        }
        return success($result);
    }
}
