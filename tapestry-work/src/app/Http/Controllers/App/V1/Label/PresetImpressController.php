<?php
/**
 * 获取系统预设印象标签列表
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/1/21
 * Time: 18:51
 */
namespace App\Http\Controllers\App\V1\Label;

use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\Shop\UserLabel;
use Illuminate\Http\Request;

class PresetImpressController extends BaseController
{
    /**
     * 获取系统预设印象标签列表
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        $page = $request->post('page', 1);
        $size = $request->post('size', 10);
        $search = $request->post('search', []);

        if (!isset($search['ec_userid']) || empty($search['ec_userid'])) {
            return success([
                'data' => [],
                'meta' => [
                    'total' => 0,
                    'page' => $page,
                    'size' => $size
                ]
            ]);
        }

        // 拼接企业信息
        $search['guide_id'] = $this->memUid;
        $search['enterprise_id'] = $this->enterpriseId;

        $list = (new UserLabel())->presetImpressList([
            'page' => $page,
            'size' => $size,
            'search' => $search
        ]);
        return success($list);
    }
}
