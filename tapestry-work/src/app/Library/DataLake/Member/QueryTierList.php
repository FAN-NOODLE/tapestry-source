<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/3/26 0026
 * Time: 14:40
 */

namespace App\Library\DataLake\Member;


use App\Library\DataLake\Base;

class QueryTierList extends Base
{
    public $errStr = 'query tier list error';

    public $api = '/api/interface/third/system/queryTierList';


    protected function setQueryParam()
    {
        $this->queryParam = [
            "brand_code" => $this->config->getBrandCode(),
            "program_code"=>$this->config->getProgramCode()
        ];
    }
}