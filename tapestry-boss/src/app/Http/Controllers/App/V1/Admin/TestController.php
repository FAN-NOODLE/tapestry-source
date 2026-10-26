<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/27
 * Time: 2:32 PM
 */

namespace App\Http\Controllers\App\V1\Admin;

use App\Http\Controllers\App\V1\BaseController;

class TestController extends BaseController
{
    public function index()
    {
        return success(['data' => $this->adminId]);
    }
}
