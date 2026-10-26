<?php

namespace App\Http\Controllers\App\V1\Announcement;

use App\Models\Announcement;
use App\Services\AnnouncementService;
use Illuminate\Validation\Rule;
use App\Http\Controllers\App\V1\BaseController;

class ListController extends BaseController
{
    /**
     * 公告列表接口
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $this->validate($this->request, [
            'page'             => 'integer|min:1',
            'size'             => 'integer|min:1|max:100',
            'search'           => 'array',
            'search.title'     => 'string',
            'search.type'      => [
                'integer',
                Rule::in(Announcement::TYPE_SHOP, Announcement::TYPE_FUNCTION_ONLINE, Announcement::TYPE_NOTICE_THINGS)
            ],
            'search.time_type' => 'string|in:released_at,created_at',
            'search.start_at'  => 'integer',
            'search.end_at'    => 'integer',
        ]);
        $input = $this->request->all();
        if (!isset($input['page'])) {
            $input['page'] = 1;
        }
        if (!isset($input['size'])) {
            $input['size'] = 10;
        }

        $notPage = false;
        if (!empty($input['not_page']) && $input['not_page'] == 1) {
            $notPage = true;
        }

        $announcementService = new AnnouncementService();

        list($total, $list) = $announcementService->lists($input['page'], $input['size'], $input['search'] ?? [], ['*'], $notPage);

        return success([
            'list'  => $list,
            'total' => $total,
            'page'  => $input['page'],
            'size'  => $input['size']
        ]);
    }
}
