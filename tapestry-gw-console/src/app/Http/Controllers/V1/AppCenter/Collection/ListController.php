<?php
/**
 * Created by PhpStorm.
 * User: zhuhongda
 * Date: 2019/5/14
 * Time: 10:04
 */

namespace App\Http\Controllers\V1\AppCenter\Collection;

use CrsSdk\Shop\Collection;
use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;

class ListController extends BaseController
{
    /**
     * 用户提交收藏有礼列表
     * @author 朱虹达 zhuhongda@vchangyi.com
     * @date 2019/05/14 11:16
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        // 接受参数
        $input = $request->only([
            'state',
            'search',
            'page',
            'size'
        ]);

        $input['size']          = $input['size'] ?? 10;
        $input['enterprise_id'] = $this->enterpriseId;

        // 查询列表
        return (new Collection())->lists($input);
    }

}
