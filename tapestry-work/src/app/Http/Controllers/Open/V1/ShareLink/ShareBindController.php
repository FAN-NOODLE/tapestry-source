<?php
/**
 * 分享链路 - 创建绑定/解除绑定
 */

namespace App\Http\Controllers\Open\V1\ShareLink;

use App\Http\Controllers\Open\V1\BaseController;
use App\Jobs\Message\ShareLinkMessageAdapter;
use App\Jobs\ShareLinkUpdateCustomerStatisticJob;
use App\Library\Ddc\DdcMessage;
use App\Library\Uc\Member;
use App\Models\PartnerEnterpriseSettings;
use App\Models\PartnerMember;
use App\Models\StaffShare;
use App\Services\Events\ShareLinkChangeEventService;
use App\Services\StoreWorkerService;
use App\Services\UserShareService;
use Illuminate\Support\Facades\Log;
use App\Services\PartnerMemberService;

class ShareBindController extends BaseController
{

    /**
     * 员工状态：1=在职，2=离职
     */
    const IN_SERVICE = 1;
    const OUT_SERVICE = 2;

    //链路绑定身份标记类型 1、普通用户  2、导购  10、合伙人
    const SHARE_BIND_TYPE_USER = 1;
    const SHARE_BIND_TYPE_GUIDE = 2;
    const SHARE_BIND_TYPE_PARTNER = 10;

    public function index()
    {
        Log::info('shareBindParams', $this->request->all());
        $input = $this->getParams();

        $shareUid = $input['share_uid'] ?? '';
        $workerInfo = [];
        $memUid = $input['mem_uid'] ?? '';
        $userType = self::SHARE_BIND_TYPE_USER;

        if (!empty($shareUid) && $input['source'] == 'scrm') {
            $qyDomain = $this->request->get('enterprise_info')['domain'];
            $workerDetail = (new Member())->detailByUserId($qyDomain, 'crs-shop', ['memUserid' => $shareUid]);
            $memUid = $workerDetail['memUid'] ?? $memUid;
        }

        if (!empty($shareUid) && $input['source'] == 'isg') {
            $memUid = $shareUid;
        }

        $partnerMemberService = new PartnerMemberService();
        if (!empty($shareUid) && $input['source'] == 'shop') {
            $input['scene_user_id'] = $shareUid;
        }
        if (!empty($input['scene_user_id'])) {
            // shareUid 为分享者商城userId
            $partnerInfo = $partnerMemberService->getPartnerInfoByUserId($input['enterprise_id'], $input['scene_user_id'], PartnerMember::PARTNER_STATUS_ENABLED);
            if (!empty($partnerInfo)) {
                $shareUid = $partnerInfo['id'];
            }
        }

        // 获取登陆者身份
        //查询是否合伙人
        if (!empty($memUid) || !empty($shareUid)) {
            $workerInfo = $partnerMemberService->getPartnerInfo([
                'enterprise_id' => $input['enterprise_id'],
                'id' => empty($memUid) ? $shareUid : $memUid,
                'status' => PartnerMember::PARTNER_STATUS_ENABLED
            ]);
        }
        Log::info("查询合伙人信息", ['partner' => $workerInfo]);
        if (empty($workerInfo)) {
            if (!empty($input['worker_id'])) {
                $workerInfo = StoreWorkerService::detail([
                    'enterprise_id' => $input['enterprise_id'],
                    'worker_id' => $input['worker_id'],
                    'worker_status' => self::IN_SERVICE
                ]);
            } elseif (!empty($memUid) || !empty($input['scene_user_id'])) {
                $workerInfo = StoreWorkerService::getWorkerByShareLink($input['enterprise_id'], strval($memUid), $input['scene_user_id'], $input['store_id']);
            }
            if (!empty($workerInfo)) {
                $userType = self::SHARE_BIND_TYPE_GUIDE;
            }
            Log::info("查询导购信息", ['worker' => $workerInfo]);
        } else {
            //查询二级分销是否开启
            //检查是否开启了合伙人分销，是否开启了全员分销
            $checkSecondSales = PartnerEnterpriseSettings::query()
                ->where('enterprise_id', $input['enterprise_id'])
                ->where('partner_status', PartnerEnterpriseSettings::PARTNER_STATUS_ENABLE)
                ->first();
            if (empty($checkSecondSales)) {
                //没开启不用调接口绑定
                $userType = self::SHARE_BIND_TYPE_USER;
                $workerInfo = [];
            } else {
                $userType = self::SHARE_BIND_TYPE_PARTNER;
            }
        }
        $memUid = $workerInfo['mem_uid'] ?? '';
        $sceneUserId = isset($workerInfo['user_id'])
            ? ($workerInfo['user_id'] ?: $input['scene_user_id'])
            : 0;

        $workerId = $workerInfo['id'] ?? $workerInfo['worker_id'] ?? 0;
        $storeId = $workerInfo['store_id'] ?? 0;

        //链路绑定
        if (($workerId || $sceneUserId) && $input['user_id'] > 0) {
            $condition = [
                'enterprise_id' => $input['enterprise_id'],
                'user_id' => $input['user_id']
            ];
            $beforeBinding = (new StaffShare())->getStaffShareInfoByCondition($condition);
            (new UserShareService($input['enterprise_id'], $storeId, $input['user_id'], strval($memUid), $sceneUserId ?? 0, $workerId, $userType))
                ->setUserShare();
            $afterBinding = (new StaffShare())->getStaffShareInfoByCondition($condition);
            if (!empty($beforeBinding) || !empty($afterBinding)) {
                $this->shareLinkEventNotify($beforeBinding, $afterBinding);
            }
        }

        return $this->success([
            'mem_uid' => $memUid,
            'scene_user_id' => $sceneUserId,
            'worker_id' => $workerId,
            'store_id' => $storeId,
            'user_type' => $userType
        ]);
    }

    /**
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    private function getParams()
    {
        $input = $this->request->only([
            'store_id',
            'mem_uid',
            'worker_id',
            'user_id',
            'share_uid',
            'source',
            'scene_user_id',
            'enterprise_id'
        ]);

        $rules = [
            'store_id' => 'integer',
            'user_id' => 'required',
        ];
        $messages = [
            'store_id.*' => '门店id参数不合法',
            'user_id.*' => '用户id参数不合法'
        ];

        Log::info('分享链路绑定数据ShareBindController', [$this->request->input()]);

        $this->validate($this->request, $rules, $messages);

        return $input;
    }

    /**
     * 用户绑定关系变更通知
     *
     * @param $beforeBinding
     * @param $afterBinding
     * @throws \Exception
     */
    public function shareLinkEventNotify($beforeBinding, $afterBinding)
    {
        Log::info("分享链路关系变更通知数据", func_get_args());
        $method = '';
        $list = [];
        if (!empty($beforeBinding) && !empty($afterBinding) && $beforeBinding['type'] != 10 && $afterBinding['type'] != 10) {
            // 绑定关系变更 导购->导购
            $method = DdcMessage::GUIDE_USER_BINDING_UPDATE;
            $list = [
                'enterprise_id' => $afterBinding['enterprise_id'],
                'third_user_id' => $afterBinding['user_id'],
                'curr_worker_id' => $afterBinding['worker_id'],
                'curr_worker_mem_uid' => $afterBinding['mem_uid'],
                'old_worker_id' => $beforeBinding['worker_id'],
                'old_worker_mem_uid' => $beforeBinding['mem_uid'],
                'bind_time' => intval($afterBinding['created_at']),
            ];
            if ($beforeBinding['worker_id'] != $afterBinding['worker_id']) {
                dispatch(new ShareLinkUpdateCustomerStatisticJob(ShareLinkMessageAdapter::updateMessage($beforeBinding, $afterBinding)));
            }
        } else if (!empty($beforeBinding) && $beforeBinding['type'] != 10) {
            // 解除绑定关系 导购解绑
            $method = DdcMessage::GUIDE_USER_BINDING_DELETE;
            $list = [
                'enterprise_id' => $beforeBinding['enterprise_id'],
                'third_user_id' => $beforeBinding['user_id'],
                'curr_worker_id' => $beforeBinding['worker_id'],
                'curr_worker_mem_uid' => $beforeBinding['mem_uid'],
                'bind_time' => intval($beforeBinding['created_at']),
            ];
            dispatch(new ShareLinkUpdateCustomerStatisticJob(ShareLinkMessageAdapter::deleteMessage($beforeBinding, time())));
        } else if (!empty($afterBinding) && $afterBinding['type'] != 10) {
            // 新增绑定关系 初次绑定导购
            $method = DdcMessage::GUIDE_USER_BINDING_ADD;
            $list = [
                'enterprise_id' => $afterBinding['enterprise_id'],
                'third_user_id' => $afterBinding['user_id'],
                'curr_worker_id' => $afterBinding['worker_id'],
                'curr_worker_mem_uid' => $afterBinding['mem_uid'],
                'bind_time' => intval($afterBinding['created_at']),
            ];
            dispatch(new ShareLinkUpdateCustomerStatisticJob(ShareLinkMessageAdapter::createMessage($afterBinding)));
        }
        (new ShareLinkChangeEventService)->pushDdcMessage($method, [$list]);
    }
}