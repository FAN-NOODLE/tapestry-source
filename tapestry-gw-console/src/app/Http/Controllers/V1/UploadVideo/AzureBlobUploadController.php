<?php
/**
 * Created by PhpStorm.
 * User: chengyujia
 * Date: 2019/09/18
 * Time: 11:04
 */

namespace App\Http\Controllers\V1\UploadVideo;

use App\Http\Controllers\V1\BaseController;
use App\Services\AzureBlobService;

class AzureBlobUploadController extends BaseController
{

    protected $requireLogin =false;
    /**
     * 获取微软云blob上传参数
     */
    public function getSas()
    {
        $str = (new AzureBlobService())->getConnectionStr();
        return success(['data' => ['sas' => $str]]);
    }
}
