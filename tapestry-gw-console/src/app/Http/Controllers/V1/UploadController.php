<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/9
 * Time: 10:00 AM
 */

namespace App\Http\Controllers\V1;

use App\Http\Controllers\Controller;
use CrsSdk\Common\QCloudCos;
use Illuminate\Http\Request;

class UploadController extends Controller
{
    public function index(Request $request, QCloudCos $QCloudCos)
    {
        $hash = $request->input('enterprise_hash', 'public');

        $headers    = $request->header();
        foreach ($headers as $key => $val) {
            if (strtolower($key) == 'enterprise-hash') {
                $hash = $val[0];
                break;
            }
        }

        $this->validate($request, [
            'file' => 'required|image:jpeg,jpg,png,gif'
        ], [
            'file.required' => '必须上传图片',
            'file.image'    => '图片格式不正确,只支持jpg、png和gif 图片'
        ]);

        $url = $QCloudCos->upload($request->file, $hash);

        return success(['data' => $url]);
    }
}
