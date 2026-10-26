<?php

namespace App\Sdk\MaPlan\Coach\Api;

use App\Exceptions\Error;
use App\Sdk\MaPlan\Coach\BaseCore;

class Base extends BaseCore
{

    public $apiUri;


    public function __construct($enterpriseId,$originType)
    {
        parent::__construct($enterpriseId,$originType);
        if (config("app.env") == "prod") {
            $this->apiUri = "https://ubaweb-newcdp.coach.com.cn";
        } else {
            $this->apiUri = "https://uat-ubaweb-newcdp.coach.com.cn";
        }
    }

    protected function getHeader($data): array
    {
        if(empty($data['task_id']) || empty($data['crowd_id'])) {
            throw new Error(100000, '子任务id和人群包id不能为空');
        }
        $taskId = $data['task_id'];
        $crowdId = $data['crowd_id'];
        $accessToken = strtoupper(md5($taskId.'&'.$crowdId));
        return ["headers" => [
            "Authorization" => "Bearer {$accessToken}",
            "Content-Type" => "application/json",
            "Cache-Control" => "no-cache"
        ]];
    }
}