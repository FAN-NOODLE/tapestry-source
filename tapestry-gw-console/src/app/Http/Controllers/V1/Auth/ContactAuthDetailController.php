<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/12/3
 * Time: 11:55 AM
 */

namespace App\Http\Controllers\V1\Auth;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Boss\Auth;

class ContactAuthDetailController extends BaseController
{
    /**
     * 外部联系人secret
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $result = (new Auth())->contactAuthDetail(['enterprise_id' => $this->enterpriseId]);

        $this->writeOperatorLog("通讯录授权");

        return success($result['data'] ?? []);
    }
}
