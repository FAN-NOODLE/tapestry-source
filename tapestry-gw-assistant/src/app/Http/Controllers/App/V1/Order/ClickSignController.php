<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/4/8 0008
 * Time: 18:08
 */

namespace App\Http\Controllers\App\V1\Order;

use App\Models\StoreMemberViewTab;
use Illuminate\Support\Facades\Auth;

class ClickSignController extends OrderBaseController
{
    /**
     * 标记推广或服务接口
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        // 验证输入参数
        $rules   = [
            'store_id' => 'required',
            'type'     => 'required|in:1,2'
        ];
        $message = [
            'store_id.required' => '店铺ID不能为空',
            'type.required'     => '类型不能为空',
            'type.in'           => '类型不符合规范'

        ];
        $this->validate($this->request, $rules, $message);

        $storeId = $this->request->post('store_id', 0);
        $type    = $this->request->post('type', 0);

        $enterpriseId = Auth::user()['enterprise_id'];
        $workerDetail = $this->storeWorkerDetail($enterpriseId, $storeId);
        // 拼接查询条件
        $where = [
            'enterprise_id'   => $enterpriseId,
            'store_id'        => $storeId,
            'store_worker_id' => $workerDetail['id']
        ];

        // 如果是推广订单查看
        if ($type == 1) {
            $data = [
                'extension_view_time' => time()
            ];
        } else {
            $data = [
                'service_view_time' => time()
            ];
        }

        // 查询是否有记录
        $total = StoreMemberViewTab::query()->where($where)->count();
        if (!$total) {
            $createData = array_merge($where, $data);
            StoreMemberViewTab::query()->create($createData);
        } else {
            StoreMemberViewTab::query()->where($where)->update($data);
        }

        return success(['data' => []]);
    }
}
