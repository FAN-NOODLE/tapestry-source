<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/26
 * Time: 10:05
 */

namespace App\Http\Controllers\V1\UserLevel;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\UserLevel;
use Illuminate\Http\Request;

class SetStatusController extends BaseController
{
    /**
     * 变更会员等级启用状态
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

        // 状态
        $status = $request->post('status', 0);
        if (!$status) {
            error(320003);
        }

        // 参数列表
        $params = [
            'enterprise_id' => $this->enterpriseId,
            'id'            => $id,
            'status'        => $status
        ];
        // 调用sdk
        $data = (new UserLevel())->setStatus($params);

        // 写入记录日志
        $name = $status==1 ? '启用' : '禁用';
        $this->writeOperatorLog('变更会员等级状态：'.$name);

        return success($data);
    }
}
