<?php

namespace App\Http\Controllers\V1\Announcement;

use CrsSdk\Boss\Announcement;
use App\Http\Controllers\V1\BaseController;

class ListController extends BaseController
{
    /**
     * 公告列表接口
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $input = $this->request->all();

        $input['search'] = [
            'time_type' => 'released_at',
            'start_at'  => 0,
            'end_at'    => time()
        ];
        $data = (new Announcement())->lists($input);

        return success($data);
    }
}
