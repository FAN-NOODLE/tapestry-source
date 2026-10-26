<?php


namespace App\Http\Controllers\Service\V1\StaffShare;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StaffShareService;

class ListController extends BaseController
{
    public function index()
    {
        $input = $this->request->only([
            'ids',
            'user_id',
            'mem_uid',
            'worker_id',
            'worker_ids',
            'user_ids',
            'mem_uids',
            'page',
            'size',
            'is_page',
            'exclude_user_ids',
            'start_time',
            'end_time'
        ]);

        $data = (new StaffShareService())->staffList($input);

        return success($data);

    }

}