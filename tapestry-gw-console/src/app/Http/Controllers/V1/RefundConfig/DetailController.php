<?php
/**
 * Created by PhpStorm.
 * User: admin
 * Date: 2019/6/24
 * Time: 17:36
 */

namespace App\Http\Controllers\V1\RefundConfig;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\RefundConfig;

class DetailController extends BaseController
{
    public function index()
    {
        $input = $this->request->input();
        $input['enterprise_id'] = $this->enterpriseId;

        $response = (new RefundConfig())->detail($input);

        return success($response);
    }
}