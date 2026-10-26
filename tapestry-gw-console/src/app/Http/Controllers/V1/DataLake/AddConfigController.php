<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/4/1 0001
 * Time: 10:16
 */

namespace App\Http\Controllers\V1\DataLake;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\DataLake;

class AddConfigController extends BaseController
{
    private $sdk;

    public function __construct()
    {
        parent::__construct();
        $this->sdk = new DataLake();
    }

    public function index()
    {
        $params = $this->request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        return success($this->sdk->add($params));
    }
}