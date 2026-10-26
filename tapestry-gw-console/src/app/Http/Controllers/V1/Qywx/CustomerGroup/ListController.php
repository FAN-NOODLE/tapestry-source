<?php

namespace App\Http\Controllers\V1\Qywx\CustomerGroup;

use CrsSdk\Work\CustomerGroup;
use App\Http\Controllers\V1\BaseController;

class ListController extends BaseController
{
    /**
     * 获取用户分组列表
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        $input                  = $this->request->input();
        $input['enterprise_id'] = $this->enterpriseId;

        $list = (new CustomerGroup())->customerGroupList($input);

        if (isset($input['page']) || isset($input['size'])) {
            return success([
                'list'  => $list['data'],
                'page'  => $list['meta']['page'],
                'size'  => $list['meta']['size'],
                'total' => $list['meta']['total'],
            ]);
        }

        return success([
            'list' => $list['data']
        ]);
    }
}