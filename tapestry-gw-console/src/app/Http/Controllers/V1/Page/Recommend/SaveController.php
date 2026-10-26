<?php

namespace App\Http\Controllers\V1\Page\Recommend;

use CrsSdk\Shop\PageHot;
use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;
use CrsSdk\Shop\PageActivityBanner;
use App\Http\Controllers\V1\BaseController;

class SaveController extends BaseController
{
    /**
     * 爆品推荐、店长推荐、猜你喜欢、主推商品 - 修改保存
     * @param Request $request
     * @return JsonResponse|void
     */
    public function index(Request $request)
    {
        $params = $request->only(['data', 'banner', 'type']);

        // 0爆品推荐  1店长推荐  2猜你喜欢  3主推商品
        if (!in_array($params['type'] ?? false, [0, 1, 2, 3])) {
            notice(700207);
        }

        // 添加修改商品推荐
        $this->addGoods($params['data'], $params['type']);

        // 写入日志
        $this->writeOperatorLog('编辑轮播图');

        return success();
    }

    /**
     * 添加保存商品关联
     * @param $data
     * @param $type
     * @return array
     */
    public function addGoods($data, $type)
    {
        $params = array_merge(['data' => $data], [
            'type'          => $type,
            'enterprise_id' => $this->enterpriseId,
        ]);

        (new PageHot())->add($params);

        return [];
    }
}
