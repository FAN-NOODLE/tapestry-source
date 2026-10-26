<?php
/**
 * 提现 - 银行卡列表
 * Created by PhpStorm.
 * User: Yang Hui
 * Date: 2020年10月16日16:56:4
 */

namespace App\Http\Controllers\Service\V1\Withdrawal;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\ShoppingGuideCashWithdrawal;
use App\Models\StoreWorkerCertification;
use Illuminate\Http\Request;
use App\Services\WithdrawalServices;

class WithdrawalBankListController extends BaseController
{
    /**
     * 提现 - 银行卡列表
     * @param Request $request
     * @return mixed
     * @throws \Exception
     * @author   Yang Hui
     * @date     2020年10月16日17:1:28
     */
    public function index(Request $request)
    {
        $param = $request->only([
            'page',
            'size',
            'bank_name'
        ]);
        $param['page'] = empty($param['page']) ? 1 : $param['page'];
        $param['size'] = empty($param['size']) ? 1000 : $param['size'];

        $data = [
            ['code' => 498, 'title' => "中国银行总行", 'order_num' => 0, 'status' => 1],
            ['code' => 500, 'title' => "中国工商银行总行营业部", 'order_num' => 0, 'status' => 1],
            ['code' => 450, 'title' => "中国农业银行股份有限公司", 'order_num' => 0, 'status' => 1],
            ['code' => 499, 'title' => "中国建设银行股份有限公司总行", 'order_num' => 0, 'status' => 1],
            ['code' => 64, 'title' => "中国邮政储蓄银行有限责任公司", 'order_num' => 0, 'status' => 1],
            ['code' => 289, 'title' => "招商银行股份有限公司", 'order_num' => 0, 'status' => 1],
            ['code' => 473, 'title' => "交通银行", 'order_num' => 0, 'status' => 1],
            ['code' => 502, 'title' => "中信银行股份有限公司", 'order_num' => 0, 'status' => 1],
            ['code' => 454, 'title' => "兴业银行总行", 'order_num' => 0, 'status' => 1],
            ['code' => 459, 'title' => "中国光大银行", 'order_num' => 0, 'status' => 1],
            ['code' => 336, 'title' => "平安银行股份有限公司", 'order_num' => 0, 'status' => 1],
            ['code' => 163, 'title' => "上海浦东发展银行", 'order_num' => 0, 'status' => 1],
            ['code' => 325, 'title' => "广发银行股份有限公司", 'order_num' => 0, 'status' => 1],
            ['code' => 434, 'title' => "华夏银行股份有限公司总行", 'order_num' => 0, 'status' => 1]
        ];
        $result = [];
        if (!empty($param['bank_name'])) {
            foreach ($data as $k => $v) {
                if (strstr($v['title'], $param['bank_name']) !== false) {
                    $result[] = $v;
                }
            }
        } else {
            $result = $data;
        }
        $result = array_slice($result, intval(($param['page']-1)*$param['size']), $param['size']);

        $returnData = [
            'data' => $result
        ];
        return success($returnData);
    }
}
