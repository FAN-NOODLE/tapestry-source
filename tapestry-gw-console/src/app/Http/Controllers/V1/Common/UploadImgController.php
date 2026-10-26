<?php

namespace App\Http\Controllers\V1\Common;

use App\Http\Controllers\V1\BaseController;
use App\Library\TencentYun\Cos\QCloudCosSdk;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Validation\ValidationException;

class UploadImgController extends BaseController
{

    /**
     * cos 上传
     *
     * @param Request      $request
     * @param QCloudCosSdk $QCloudCosSdk
     *
     * @return JsonResponse
     * @throws ValidationException
     */
    public function index(Request $request, QCloudCosSdk $QCloudCosSdk): JsonResponse
    {

        $type = $request->input('type', 'JPG, jpg, PNG, "png"');
        $this->validate($request, [
            'file' => 'required|image:' . $type
        ], [
            'file.required' => '必须上传图片',
            'file.image'    => '图片格式不正确,只支持' . $type
        ]);

        // 判断大小是否符合
        $tmpFile = $request->file->getRealPath();
        if (filesize($tmpFile) >= 10485760) {
            notice(2100001, ['10M']);
        }
        $hash = $request->header('Enterprise-Hash', '');
        $url  = $QCloudCosSdk->upload($request->file, $hash);

        return success(['url' => $url]);
    }
}