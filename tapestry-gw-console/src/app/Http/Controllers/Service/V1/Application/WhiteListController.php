<?php

namespace App\Http\Controllers\Service\V1\Application;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Setting;
use Illuminate\Http\Request;

class WhiteListController extends BaseController
{

    /**
     * 白名单
     *
     * @param Request $request
     *
     * @return mixed
     */
    public function index(Request $request)
    {
        $params = $request->only(['package']);

        $list = Setting::query()
            ->select('id', 'package', 'config', 'describe');

        $where = [];
        if (!empty($params['package'])) {
            $where = ['package' => $params['package']];
            $list  = $list->where($where)->first();
        } else {
            $list = $list->where($where)->get()->toArray();
        }

        return success(['data' => empty($list) ? [] : $list]);
    }

}
