<?php
/**
 * Created by PhpStorm.
 * Function: 下载文件公共类
 * User: 龚国新
 * Date: 2019/10/22
 * Time: 18:14
 */

namespace App\Http\Controllers\V1\Common;


use App\Http\Controllers\Controller;

class DownLoadFileController extends Controller
{
    public function index(){
        $input = $this->request->only(['filepath']);
        if(!empty($input['filepath'])){
            $filename = basename($input['filepath']);
            $file = file_get_contents($input['filepath']);

            return response()->make($file, 200, [
                'Content-Type'              => 'application/octet-stream',
                'Content-Disposition'       => 'attachment; filename="'.$filename.'"',
                'Content-Transfer-Encoding' => 'binary'
            ]);
        }else{
            notice(310001);
        }
    }
}