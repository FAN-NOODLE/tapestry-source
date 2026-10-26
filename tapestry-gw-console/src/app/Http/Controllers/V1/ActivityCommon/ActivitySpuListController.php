<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/5/13 0013
 * Time: 14:35
 */

namespace App\Http\Controllers\V1\ActivityCommon;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroActivity\Activity;

class ActivitySpuListController extends BaseController
{

    public function index()
    {
        $param = $this->request->input();
        $param['enterprise_id'] = $this->enterpriseId;
        try {
            $list = (new Activity)->activitySpuList($param);
        } catch (\Exception $e) {
            throw new Error($e->getCode(), $e->getMessage());
        }

        return success($list ?? []);
    }

}
