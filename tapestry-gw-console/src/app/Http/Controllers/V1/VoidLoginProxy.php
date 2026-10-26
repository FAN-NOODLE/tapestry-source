<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/10/19 0019
 * Time: 16:43
 */

namespace App\Http\Controllers\V1;


class VoidLoginProxy extends BaseController
{
    private $mock = true;

    public function __construct()
    {
        if (config('app.env') == 'dev' && $this->mock) {
            $this->requireLogin = false;
            $this->enterpriseId = 2;
        }
        parent::__construct();

    }

}