<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/11/16 15:08
 */

namespace App\Http\Controllers\Service\V1\PartnerMember;


use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\PartnerMember;

class DetailController  extends BaseController
{
    public function index()
    {
        $params = $this->request->all();

        // 查询条件
        $where = [];

        if (isset($params['partner_id']) && $params['partner_id']) {
            $where['id'] = $params['partner_id'];
        }
        $result = PartnerMember::query()
            ->where($where)
            ->first();
        $result = !empty($result) ? $result->toArray() : [];

        return success($result);
    }
}