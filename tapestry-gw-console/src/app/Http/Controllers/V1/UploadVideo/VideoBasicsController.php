<?php
/**
 * 初始化文件上传基础数据
 */

namespace App\Http\Controllers\V1\UploadVideo;

use Illuminate\Http\Request;
use CrsSdk\Shop\UploadVideo;
use App\Http\Controllers\V1\BaseController;

class VideoBasicsController extends BaseController
{

    /**
     * @param Request $request
     * @return mixed
     * @throws \Exception
     */
    public function index(Request $request)
    {
        // 接受参数
        $input                  = $request->input();
        $input['enterprise_id'] = $this->enterpriseId;

        return (new UploadVideo())->basics($input);
    }
}
