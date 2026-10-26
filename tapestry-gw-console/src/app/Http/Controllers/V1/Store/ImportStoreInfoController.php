<?php

namespace App\Http\Controllers\V1\Store;

use App\Http\Controllers\V1\BaseController;
use App\Traits\ExcelTrait;
use CrsSdk\Sdk;
use Illuminate\Http\Request;

class ImportStoreInfoController extends BaseController
{
    use ExcelTrait;
    /**
     * @param Request $request
     *
     * @return mixed|\Psr\Http\Message\StreamInterface
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request)
    {
        // 参数验证
        $fileUrl = $this->request->post('file_url');
        $savePath = storage_path('/' ) . ltrim(strrchr($fileUrl,'/'),'/');
        file_put_contents($savePath, file_get_contents($fileUrl));
        try{
            $data = $this->excelRead($savePath);
        }catch (\Exception $exception){
            error(1002,[1002=>"文件读取失败"],[],200);
        }
        $params = [
            'enterprise_id'=>$this->enterpriseId,
            'data'  =>$data
        ];
        unlink($savePath);
        $result = Sdk::work()
            ->setApi('service/v1/store/batch-update-store-info')
            ->post($params);

        return success($result);
    }
}
