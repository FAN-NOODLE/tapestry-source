<?php


namespace App\Http\Controllers\Service\V1\CustomerGroup;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\CustomerGroupService;
use Illuminate\Support\Facades\Log;

class DeleteController extends BaseController
{
    /**
     * 删除用户分组
     * @return array
     */
    public function index(){
        $id = $this->request->input('id', 0);

        if (empty($id)) {
            notice(810002);
        }

        // 删除用户组
        try {

            (new CustomerGroupService())->customerGroupDel($id);

        } catch (\Exception $exception) {
            Log::info('用户分组删除失败:'.$exception->getMessage(), $this->request->input());
            notice(810003);
        }

        return success(['data' => []]);
    }
}