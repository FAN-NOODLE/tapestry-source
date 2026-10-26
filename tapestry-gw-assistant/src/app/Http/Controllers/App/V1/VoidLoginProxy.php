<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/10/19 0019
 * Time: 16:19
 */

namespace App\Http\Controllers\App\V1;


class VoidLoginProxy extends BaseController
{
    private $mock = true;

    public function __construct()
    {
        if (config('app.env') == 'dev' && $this->mock) {
            $this->isRequireLogin = false;
            $this->enterpriseId = 2;
            $this->mobile=17391604660;
        }
        parent::__construct();

    }
}