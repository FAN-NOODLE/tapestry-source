<?php

namespace App\Sdk\MaPlan\Coach\Api;

use App\Sdk\MaPlan\Coach\Request;
use GuzzleHttp\Exception\GuzzleException;

class MomentSendResult extends Base
{
    protected $uri = "api/coach-wecom/callback/moment_send_result";

    public function __construct($enterpriseId,$originType)
    {
        parent::__construct($enterpriseId,$originType);
        $this->uri = $this->apiUri . "/" . $this->uri;
    }

    /**
     * @param $data
     * @return array|mixed
     * @throws GuzzleException
     */
    public function push($data)
    {
        unset($data['enterprise_id']);
        return Request::Instance($this->enterpriseId,$this->originType)->sendPost($this->uri, $data,$this->getHeader($data));

    }
}