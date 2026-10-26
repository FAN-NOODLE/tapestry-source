<?php

namespace App\Http\Controllers\V1\Order;

use App\Http\Controllers\V1\BaseController;
use App\Services\OrderService;
use CrsSdk\Shop\Order;
use Exception;

class SendGoodsBatchController extends BaseController
{
    /**
     * 上传文件批量发货接口
     *
     * @return \Illuminate\Http\JsonResponse|void
     * @throws \Illuminate\Validation\ValidationException
     * @throws \PHPExcel_Exception
     * @throws \PHPExcel_Reader_Exception
     * @throws Exception
     */
    public function index()
    {
        $input = $this->request->only(['order_file']);

        $this->validate($this->request, [
            'order_file' => 'required|file',
        ], [
            'order_file.*' => '请上传文件',
        ]);

        $orderService = new OrderService();
        // 获取文件路径
        $fileName = $orderService->uploadFile($input['order_file']);

        // 获取文件内容
        $fileInfo = $orderService->excelImport($fileName);

        // 删除表头
        unset($fileInfo[0]);
        $orderList = [];
        foreach ($fileInfo as $info) {
            $orderList[] = array_combine($orderService->keyHand, $info);
        }

        if (empty($orderList)) {
            error(320263);
        }

        $data = (new Order())->sendGoodsBatch([
            'enterprise_id'   => $this->enterprise['id'],
            'enterprise_hash' => $this->enterprise['hash'],
            'data'            => $orderList
        ]);

        return success($data);
    }

    /**
     * 下载手工发货模板
     *
     * @return \Symfony\Component\HttpFoundation\BinaryFileResponse
     */
    public function downloadTemplate()
    {
        $path = resource_path('static/手工发货模板.xlsx');

        return response()->download($path);
    }
}
