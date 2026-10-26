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
use Illuminate\Http\Request;

class SaveController extends BaseController
{
    /**
     * 佣金设置信息
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        // 佣金设置数据
        $commissionSetting = $request->post('commission_setting', []);

        // 参数列表
        $saveData = [
            'enterprise_id'      => $this->enterpriseId,
            'commission_setting' => $commissionSetting
        ];
        // 调用sdk
        $data = (new CommissionSetting())->save($saveData);

        // 写入记录日志
        $this->writeOperatorLog('保存佣金设置');

        return success($data);
    }
}
