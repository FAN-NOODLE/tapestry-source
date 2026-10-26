<?php
/**
 * Created by PhpStorm.
 * User: guokai
 * Date: 2019/04/18
 * Time: 18:05
 */

namespace App\Http\Controllers\V1\Page\ReduceMoney;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\PageReduce;
use Illuminate\Http\Request;

class SaveController extends BaseController
{
    /**
     * 修改立减金
     *
     * @param Request $request
     * @return array|mixed
     */
    public function index(Request $request)
    {
        // 接受查询参数
        $data = $request->input();

        // 必填字段验证
        if (!isset($data['reduce_id']) && empty($data['reduce_id'])) {
            error(760001);
        }

        $result = (new PageReduce())->save($data);

        try{
            $this->writeOperatorLog('修改立减金');
        }catch (\Exception $e){
        }

        return success($result);
    }
}
