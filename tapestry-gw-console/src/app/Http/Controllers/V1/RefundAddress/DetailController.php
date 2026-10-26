<?php
/**
 * Created by PhpStorm.
 * User: chengyujia
 * Date: 2019/07/12
 * Time: 10:04
 */


namespace App\Http\Controllers\V1\RefundAddress;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\RefundAddress;

class DetailController extends BaseController
{
    /**
     * @return array
     */
    public function index()
    {
        $params = request()->input();
        $params['enterprise_id'] = $this->enterpriseId;

        return success((new RefundAddress())->detail($params));
    }
}
