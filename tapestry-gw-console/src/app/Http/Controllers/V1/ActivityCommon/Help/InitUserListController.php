<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/5/19 0019
 * Time: 20:28
 */

namespace App\Http\Controllers\V1\ActivityCommon\Help;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroActivity\Help;
use CrsSdk\MicroGoodsV1\Activity;

class InitUserListController extends BaseController
{
    /**
     * 发起助力用户列表
     *
     * @return \Illuminate\Http\JsonResponse
     * @author  baizhongyang <baizhongyang@vchangyi.com>
     * @date    2020/6/15 11:20
     */
    public function index()
    {
        $params = $this->request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $result = (new Help())->initUserList($params);
        return success($result);
    }

}
