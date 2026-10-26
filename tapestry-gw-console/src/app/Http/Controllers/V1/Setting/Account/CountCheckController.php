<?php
namespace App\Http\Controllers\V1\Setting\Account;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroAuthority\Account;
use Illuminate\Support\Facades\Log;

/**
 * 检测企业子账号数量
 * Class CountCheckController
 * @package App\Http\Controllers\V1\account
 */
class CountCheckController extends BaseController
{
    public function index()
    {
        $input['enterprise_id'] = $this->enterpriseId;

        (new Account())->checkAccountNum($input);

        return success([]);
    }

}