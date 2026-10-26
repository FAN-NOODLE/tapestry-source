<?php

namespace App\Sdk\Active\CRM\CoachCore;

use App\Sdk\Active\CRM\CoachCore\Auth\TokenCache;
use App\Sdk\Active\CrmHelper;

abstract class Base
{
    protected $errStr = '';

    protected $queryParam = [];

    protected $header = [];

    protected $api = '';

    protected $throttle = false;

//    protected $type = 'json';

    protected $tokenUrl = false;

    public function __construct()
    {
        $this->setQueryParam();
//        $this->setHeader();
    }

    public function getApiUat(): string
    {
        if ($this->tokenUrl) {
            return $this->api;
        }

        if (config('app.env') == 'prod') {
            return $this->api;
        }
        $prefix = "";
        return $prefix . substr($this->api, strlen($prefix));

    }

    public function getApiProd(): string
    {
        if ($this->tokenUrl) {
            return $this->api;
        }
        if (config('app.env') == 'prod') {
            $prefix = "";

            return $prefix . substr($this->api, strlen($prefix));
        }
        return $this->api;

    }

    public function getTokenApiUat(): string
    {
        if ($this->tokenUrl) {
            return $this->api;
        }

        if (config('app.env') == 'prod') {
            return $this->api;
        }
        $prefix = "";
        return $prefix . substr($this->api, strlen($prefix));


    }

    public function getTokenApiProd(): string
    {
        if ($this->tokenUrl) {
            return $this->api;
        }

        if (config('app.env') == 'prod') {
            return $this->api;
        }
        $prefix = "";
        return $prefix . substr($this->api, strlen($prefix));

    }

    public function getQueryParam($requestType = 'json'): array
    {
        if ($requestType == 'json') {
            return ['body' => json_encode($this->queryParam)];
        } else if ($requestType == 'form_data') {
            return ['form_params' => $this->queryParam];
        } else if ($requestType == 'xml') {
            return ['body' => $this->queryParam];
        } else {
            return ['data' => $this->queryParam];
        }

    }

    protected abstract function setQueryParam();

    public function handleQueryParam(): array
    {
        return ['data' => $this->queryParam];

    }

    public function getHeader(): array
    {
        return $this->header;
    }

    public function setHeader($enterpriseId, $xCountryCode = 'CN')
    {
        $headers = [
            'Authorization' => "Bearer " . (new TokenCache($enterpriseId))->getAccessToken(),
            'X-Country-Code' => $xCountryCode,
            'Content-Type' => 'application/json',
        ];
        if ($this->throttle) {
            $headers['ruid'] = json_encode(['uuid' => CrmHelper::createUid()]);

        }
        $this->header = ['headers' => $headers, 'http_errors' => false];

    }

    public function getErrorStr(): string
    {
        return $this->errStr;
    }

    public function getThrottle(): bool
    {
        return $this->throttle;
    }

    public function setHeaderCore($enterpriseId)
    {
        $headers = [
            'Content-Type' => 'application/json',
        ];
        if ($this->throttle) {
            $headers['ruid'] = json_encode(['uuid' => CrmHelper::createUid()]);

        }
        $this->header = ['headers' => $headers];

    }


}