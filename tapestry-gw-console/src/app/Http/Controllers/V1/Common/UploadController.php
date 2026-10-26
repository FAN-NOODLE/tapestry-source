<?php

namespace App\Http\Controllers\V1\Common;

use CrsSdk\Work\Common;
use App\Exceptions\Error;
use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\Controller;

class UploadController extends Controller
{
    /**
     * 上传文件
     * @param Request $request
     * @return JsonResponse
     * @throws \Exception
     */
    public function upload(Request $request)
    {
        if (!$request->hasFile('file')) {
            throw new Error(51000, '上传文件不存在');
        }

        $file = $request->file('file');

        $size = $file->getSize();

        $maxSize = 2 * 1024 * 1024;

        if ($size > $maxSize) {
            throw new Error(52000, '大小超过限制[2M]');
        }

        $path = base_path('keys/');

        if (!is_dir($path)) {
            mkdir($path, 0777, true);
        }

        $originalName = $file->getClientOriginalName();

        $extension = $file->getClientOriginalExtension();

        if ($extension != 'txt') {
            throw new Error(53000, '请上传txt文件');
        }

        $file->move($path, $originalName);

        $result = (new Common())->upload(['file' => $path . $originalName]);

        if (!isset($result['data'])) {
            throw new Error(54000, '上传文件异常');
        }

        @unlink($path . $originalName);

        return success($result['data']);
    }
}