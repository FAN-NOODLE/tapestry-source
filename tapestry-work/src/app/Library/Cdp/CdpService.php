<?php

namespace App\Library\Cdp;

class CdpService
{
    private $client;

    /**
     * CdpRequests constructor.
     */
    public function __construct()
    {
        $this->client = new CdpRequests();
    }

    /**
     * @param array $params
     *
     * @return mixed
     */
    public function userAct(array $params)
    {
        return $this->client->doPost($params);
    }
}
