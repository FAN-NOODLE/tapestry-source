<?php

namespace App\Sdk\MaPlan\Coach\Api;

use App\Sdk\MaPlan\Coach\Request;
use GuzzleHttp\Exception\GuzzleException;

class SendFinish extends Base
{
    protected $uri = "api/coach-wecom/callback/sa";

    public function __construct($enterpriseId,$originType)
    {
        parent::__construct($enterpriseId,$originType);
        $this->uri = $this->apiUri . "/" . $this->uri;
    }

    /**
     * @throws GuzzleException
     */
    public function push($params)
    {
        $data = $params;
        $taskId = $data['events'][0]['c_task_id'];
        $crowdId = $data['events'][0]['c_crowd_id'];
        $accessToken=strtoupper(md5($taskId.'&'.$crowdId));
        $header= ["headers" => [
            "Authorization" => "Bearer {$accessToken}",
            "Content-Type" => "application/json",
            "Cache-Control" => "no-cache"
        ]];
        unset($data['enterprise_id']);
        return Request::Instance($this->enterpriseId,$this->originType)->sendPost($this->uri, $data,$header);

    }
}