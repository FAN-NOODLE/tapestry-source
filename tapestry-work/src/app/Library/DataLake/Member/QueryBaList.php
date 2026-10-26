<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/3/31 0031
 * Time: 20:59
 */

namespace App\Library\DataLake\Member;


use App\Library\DataLake\Base;

class QueryBaList extends Base
{
    public $errStr = 'query ba list error';
    public $api = '/api/interface/third/member/queryBAList';

    private $page = 1;
    private $size = 10;
    private $baCode;

    function __construct(string $baCode, ?int $page = 1, ?int $size = 10)
    {
        $this->baCode = $baCode;
        $this->page   = $page;
        $this->size   = $size;
        parent::__construct();
    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            "brand_code"   => $this->config->getBrandCode(),
            "program_code" => $this->config->getProgramCode(),
            "BA_code"      => $this->baCode,
            "limit"        => $this->size,
            "offset"       => ($this->page - 1) * $this->size
        ];

    }
}