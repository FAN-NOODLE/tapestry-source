<?php

namespace App\Http\Controllers\App\V1\Announcement;

use App\Models\Announcement;
use Illuminate\Validation\Rule;
use App\Http\Controllers\App\V1\BaseController;

class DeleteController extends BaseController
{
    /**
     * 公告删除接口
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $this->validate($this->request, [
            'announcement_id' => 'required|integer'
        ]);
        $input = $this->request->all();

        $announcement = Announcement::query()->where('id', '=', $input['announcement_id'])->first();
        if ($announcement === null) {
            return notice(100000, ['公告数据不存在']);
        }
        $announcement->delete();

        return success();
    }
}
