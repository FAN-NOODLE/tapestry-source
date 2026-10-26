<?php


namespace App\Http\Controllers\V1\Message;

use App\Http\Controllers\V1\BaseController;
use App\Modules\Message\MessageFactory;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Validator;

class AccountController extends BaseController
{
    /**
     * @return \Illuminate\Http\JsonResponse
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidArgumentException
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function send()
    {
        Validator::make($this->request->all(), [
            'enterprise_id'           => 'required',
            'user_id'                 => 'required',
            'scene_type'              => 'required|in:1000,1001,1002,1003,1004,1005,1006,1007,1008,1009,1010,1011,2000,2001,2002,3000',
            'page_query'              => 'required',
            'msg'                     => 'required|array',
            'msg.spu_name'            => 'required_if:scene_type,1000,1001,1002,1003,1005,1006,1007,1008,1009,1010',
            'msg.start_at'            => 'required_if:scene_type,1004',
            'msg.activity_progress'   => 'required_if:scene_type,1005',
            'msg.end_at'              => 'required_if:scene_type,1005',
            'msg.luck_time'           => 'required_if:scene_type,1006',
            'msg.need_num'            => 'required_if:scene_type,1007,1008',
            'msg.group_progress'      => 'required_if:scene_type,1007,1008',
            'msg.refund_money'        => 'required_if:scene_type,1008',
            'msg.rest_time'           => 'required_if:scene_type,1009',
            'msg.sell_price'          => 'required_if:scene_type,1009,1010',
            'msg.activity_price'      => 'required_if:scene_type,1009,1010',
            'msg.haggle_progress'     => 'required_if:scene_type,1009,1010',
            'msg.live_title'          => 'required_if:scene_type,1011',
            'msg.start_time'          => 'required_if:scene_type,1011',
            'msg.wait_pay_money'      => 'required_if:scene_type,2000',
            'msg.create_at'           => 'required_if:scene_type,2000',
            'msg.order_no'            => 'required_if:scene_type,2000,2001,2002',
            'msg.order_goods_info'    => 'required_if:scene_type,2000,2001',
            'msg.remain_at'           => 'required_if:scene_type,2000',
            'msg.logistic_name'       => 'required_if:scene_type,2001',
            'msg.logistic_no'         => 'required_if:scene_type,2001',
            'msg.delivery_address'    => 'required_if:scene_type,2001',
            'msg.refund_apply_status' => 'required_if:scene_type,2002',
            'msg.refund_fee'          => 'required_if:scene_type,2002',
            'msg.refund_goods_info'   => 'required_if:scene_type,2002',
            'msg.score_change'        => 'required_if:scene_type,3000',
            'msg.score_total'         => 'required_if:scene_type,3000',
            'msg.type'                => 'required_if:scene_type,3000',
        ], [
            'enterprise_id.required' => '企业id必填',
            'user_id.required'       => '用户id必填',
            'scene_type.required'    => '类型必填',
            'scene_type.in'          => '类型可选值为[1000,1001,1002,1003,1004,1005,1006,1007,1008,1009,1010,1011,2000,2001,2002,3000]',
            'msg.required'           => 'msg必填',
            'msg.array'              => 'msg类型是数组',
            'page_query.required'    => '页面跳转参数必填',
        ])->validate();

        $input = $this->request->only(['enterprise_id', 'user_id', 'scene_type', 'msg', 'page_query']);

        $message = MessageFactory::getInstance(
            MessageFactory::CDP,
            $input['enterprise_id'],
            [$input['user_id']],
            $input['scene_type'],
            $input['msg'],
            $input);

        Log::info('公众号推送消息-公众号策略', $input);
        $data = $message->send();

        return success($data ?? []);
    }
}