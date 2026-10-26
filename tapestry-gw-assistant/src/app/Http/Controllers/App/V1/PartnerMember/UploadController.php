<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/11/16 17:46
 */

namespace App\Http\Controllers\App\V1\PartnerMember;


use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\Common\QCloudCos;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class UploadController extends BaseController
{
    protected $isRequireLogin = false;

    public function index(Request $request, QCloudCos $QCloudCos)
    {
        Log::info('上传图片', [json_encode($request->file)]);
        $this->validate($request, [
            'file' => 'required|mimes:jpeg,jpg,png,gif,webp'
        ], [
            'file.required' => '必须上传图片',
            'file.mimes'    => '图片格式不正确,只支持jpeg、jpg、png、gif和webp图片'
        ]);
        $hash = $request->header('Enterprise-Hash','');
        $url = $QCloudCos->upload($request->file,$hash);

        return success(['url' => $url]);
    }
}