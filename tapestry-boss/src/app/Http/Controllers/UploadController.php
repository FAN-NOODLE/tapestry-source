<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/9
 * Time: 10:00 AM
 */

namespace App\Http\Controllers;

use CrsSdk\Common\QCloudCos;
use CrsSdk\Common\Qiniu;
use Illuminate\Http\Request;

class UploadController extends Controller
{
    public function index(Request $request, QCloudCos $QCloudCos)
    {
        $this->validate($request, [
            'file' => 'required'
        ], [
            'file.required' => '必须上传图片',
        ]);

        $url = $QCloudCos->upload($request->file);

        return success(['data' => $url]);
    }
}
