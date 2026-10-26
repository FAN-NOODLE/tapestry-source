<?php
/**
 * Created by PhpStorm.
 * Function: 企业员工业绩排行
 * User: 龚国新
 * Date: 2020/2/14
 * Time: 11:43
 */

namespace App\Http\Controllers\App\V1\StaffOrder;


use App\Http\Controllers\App\V1\BaseController;
use App\Models\StaffOrderRelation;
use App\Services\StaffManagerService;
use App\Services\StaffOrderService;
use CrsSdk\Shop\Goods;
use CrsSdk\Shop\StaffOrder;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;

class OrderListController extends BaseController
{
  //  protected $requireLogin = false;
    /**
     * 企业员工对应订单列表
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
   //     $_SERVER['enterprise_hash']='054edce1d4cec91c76001d968841c08f';
    //      $this->enterpriseId=2;
    //      $this->memUid='561CF4F37F00000109CD40645574BF03';
        $input = $request->only(['page', 'size','mem_uid']);
        $page  = $input['page'] ?? 1;
        $size  = $input['size'] ?? 10;

        if (null == $this->memUid) {
            $this->memUid = $input['mem_uid']??'';
        }

        $obj = StaffOrderRelation::query()
            ->where('enterprise_id', $this->enterpriseId)
            ->where('mem_uid', $this->memUid)
            ->where('order_no', '<>', 1)
            ->orderByDesc('created_at')
            ->paginate($size, ['*'], '', $page);

        $list = $this->format($obj->items());

        $data = [
            'meta' => [
                'total' => $obj->total(),
                'page'  => (int)$page,
                'size'  => (int)$size
            ],
            'data' => $list
        ];

        return success($data);
    }

    /**
     * 格式化代码
     * @param array $data
     *
     * @return array
     * @throws \Exception
     */
    public function format($data = [])
    {
        if (empty($data)) {
            return $data;
        }

        $orderNos  = array_column($data, 'order_no');
        $orderData = (new StaffOrder())->orderInfoList(['order_nos' => $orderNos]);

        $orderList = array_column($orderData['data'], null, 'order_no');
        $list      = [];
        foreach ($data as $value) {
            $orderInfo = $orderList[$value['order_no']] ?? [];
            $list[]    = [
                'id'           => $orderInfo['id'] ?? 0,
                'order_no'     => $value['order_no'],
                'status'       => $orderInfo['status'] ?? 0,
                'created_at'   => $orderInfo['created_at'] ?? 0,
                'pay_money'    => $value['pay_money'] ?? 0,
                'mem_username' => $orderInfo['nick_name'] ?? 0,
                'avatar'       => $orderInfo['avatar'] ?? '',
                'good_num'     => $orderInfo['goods_num'] ?? 0,
            ];
        }

        return $list;
    }
}