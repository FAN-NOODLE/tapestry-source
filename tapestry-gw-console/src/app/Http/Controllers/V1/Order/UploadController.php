<?php
/**
 * 上传文件同步订单的发货物流信息
 * Created by PhpStorm.
 * User: cjp
 * Date: 2019/6/4
 * Time: 11:03
 */
namespace App\Http\Controllers\V1\Order;

use App\Http\Controllers\V1\BaseController;
use App\Services\OrderService;
use CrsSdk\Shop\Order;
use Illuminate\Http\Request;

class UploadController extends BaseController
{

    /**
     * @param Request      $request
     * @param OrderService $orderService
     *
     * @return mixed|\Psr\Http\Message\StreamInterface
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request, OrderService $orderService)
    {
        $input = $request->only(['order_file']);

        $rules = [
            'order_file' => 'required|file',
        ];

        $message = [
            'order_file.*' => '请上传文件',
        ];

        $this->validate($request, $rules, $message);

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

        return (new Order())->importOrder([
            'enterprise_hash' => $this->enterprise['hash'],
            'data'            => $orderList
        ]);
    }
}