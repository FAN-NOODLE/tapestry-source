<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2021/1/13 10:39
 */

namespace App\Http\Controllers\V1\Common;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Common\QCloudCos;
use Illuminate\Http\Request;

class UploadVideoController extends BaseController
{
    protected $isRequireLogin = false;

    public function index(Request $request, QCloudCos $QCloudCos)
    {
        $type=$request->input('type','mp4');
        //mimes:mp4,mov,ogg
        $this->validate($request, [
            'file' => 'required|mimes:'.$type
        ], [
            'file.required' => '必须上传视频',
            'file.mimes'    => '视频格式不正确,只支持'.$type
        ]);

        // 大小前端来限制
        $hash = $request->header('Enterprise-Hash','');
        $url = $QCloudCos->upload($request->file,$hash);

        return success(['url' => $url]);
    }
}