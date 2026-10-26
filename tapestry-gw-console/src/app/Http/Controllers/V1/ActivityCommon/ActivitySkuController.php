<?php
/**
 * 弃用
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/5/13 0013
 * Time: 14:35
 */

namespace App\Http\Controllers\V1\ActivityCommon;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroActivity\Activity;

class ActivitySkuController extends BaseController
{

    /**
     * 弃用改方法，整合到了sku列表接口中
     */
    public function index()
    {
        $param = $this->request->input();
        $param['enterprise_id'] = $this->enterpriseId;

        try {
            $list = (new Activity)->activitySku($param);
        } catch (\Exception $e) {
            throw new Error($e->getCode(), $e->getMessage());
        }

        return success($list ?? []);
    }

}
