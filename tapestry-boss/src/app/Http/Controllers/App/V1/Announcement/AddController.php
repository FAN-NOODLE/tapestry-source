<?php

namespace App\Http\Controllers\App\V1\Announcement;

use App\Models\Announcement;
use Illuminate\Validation\Rule;
use App\Http\Controllers\App\V1\BaseController;

class AddController extends BaseController
{
    /**
     * 公告新增接口
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $this->validate($this->request, [
            'title'       => 'required|string|min:1',
            'type'        => [
                'required',
                'integer',
                Rule::in(Announcement::TYPE_SHOP, Announcement::TYPE_FUNCTION_ONLINE, Announcement::TYPE_NOTICE_THINGS)
            ],
            'content'     => 'required|string',
            'released_at' => 'required|integer'
        ]);
        $input = $this->request->all();
        if (!isset($input['weight'])) {
            $input['weight'] = 0;
        }

        Announcement::query()->create([
            'title'       => $input['title'],
            'type'        => $input['type'],
            'content'     => $input['content'],
            'released_at' => $input['released_at'],
            'weight'      => $input['weight']
        ]);

        return success();
    }
}
