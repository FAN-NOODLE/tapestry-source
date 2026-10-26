<?php
/*
* +------------------------------------+
* | Created: phpstorm                  
* +------------------------------------+
* | Author : wyg                       
* +------------------------------------+
* | Date   : 2019/12/31                   
* +------------------------------------+
* | Time   : 16:20                   
* +------------------------------------+
* | Email  : <wangyongguo@vchangyi.com>
* +------------------------------------+
*/


namespace App\Http\Controllers\Service\V1\Common;


use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use Illuminate\Http\Request;

class UploadController extends BaseController
{
    /**
     * 上传文件
     *
     * @param Request $request
     *
     * @return array
     */
    public function upload(Request $request)
    {
        if (!$request->hasFile('file')) {
            throw new Error(56000, '上传文件不存在');
        }

        $file = $request->file('file');

        $size = $file->getSize();

        $maxSize = 2 * 1024 * 1024;

        if ($size > $maxSize) {
            throw new Error(57000, '大小超过限制');
        }

        $path = storage_path('keys/');

        if (!is_dir($path)) {
            mkdir($path, 0777, true);
        }

        @chown($path, 'cy-crs');

        $fileName = $file->getClientOriginalName();

        $extension = $file->getClientOriginalExtension();

        if ($extension != 'txt') {
            throw new Error(58000, '请上传txt文件');
        }

        $file->move($path, $fileName);

        return success(['data' => []]);
    }
}