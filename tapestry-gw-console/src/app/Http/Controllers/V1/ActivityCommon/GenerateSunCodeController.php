<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/5/13 0013
 * Time: 14:35
 */

namespace App\Http\Controllers\V1\ActivityCommon;


class GenerateSunCodeController extends ActivityBaseController
{

    public function index()
    {
        $param = $this->params;
        $param['type'] = $this->typeCode; //此处类型为数字

        try {
            $this->activityApi->generateSunCode($param);
        } catch (\Exception $e) {
            $this->notice($e,$param,'generate sun code fail');
        }

        return success([]);
    }

}
