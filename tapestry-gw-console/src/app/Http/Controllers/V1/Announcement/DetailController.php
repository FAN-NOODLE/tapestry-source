<?php

namespace App\Http\Controllers\V1\Announcement;

use CrsSdk\Boss\Announcement;
use App\Http\Controllers\V1\BaseController;

class DetailController extends BaseController
{
    /**
     * 公告详情接口
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $input = $this->request->all();

        $data = (new Announcement())->detail($input);

        return success($data);
    }
}
