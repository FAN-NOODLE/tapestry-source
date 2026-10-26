<?php


namespace App\Http\Controllers\V1\DownloadCenter;


use App\Http\Controllers\V1\BaseController;
use App\Exceptions\Error;
use App\Micros\DownloadCenter\DownloadCenterMicro;
use Exception;
use Illuminate\Support\Facades\Log;
use Illuminate\Validation\ValidationException;
use Symfony\Component\HttpFoundation\BinaryFileResponse;

class DownloadController extends BaseController
{
    protected $requireLogin = false;

    /**
     * @return BinaryFileResponse
     * @throws ValidationException
     * @throws Exception
     */
    public function index()
    {

        // 接值
        $input = $this->request->input();
        // 数据验证
        $this->validate($this->request, [
            'id' => 'required'
        ], ['id.required' => '下载任务ID不能为空']);

        if (!empty($input['hash'])) {
            $_SERVER['enterprise_hash'] = $input['hash'];
        }

        $input['enterprise_id'] = $this->enterpriseId;
        $data                   = (new DownloadCenterMicro())->detail($input);


        $file = storage_path('crs_download/download/' . $data['file_name']);
        // $file = $data['file_name'];
        Log::info('文件路径:' . $file);

        if (!file_exists($file)) {
            throw new Error(3333, '该任务因为超过30天已被清理,请在系统页面重新进行导出/下载!');
        }

        $fileName = $data['name'];

        return response()->download($file, $fileName);
    }
}
