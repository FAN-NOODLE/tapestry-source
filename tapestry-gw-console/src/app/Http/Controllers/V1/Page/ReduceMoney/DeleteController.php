<?php
/**
 * Created by PhpStorm.
 * User: guokai
 * Date: 2019/04/18
 * Time: 18:25
 */

namespace App\Http\Controllers\V1\Page\ReduceMoney;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\PageReduce;
use Illuminate\Http\Request;

class DeleteController extends BaseController
{
    /**
     * 删除商品
     *
     * @param Request $request
     *
     * @return array|\Illuminate\Http\JsonResponse
     */
    public function index(Request $request)
    {
        // 接受查询参数
        $data = $request->input();

        // 必填字段验证
        if (!isset($data['reduce_id']) && empty($data['reduce_id'])) {
            error(720001);
        }

        try{
            $this->writeOperatorLog('删除立减金');
        }catch (\Exception $e){
        }

        $result = (new PageReduce())->delete(['reduce_id' => $data['reduce_id']]);

        if (empty($result['data']['error_ids'])) {
            return success(['data' => []]);
        }

        // 微信立减金被占位符占用，无法删除
        if (!empty($result['data']['placeholder_error_ids'])) {
            error(760002, '占位符');
        }
    }
}
