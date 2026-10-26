<?php


namespace App\Http\Controllers\Service\V1\CustomerGroup;

use App\Http\Controllers\Service\V1\BaseController;
use App\Http\Requests\CustomerGroupRequest;
use App\Services\CustomerGroupService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class AddController extends BaseController
{
    /**
     * 添加用户分组
     * @param Request $request
     * @param CustomerGroupService $customerGroupService
     * @return array
     */
    public function index(
        Request $request,
        CustomerGroupService $customerGroupService
    )
    {
        // 参数
        $input = $request->only([
            'enterprise_id',
            'name',
            'user_visit',
            'coupon',
            'order',

        ]);

        CustomerGroupRequest::add($input);

        $input['user_visit'] = json_encode($input['user_visit']);
        $input['coupon'] = json_encode($input['coupon']);
        $input['order'] = json_encode($input['order']);
        $input['name'] = trim($input['name']);

        $deTailByName = $customerGroupService->customerGroupDetailByName($input['name']);
        if ($deTailByName) {
            notice(810009);
        }

//        $detailByCondition = $customerGroupService->customerGroupDetailByCondition($input);
//        if ($detailByCondition) {
//            notice(810010);
//        }

        // 插入
        try {

            $data = $customerGroupService->insertCustomerGroup($input);

        } catch (\Exception $exception) {

            Log::info('用户分组数据插入错误'.$exception->getMessage(), $input);
            notice(810001);
        }

        return success([
            'data' => [
                'id' => $data['id'] ?? 0
            ]
        ]);
    }

}