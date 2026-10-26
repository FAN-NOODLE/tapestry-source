<?php

namespace App\Http\Controllers\V1\Industry;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Boss\Industry;

class ListController extends BaseController
{
    protected $requireLogin = false;

    /**
     * 行业列表
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $input = $this->request->all();

        $data = (new Industry())->lists($input);

        return success($data);
    }
}
