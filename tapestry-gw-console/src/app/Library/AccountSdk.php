<?php

namespace App\Library;


use CrsSdk\MicroAuthority\Account;

class AccountSdk extends Account
{
    /**
     * @param array $params
     * @return mixed
     * @throws \Exception
     */
    public function accountAccessControl($params = [])
    {
        return $this->httpPost('micro-authority/v1/account/account-access-control', $params);
    }
}
