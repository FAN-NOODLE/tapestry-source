<?php
/**
 * 文件上传成功或失败
 * Created by PhpStorm.
 * User: chengyujia
 * Date: 2019/09/18
 * Time: 11:04
 */

namespace App\Http\Controllers\V1\UploadVideo;

use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;
use CrsSdk\Shop\UploadVideo;

class UploadVideoStatusController extends BaseController
{

    /**
     * @param Request $request
     *
     * @return mixed
     * @throws \Exception
     */
    public function index(Request $request)
    {
        // 接受参数
        $input = $request->input();

        return (new UploadVideo())->UploadVideoStatus($input);
    }
}
