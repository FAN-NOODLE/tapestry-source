<?php

namespace App\Http\Controllers\V1\MaterialCenter;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\MaterialCenter;

class DownLoadFileController extends BaseController
{
    /**
     * 下载素材
     * @return mixed
     * @throws \Exception
     */
    public function downloadFile()
    {
        $url = $this->request->input('file_url');

        // 文件名
        $filename = basename($url);

        // 文件实体
        $fileBody = file_get_contents($url);

        // 响应字节流
        return response()->make($fileBody, 200, [
            'Content-Type'              => 'application/octet-stream',
            'Content-Disposition'       => 'attachment; filename="' . $filename . '"',
            'Content-Transfer-Encoding' => 'binary'
        ]);

//        $params                  = $this->request->all();
//        $params['enterprise_id'] = $this->enterpriseId;
//        $this->writeOperatorLog('下载文件');
//        return (new MaterialCenter)->downloadFile($params);
    }
}