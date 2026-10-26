<?php
/**
 * Created by PhpStorm.
 * User: chengyujia
 * Date: 2019/09/18
 * Time: 11:04
 */

namespace App\Http\Controllers\V1\UploadVideo;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Common\QCloudCos;

class SignatureController extends BaseController
{

    /**
     * @return array
     */
    public function index()
    {
        $sign = (new QCloudCos())->getSign();
        return success(['data' => ['sign' => $sign]]);
    }
}
