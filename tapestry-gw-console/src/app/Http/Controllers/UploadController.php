<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/9
 * Time: 10:00 AM
 */

namespace App\Http\Controllers;

use CrsSdk\Common\QCloudCos;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class UploadController extends Controller
{
    public function index(Request $request, QCloudCos $QCloudCos)
    {
        $hash = $request->input('enterprise_hash', 'public');

        $headers = $request->header();
        foreach ($headers as $key => $val) {
            if (strtolower($key) == 'enterprise-hash') {
                $hash = $val[0];
                break;
            }
        }

        $this->validate($request, [
            'file' => 'required|image:jpeg,jpg,png,gif,webp'
        ], [
            'file.required' => '必须上传图片',
            'file.image'    => '图片格式不正确,只支持jpg、png、gif 和 webp 图片'
        ]);

        $file = $request->file;

        $filePath = $file->path();
        $md5Uniqid = md5(uniqid());
        $picName = $hash . '/' . date('Y/m/d/', time()) . $md5Uniqid;
        $key = $picName . '.' . $file->getClientOriginalExtension();

        $url = $QCloudCos->upload($request->file, $hash, $key, $filePath);

        $type = $file->getClientOriginalExtension();

        // 上传webp格式
        if ('webp' != $type) {
            try {
                $img = new \Imagick($filePath);

                $img->setImageFormat("webp");
                if ("png" == $type) {
                    $img->setImageCompressionQuality(100);
                } else {
                    $img->setImageCompressionQuality(97);
                }

                $webpKey = $picName . '.webp';

                $webpPath = resource_path() . '/pictures/' . $md5Uniqid . '.webp';

                if (!is_dir(resource_path() . '/pictures')) {
                    mkdir(resource_path() . '/pictures');
                }

                $img->writeImage($webpPath);
                $webpUrl = $QCloudCos->upload($request->file, $hash, $webpKey, $webpPath);
                Log::info("webp格式转换成功" . $webpUrl);

                if (file_exists($webpPath)) {
                    unlink($webpPath);
                }
            } catch (\Exception $exception) {
                Log::info("webp格式转换失败" . $exception->getMessage());
            }
        }

        return success(['data' => $url]);
    }
}
