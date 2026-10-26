<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/23
 * Time: 17:32
 */

namespace App\Http\Controllers\V1\Commission;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\CommissionSetting;

class DetailController extends BaseController
{
    /**
     * 查看佣金设置
     *
     * @return array
     */
    public function index()
    {
        // 调用sdk
        $data = (new CommissionSetting())->detail(['enterprise_id' => $this->enterpriseId]);

        return success($data);
    }
}
