<?php


namespace App\Http\Controllers\V1\Wx;


use App\Http\Controllers\V1\BaseController;
use App\WxServices\DynamicMessagesWxService;
use Illuminate\Http\Request;

class DynamicMsgCreateController extends BaseController
{
    public function index(Request $request)
    {
        $enterpriseId = $request->input('enterprise_id', 0);

        return (new DynamicMessagesWxService())->createActivityId($enterpriseId);
    }
}
