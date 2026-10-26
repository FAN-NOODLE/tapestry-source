<?php


namespace App\Http\Controllers\V1\WisdomGuide;


use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\GoMicroScrm\BaseClient;
use CrsSdk\GoMicroScrm\Guide;

class LostConfigController extends BaseController
{
    public function index()
    {
        $params = $this->request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        try {
            return (new Guide())->lostConfig($params);
        }catch (\Exception $exception){
            throw new Error($exception->getCode(), $exception->getMessage());
        }
    }

    public function info()
    {
        $params = $this->request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        try {
            return (new Guide())->lostConfigInfo($params);
        }catch (\Exception $exception){
            throw new Error($exception->getCode(), $exception->getMessage());
        }
    }
}