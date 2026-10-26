<?php

namespace App\Http\Controllers\V1\Admin;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Boss\Enterprise;
use CrsSdk\Work\Auth;

class JsApiConfigController extends BaseController
{
    /**
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $enterprise = (new Enterprise())->detail(['enterprise_id' => $this->enterpriseId])['data'] ?? [];

        $domain = $enterprise['domain'] ?? '';

        $url = $this->request->input('url', '');

        $data = (new Auth())->agentConfig([
                'url'           => $url,
                'domain'        => $domain,
                'appIdentifier' => 'crs_store_assistant'
            ]) ?? [];

        return success($data);
    }
}
