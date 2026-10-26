<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/3/26 0026
 * Time: 16:12
 */

namespace App\Http\Controllers\Service\V1\Setting;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\DataLakeConfigService;

abstract class DataLakeConfigBase extends BaseController
{
    protected $dataLakeConfigService;
    protected $enterpriseId;
    public function __construct()
    {
        parent::__construct();
        $this->checkParam();
        $this->dataLakeConfigService = new DataLakeConfigService();
        $this->enterpriseId = $this->request->input('enterprise_id');
    }

    abstract function checkParam();
}