<?php

namespace App\Sdk\Active\CRM\KsCore;

use App\Library\Util\Logger;

abstract class Base
{
    protected $errStr = '';

    protected $queryParam = [];

    protected $header = [];

    protected $api = '';

    protected $throttle = false;

    protected $type = '';

    protected $tokenUrl = false;

    public function __construct()
    {
        $this->setQueryParam();
    }

    protected abstract function setQueryParam();

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
        if (config('app.env') == 'prod') {
            $prefix = "";

            return $prefix . substr($this->api, strlen($prefix));
        }
        return $this->api;

    }

    public function getQueryParam(): array
    {
        if ($this->type == 'json') {
            return ['body' => json_encode($this->queryParam)];
        } else if ($this->type == 'form_data') {
            return ['form_params' => $this->queryParam];
        } else {
            return ['params' => $this->queryParam];
        }

    }

    public function getHeader(): array
    {
        return $this->header;
    }

    public function getErrorStr(): string
    {
        return $this->errStr;
    }

    public function getThrottle(): bool
    {
        return $this->throttle;
    }

    public function setType()
    {
        return $this->type='form_data';
    }


}