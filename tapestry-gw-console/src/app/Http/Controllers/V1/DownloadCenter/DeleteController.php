<?php


namespace App\Http\Controllers\V1\DownloadCenter;


use App\Http\Controllers\V1\BaseController;
use App\Micros\DownloadCenter\DownloadCenterMicro;

class DeleteController extends BaseController
{
    public function index()
    {
        // 接值
        $input = $this->request->input();
        $input['enterprise_id'] = $this->enterpriseId;
        // 数据验证
        $this->validate($this->request, [
            'id' => 'required'
        ],['id.required' => '下载任务ID不能为空']);

        (new DownloadCenterMicro())->delete($input);

        return success();
    }
}
