<?php


namespace App\Http\Controllers\Service\V1\CustomerGroup;


use App\Http\Controllers\Service\V1\BaseController;
use App\Http\Requests\CustomerGroupRequest;
use App\Services\CustomerGroupService;
use Illuminate\Support\Facades\Log;

class EditController extends BaseController
{
    public function index()
    {
        // 参数
        $input = $this->request->only([
            'id',
            'enterprise_id',
            'name',
            'user_visit',
            'coupon',
            'order',

        ]);

        if (!isset($input['id']) || empty($input['id'])) {
            notice(810002);
        }

        CustomerGroupRequest::add($input);

        $input['user_visit'] = json_encode($input['user_visit']);
        $input['coupon'] = json_encode($input['coupon']);
        $input['order'] = json_encode($input['order']);
        $input['name'] = trim($input['name']);

        // 重复的名称
        $deTailByName = (new CustomerGroupService())->customerGroupDetailByName($input['name'], $input['id']);
        if ($deTailByName) {
            notice(810009);
        }

          // 重复的条件
//        $detailByCondition = $customerGroupService->customerGroupDetailByCondition($input);
//        if ($detailByCondition) {
//            notice(810010);
//        }

        // 更新用户分组信息
        try {

            (new CustomerGroupService())->customerGroupEdit($input);

        } catch (\Exception $exception) {

            Log::info('更新用户分组失败'.$exception->getMessage(), $input);

        }

        return success(['data' => []]);
    }

}