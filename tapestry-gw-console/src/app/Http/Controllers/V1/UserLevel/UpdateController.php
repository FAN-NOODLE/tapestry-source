<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/26
 * Time: 10:04
 */

namespace App\Http\Controllers\V1\UserLevel;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\UserLevel;
use Illuminate\Http\Request;

class UpdateController extends BaseController
{
    /**
     * 编辑会员等级
     *
     * @param Request $request
     *
     * @return array|mixed
     */
    public function index(Request $request)
    {
        // 等级id
        $id = $request->post('id', 0);
        if (!$id) {
            error(320001);
        }

        // 等级名称
        $name = $request->post('name', '');
        if (!$name) {
            error(320002);
        }

        // 参数列表
        $params = [
            'enterprise_id' => $this->enterpriseId,
            'id'            => $id,
            'name'          => $name
        ];

        // 消费金额
        $requestParams = $request->post();
        if (isset($requestParams['amount'])) {
            $params['amount'] = $requestParams['amount'];
        }

        // 调用sdk
        $data = (new UserLevel())->update($params);

        // 写入记录日志
        $this->writeOperatorLog('编辑会员等级：'.$name);

        return success($data);
    }
}
