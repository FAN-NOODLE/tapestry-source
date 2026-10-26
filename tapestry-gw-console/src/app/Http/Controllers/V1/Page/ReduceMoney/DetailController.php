<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2019/4/19
 * Time: 22:04
 */

namespace App\Http\Controllers\V1\Page\ReduceMoney;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\PageReduce;
use Illuminate\Http\Request;

class DetailController extends BaseController
{
    /**
     * 详情商品
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

        $result = (new PageReduce())->detail(['reduce_id' => $data['reduce_id']]);

        return success($result);
    }
}