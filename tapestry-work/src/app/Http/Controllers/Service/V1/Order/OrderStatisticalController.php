<?php
/**
 * Created by PhpStorm.
 * User: baizhongyang
 * Time: 13:14
 */

namespace App\Http\Controllers\Service\V1\Order;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\CrsOrder;
use App\Services\GuideManagerService;
use App\Services\PartnerMemberService;
use CrsSdk\Shop\Order;
use Illuminate\Support\Facades\DB;

class OrderStatisticalController extends BaseController
{
    /**
     * 订单统计
     *
     * @return array
     * @author  baizhongyang <baizhongyang@vchangyi.com>
     * @date    2020/8/19 14:13
     */
    public function index()
    {
        $param        = $this->getParam();
        $qyDomain     = $param['qy_domain'];
        $adminId      = $param['admin_id'];
        $mobile       = $param['mobile'];
        $enterpriseId = $param['enterprise_id'];
        $storeId      = $param['store_id'];

        switchAppHost(['enterprise_id' => $enterpriseId]);

        // 获取导购列表
        $guideList           = GuideManagerService::getGuideListForStatistic($param, $mobile, $qyDomain, $enterpriseId, $adminId, $storeId);
        $param['worker_ids'] = array_column($guideList['guide_list'], 'id');
        //获取合伙人ID
        $param['partner_ids'] = (new PartnerMemberService())->getPartnerIdsByWorkerIds(
            $enterpriseId,
            $param['worker_ids']
        );

        // 获取订单统计数据
        $data = (new Order())->wisdomOrderStatistical($param);

        return success($data);
    }

    /**
     * 验证参数
     *
     * @return array|null|string
     * @author  baizhongyang <baizhongyang@vchangyi.com>
     * @date    2020/8/18 15:22
     */
    private function getParam()
    {
        $input    = $this->request->input();
        $rules    = [
            'enterprise_id' => 'integer|required',
            'store_id'      => 'integer',
            'mobile'        => 'required',
            'qy_domain'     => 'required',
            'admin_id'      => 'required',
            'start_time'    => 'required',
            'end_time'      => 'required',
        ];
        $messages = [
            'enterprise_id.*'     => '企业id参数不合法',
            'store_id.*'          => '门店id参数不合法',
            'mobile.*'            => '用户手机号参数不合法',
            'qy_domain.required'  => 'qy_domain参数缺失',
            'admin_id.required'   => 'admin_id参数缺失',
            'start_time.required' => 'start_time参数缺失',
            'end_time.required'   => 'end_time参数缺失',
        ];
        $this->validate($this->request, $rules, $messages);

        return $input;
    }
}