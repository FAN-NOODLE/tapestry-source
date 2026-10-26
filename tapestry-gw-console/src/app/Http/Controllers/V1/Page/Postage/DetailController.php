<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2019/4/3
 * Time: 16:43
 */

namespace App\Http\Controllers\V1\Page\Postage;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\GlobalPostage;

class DetailController extends BaseController
{
    public function index()
    {
        try {
            $list = (new GlobalPostage())->lists(['enterprise_id' => $this->enterpriseId]);

            $return = [
                'enterprise_id' => $list['enterprise_id'] ?? $this->enterpriseId,
                'rule'          => $list['rule'] ?? 0,
                'on_off'        => $list['on_off'] ?? 2,
            ];
            return success($return);
        } catch (\Exception $e) {
            error(450000, [$e->getMessage()]);
        }
    }
}
