<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/11/25 11:12
 */

namespace App\Http\Controllers\App\V1\Common;


use App\Exceptions\Error;
use App\Http\Controllers\Controller;
use CrsSdk\Boss\Enterprise;
use Illuminate\Http\Request;

class GetEnterpriseHashController extends Controller
{
        public function index(Request $request)
        {
            $params = $request->all();

            if (empty($params['scene'])) {//二维码参数不能为空
                throw new Error(110002, '二维码参数不能为空');
            }
            $decryptData = decryptShortLink($params['scene']);

            if (empty($decryptData['args'])) {
                $decryptData = decryptShortLinkByTest($params['scene']);
                if (empty($decryptData['args'])) {
                    notice(657008);
                }
            }
            if (!empty($decryptData['args'])) {
                foreach ($decryptData['args'] as $v) {
                    $params[$v['name']] = $v['value'];
                }
            }
            $enterpriseInfo = (new Enterprise())->detail(['enterprise_id' => $params['enterprise_id']]);
            $enterprise_hash = $enterpriseInfo['data']['hash'] ?? '';

            return success(['enterprise_hash'=>$enterprise_hash]);
        }
}