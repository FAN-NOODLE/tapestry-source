<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/8/14 0014
 * Time: 17:11
 */

namespace App\Library\DataLake\Member;

use App\Library\DataLake\Base;
use App\Library\DataLake\Config;

class QueryMemberPointsDetailList extends Base
{
    public $errStr = 'query member points list error';
    public $api = '/api/interface/third/member/queryMemberPointsDetailList';

    private $queryType = "";
    private $value = "";
    private $pointTypeGroup = "";
    private $Startdate = "";
    private $Enddate = "";
    private $pageNo = "";
    private $limit = "";


    public function __construct(
         $queryType,  $value,  $pointTypeGroup,
         $Startdate,  $Enddate,  $pageNo,  $limit
    )
    {
        $this->queryType      = $queryType;
        $this->value          = $value;
        $this->pointTypeGroup = $pointTypeGroup;
        $this->Startdate      = $Startdate;
        $this->Enddate        = $Enddate;
        $this->pageNo         = $pageNo;
        $this->limit          = $limit;
        parent::__construct();
    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            "brand_code"     => $this->config->getBrandCode(),
            "program_code"   => $this->config->getProgramCode(),
            "queryType"      => $this->queryType,
            "value"          => $this->value,
            "pointTypeGroup" => $this->pointTypeGroup,
            "startDate"      => $this->Startdate,
            "endDate"        => $this->Enddate,
            "pageNo"         => $this->pageNo,
            "limit"          => $this->limit
        ];

    }

}