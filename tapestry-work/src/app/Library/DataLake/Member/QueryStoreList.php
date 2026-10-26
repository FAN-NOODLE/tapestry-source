<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/3/31 0031
 * Time: 20:59
 */

namespace App\Library\DataLake\Member;


use App\Library\DataLake\Base;

class QueryStoreList extends Base
{
    public $errStr = 'query store list error';
    public $api = '/api/interface/third/member/queryStoreList';

    private $storeNo= "";
    private $page = 1;
    private $size = 10;


    function __construct(string $storeNo="", ?int $page = 1, ?int $size = 10)
    {
        $this->storeNo = $storeNo;
        $this->page   = $page;
        $this->size   = $size;
        parent::__construct();
    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            "brand_code"   => $this->config->getBrandCode(),
            "program_code" => $this->config->getProgramCode(),
            "storeNo"      => $this->storeNo,
            "limit"        => $this->size,
            "offset"       => ($this->page - 1) * $this->size
        ];

    }
}