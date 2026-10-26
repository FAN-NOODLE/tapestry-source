<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/10/19
 * Time: 11:20 AM
 */

namespace App\Http\Controllers\Service\V1\Auth;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\JsApi;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class AgentConfigController extends BaseController
{
    /**
     * js config
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $domain        = $request->input('domain', '');
        $appIdentifier = $request->input('appIdentifier', '');
        $url           = $request->input('url', '') ?: route('service-agent-config');


        try {
            $data = (new JsApi())->agentConfig($domain, $appIdentifier, [
                'url' => 'https://crs-test.vchangyi.com/f6691183139417287f27dcfcef01489d/myapp/organization/team-manager/organize-new/list'
            ]);
        } catch (\Exception $exception) {
            Log::info('JS API config 签名生成错误' . $exception->getMessage());
            $data = [];
        }
        return success($data);
    }
}
