<?php

namespace App\Http\Controllers\V1\Page\Recommend;

use CrsSdk\Shop\PageHot;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class ListController extends BaseController
{
    /**
     * 推荐列表
     * @return JsonResponse|void
     */
    public function index()
    {
        $type = $this->request->input('type');

        if (empty($type)) {
            notice(310206);
        }

        $list = (new PageHot())->lists([
            'type'          => $type,
            'enterprise_id' => $this->enterpriseId,
        ]);

        return success([
            'list' => $list,
        ]);
    }
}