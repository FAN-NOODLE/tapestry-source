<?php
/**
 * Created By PhpStorm
 * @Date: 2021/6/24
 * @Time: 12:01
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: StatisticController.php
 */

namespace App\Http\Controllers\App\Brand\CommunityStatistics;


use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\Work\DataCompass;

class StatisticController extends BaseController
{
    /**
     * 社群数据统计
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     * @Date: 2021/6/24
     * @Time: 14:54
     */
    public function index()
    {
        $params                  = $this->request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->loginMobile();
        $params['qy_domain']     = $this->loginAdminEnterprise['domain'];
        $params['admin_id']      = $this->adminId;

        $statistics = (new DataCompass())->countForApp($params);

        return success($statistics);
    }

    /**
     * 未参加群活动人员列表
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     * @Date: 2021/6/24
     * @Time: 14:55

     */
    public function getNotJoinActivityMemberList()
    {
        $params                  = $this->request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->loginMobile();
        $params['qy_domain']     = $this->loginAdminEnterprise['domain'];
        $params['admin_id']      = $this->adminId;

        $statistics = (new DataCompass())->getNotJoinActivityMemberList($params);

        return success($statistics);
    }
}