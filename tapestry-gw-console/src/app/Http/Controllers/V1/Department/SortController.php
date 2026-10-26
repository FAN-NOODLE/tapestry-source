<?php

namespace App\Http\Controllers\V1\Department;

use CrsSdk\Work\Department;
use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;

class SortController extends BaseController
{

    public function index(Request $request)
    {
        // 逻辑:交换两个排序值
        $input = $request->only([
            'dp_id',
            'dp_display_order',
            'last_dp_id',
            'last_dp_display_order'
        ]);

        // 调用SDK
        $params = [
            'qy_domain' => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
        ];
        $Department = new Department();

        $params['params'] = [
            'dpId' => $input['dp_id'],
            'dpDisplayorder' => $input['last_dp_display_order']
        ];
        $res1 = $Department->edit($params);
        $params['params'] = [
            'dpId' => $input['last_dp_id'],
            'dpDisplayorder' => $input['dp_display_order']
        ];
        $res2 = $Department->edit($params);

        // 保存成功
        if (isset($res1['dpId']) && isset($res2['dpId'])) {
            return success([]);
        }

        notice(320104);

    }

}