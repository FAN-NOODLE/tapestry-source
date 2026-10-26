<?php
/**
 * 下载上传物流信息失败的订单
 * Created by PhpStorm.
 * User: cjp
 * Date: 2019/6/5
 * Time: 15:48
 */
namespace App\Http\Controllers\V1\Order;

use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;
use App\Services\OrderService;
use Maatwebsite\Excel\Facades\Excel;

class FailOrderDownloadController extends BaseController
{
    public function index(OrderService $orderService)
    {
        //初始化
        $fileInfo = [];
        if ($this->request->filled('fail_list')) {
            $fileInfo = json_decode($this->request->input('fail_list'), true);
            if (!is_array($fileInfo)) $fileInfo = [];
        }

        // 把order_no转成string
        foreach ($fileInfo as $key => $value) {
            $value[$key]['order_no'] = "\t" . (string)$value['order_no'] . "\t";
        }
        array_unshift($fileInfo, $orderService->valueHand);

        $name = "上传失败列表" . date('YmdHi', time());

        return Excel::create($name, function ($excel) use ($fileInfo) {
            $excel->sheet('score', function ($sheet) use ($fileInfo) {
                $sheet->rows($fileInfo);
            });
        })->export('xlsx');
    }
}