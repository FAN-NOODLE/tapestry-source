<?php
/**
 * 内购业绩订单列表
 * User: liurui
 * Date: 2018/11/26
 * Time: 14:10
 */

namespace App\Http\Controllers\V1\StaffOrder;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\StaffOrder;
use Illuminate\Http\Request;

class AchievementListController extends BaseController
{

    /**
     * @param Request $request
     * @return mixed
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $input = $request->only(['search', 'page', 'size']);
        $input['search']['enterprise_id'] = $this->enterpriseId;
        $input['search']['qy_domain'] = $this->qyDomain;
        $input['search']['app_identifier'] = $this->appIdentifier;

        $data= (new StaffOrder())->achievementList($input);

        return  $data;
    }
}
