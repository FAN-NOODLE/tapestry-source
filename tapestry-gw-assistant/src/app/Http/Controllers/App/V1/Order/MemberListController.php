<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/4/8 0008
 * Time: 18:08
 */

namespace App\Http\Controllers\App\V1\Order;

use App\Models\StoreMemberViewTab;
use CrsSdk\Shop\StorePerformance;
use Illuminate\Support\Facades\Auth;

class MemberListController extends OrderBaseController
{
    /**
     * 推广人或服务人列表
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

        $search = [
            'search' => [
                'enterprise_id' => $enterpriseId,
                'store_id'      => $storeId,
                'type'          => $type
            ]
        ];

        $returnData = (new StorePerformance())->memberList($search);

        return success($returnData);
    }
}
