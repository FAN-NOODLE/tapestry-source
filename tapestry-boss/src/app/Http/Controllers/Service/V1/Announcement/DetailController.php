<?php

namespace App\Http\Controllers\Service\V1\Announcement;

use App\Services\AnnouncementService;
use App\Http\Controllers\Service\V1\BaseController;

class DetailController extends BaseController
{
    /**
     * 公告详情接口
     *
     * @return \Illuminate\Http\JsonResponse|void
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $this->validate($this->request, [
            'announcement_id' => 'required|integer'
        ]);
        $input = $this->request->all();

        $announcementService = new AnnouncementService();
        $announcement = $announcementService->detail($input['announcement_id']);
        if ($announcement === null) {
            return notice(100000, ['公告数据不存在']);
        }

        return success($announcement->toArray());
    }
}
