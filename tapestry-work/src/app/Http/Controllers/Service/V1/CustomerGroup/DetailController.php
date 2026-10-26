<?php


namespace App\Http\Controllers\Service\V1\CustomerGroup;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\CustomerGroupService;

class DetailController extends BaseController
{
    /**
     * 获取用户分组详情
     * @return array
     */
    public function index()
    {
        $id = $this->request->input('id', 0);

        if (empty($id)) {
            notice(810002);
        }

        $detail = (new CustomerGroupService())->customerGroupDetail($id);

        return success(['data' => $detail]);
    }

}