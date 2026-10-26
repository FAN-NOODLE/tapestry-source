<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/11/16 11:44
 */

namespace App\Http\Controllers\App\V1\PartnerMember;


use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use App\Services\Wechat\WxOpenID;
use CrsSdk\Boss\MiniProgram;
use CrsSdk\Work\StoreWorker;
use CrsSdk\Work\PartnerMember;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Crypt;
use Illuminate\Support\Facades\Log;

class PersonalAppletLoginController extends BaseController
{
    protected $isRequireLogin = false;

    public function index(Request $request)
    {
        $params = $request->all();

        if (empty($params['code'])) {//code的不能为空
            error(652012);
        }

        if (empty($params['scene'])) {//二维码参数不能为空
            error(657004);
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

        //获取微信openid
        $result = (new WxOpenID($params['code'], $this->identifier))->get();

        $openId = $result['openid'];
        $sessionKey = $result['session_key'];

        if (!is_array($result)) {
            Log::error($result);
            error(657005);
        }

        // 授权返回结果错误
        if (empty($result['openid']) || empty($result['session_key'])) {
            Log::error('微信授权数据错误:' . var_export($result, true));
             error(657006);
        }

        //交互逻辑直接调用work模块
        $data = (new PartnerMember())->check(
            [
                'open_id' => $openId,
                'enterprise_id' => $params['enterprise_id']
            ]
        );

        //获得导购的头像和名称
        $workerInfo = (new StoreWorker())->detail(
            [
                'mem_uid' => $params['mem_uid'],
                'enterprise_id' => $params['enterprise_id'],
                'store_id' => $params['store_id'] ?? 0,
                'worker_id' => $params['worker_id']
            ]
        );
        if (empty($workerInfo)) {
            error(657007);
        }

        $tempToken = [
            'time' => time(),
            'open_id' => $openId,
            'sessionKey' => $sessionKey,
            'enterprise_id' => $params['enterprise_id'],
            'mem_uid' => $params['mem_uid'],
            'store_id' => $params['store_id'],
            'status' => $data['status'],
            'worker_id'=>$params['worker_id'],
        ];
        $data['cache_key'] = Crypt::encrypt(json_encode($tempToken));
        $data['worker']['avatar'] = $workerInfo['avatar'];
        $data['worker']['name'] = $workerInfo['name'];
        $data['worker']['mem_uid'] = $workerInfo['mem_uid'];
        $data['worker']['worker_id'] = $workerInfo['id'];
//        try {
//            $shopDetail = (new MiniProgram())->detail([
//                'enterprise_id' => $params['enterprise_id']
//            ]);
//        } catch (Error $error) {
//            if ($error->getCode() == 300100) {
//                notice(657012);
//            }
//        }

//        $data['worker']['enterprise_name'] =  $shopDetail['data']['name'] ?? '';
        //$enterpriseInfo = (new Enterprise())->detail(['enterprise_id' => $params['enterprise_id']]);
        //获得企业logo
        //$logoInfo  = (new Enterprise())->enterpriseMerchantSettingInfo(['enterprise_id' => $params['enterprise_id']]);
        //$data['worker']['enterprise_logo'] = $logoInfo['data']['logo'] ?? '';

        return success($data);
    }


}