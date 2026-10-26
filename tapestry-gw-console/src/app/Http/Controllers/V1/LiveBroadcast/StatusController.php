<?php
/**
 * Created by PhpStorm.
 * User: admin
 * Date: 2019/8/16
 * Time: 16:17
 */

namespace App\Http\Controllers\V1\LiveBroadcast;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroActivity\Live;

class StatusController extends BaseController
{
    public function index()
    {
        $input = $this->request->input();

        $input['enterprise_id'] = $this->enterpriseId;
        try {
            $this->writeOperatorLog('更新直播状态');
        } catch (\Exception $e) {
        }
        $return = (new Live())->status($input);
        return success($return);
    }
}
