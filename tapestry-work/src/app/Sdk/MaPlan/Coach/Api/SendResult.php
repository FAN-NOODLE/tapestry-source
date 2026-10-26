<?php

namespace App\Sdk\MaPlan\Coach\Api;

use App\Sdk\MaPlan\Coach\Request;
use GuzzleHttp\Exception\GuzzleException;

class SendResult extends Base
{

    protected $uri = "api/coach-wecom/callback/task";

    public function __construct($enterpriseId,$originType)
    {
        parent::__construct($enterpriseId,$originType);
        $this->uri = $this->apiUri . "/" . $this->uri;
    }

    /**
     * @param $params
     * @return array|mixed
     * @throws GuzzleException
     */
    public function push($params)
    {
        unset($params['enterprise_id']);
        return Request::Instance($this->enterpriseId,$this->originType)->sendPost($this->uri, $params,$this->getHeader($params));
    }
}