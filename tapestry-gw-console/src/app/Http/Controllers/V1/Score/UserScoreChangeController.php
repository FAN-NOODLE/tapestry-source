<?php

namespace App\Http\Controllers\V1\Score;

use CrsSdk\Shop\Score;
use App\Http\Controllers\V1\BaseController;

class UserScoreChangeController extends BaseController
{
    /**
     * 管理员录入用户积分
     * @return mixed
     * @throws \Exception
     */
    public function userScoreChange()
    {
        $params = $this->request->all();

        $params['admin_id']      = $this->enterpriseAdminId;
        $params['admin_name']    = $this->enterpriseAdmin['username'];
        $params['enterprise_id'] = $this->enterpriseId;

        $this->writeOperatorLog('管理员录入用户积分');

        $result = (new Score())->userScoreChange($params);

        return success($result);
    }
}