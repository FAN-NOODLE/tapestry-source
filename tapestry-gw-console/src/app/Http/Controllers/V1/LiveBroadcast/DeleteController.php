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

class DeleteController extends BaseController
{
    public function index()
    {
        $input = $this->request->input();

        $input['enterprise_id'] = $this->enterpriseId;
        try {
            $this->writeOperatorLog('删除直播');
        } catch (\Exception $e) {
        }
        $return = (new Live())->del($input);
        return success($return);
    }
}
