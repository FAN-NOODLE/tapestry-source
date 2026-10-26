<?php

namespace App\Http\Controllers\Service\V1\SubscribeMessage;

use App\Exceptions\Error;
use App\Http\Controllers\Controller;
use App\Models\SubscribeMessageRecord;
use App\Models\SubscribeMessageTemplate;
use App\Services\EnterpriseService;
use App\Services\SubscribeMsgService;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class SendMessageController extends Controller
{
    /**
     * 发送订阅消息
     * @param Request $request
     * @return JsonResponse
     * @throws Exception
     */
    public function index(Request $request): JsonResponse
    {
        $input = $request->only(['enterprise_id', 'open_id', 'type', 'data']);
        if (empty($input['enterprise_id']) || empty($input['open_id']) || empty($input['type']) || empty($input['data'])) {
            throw new Error(10001, '参数错误');
        }
        $data = $input['data'];
        $info = SubscribeMessageRecord::query()
            ->where('enterprise_id', $input['enterprise_id'])
            ->where('openid', $input['open_id'])
            ->where('type', $input['type'])
            ->where('rel_id', $data['reserve_order_id'] ?? 0)
            ->where('send_num', 0)
            ->orderBy("id")
            ->first();
        if (empty($info)) {
            return successJson();
        }
        $tplInfo      = SubscribeMessageTemplate::query()
            ->where('enterprise_id', $input['enterprise_id'])
            ->where('pri_tmpl_id', $info->pri_tmpl_id)
            ->first();
        $dataNumArr   = explode(",", $tplInfo->data_num_str);
        $subscribeSvc = new SubscribeMsgService();
        list($appId, $appSecret) = $subscribeSvc->getMiniInfo($input['enterprise_id']);
        $msg              = [];
        $enterpriseConfig = EnterpriseService::getEnterpriseConfig($input['enterprise_id'], "store_name");
        $configInfo       = current($enterpriseConfig['data']);
        $delayTime        = 3600;
//        if ($configInfo['name']=="sw"){ //sw时间段改为半小时，需要放开注释，1小时需要注释掉
//            $delayTime=1800;
//        }
//        if ($configInfo['name'] == "coach") {//uat环境的coach针对三和一小程序的模板消息发送
//            $appId     = "wx3618ede4dbe080fa";
//            $appSecret = "REDACTED_FOR_SOURCE_SHARING";
//        }
        $reserve_end_time= $data['reserve_end_time']??0;
        if ($reserve_end_time==0){
            $reserve_end_time=$data['reserve_time'] + $delayTime;
        }
        switch ($input['type']) {
            case SubscribeMessageRecord::TO_STORE:
                $msg = [
                    'thing' . $dataNumArr[0] => ['value' => $data['store_name']],//预约门店
                    'thing' . $dataNumArr[1] => ['value' => $data['service_name']],//预约服务
                    'time' . $dataNumArr[2]  => ['value' => date("Y-m-d", $data['reserve_time'])],//到店日期
                    'time' . $dataNumArr[3]  => ['value' => date("H:i", $data['reserve_time']) . '~' . date("H:i", $reserve_end_time)],//到店时间
                    'thing' . $dataNumArr[4] => ['value' => '请您携带相关电子凭证出席，期待您的莅临！'],//温馨提示
                ];
                break;
            case SubscribeMessageRecord::SERVICE_CHANGE:
                $msg = [
                    'thing' . $dataNumArr[0] => ['value' => $data['store_name']],//预约门店
                    'thing' . $dataNumArr[1] => ['value' => $data['service_name']],//预约服务
                    'time' . $dataNumArr[2]  => ['value' => date("Y-m-d", $data['reserve_time'])],//到店日期
                    'time' . $dataNumArr[3]  => ['value' => date("H:i", $data['reserve_time']) . '~' . date("H:i", $reserve_end_time)],//到店时间
                    'thing' . $dataNumArr[4] => ['value' => '您预约信息已变更，如有疑问请联系专属顾问'],//温馨提示
                ];
                break;
            case SubscribeMessageRecord::SERVICE_COMMENT:
                $msg = [
                    'thing' . $dataNumArr[0] => ['value' => $data['service_name']],//预约服务
                    'thing' . $dataNumArr[1] => ['value' => $data['store_name']],//预约门店
                    'time' . $dataNumArr[2]  => ['value' => date("Y-m-d H:i:s", $data['reserve_time'])],//订单完结时间
                    'thing' . $dataNumArr[3] => ['value' => '请您点击此处进行评价'],//温馨提示
                ];
                break;
            default:
        }
        $params = [
            'touser'      => $input['open_id'],
            'template_id' => $info->pri_tmpl_id,
            'page'        => $info->page,
            'data'        => $msg,
        ];

        $res = $subscribeSvc->send($appId, $appSecret, $params);
        if ($res['errcode'] !== 0) {
            $info->fail_reason = $res['errcode'] . '|' . $res['errmsg'];
        } else {
            $info->send_num++;
        }
        $info->save();
        return successJson();
    }
}