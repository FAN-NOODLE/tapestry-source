<?php

namespace App\Http\Controllers\Service\V1\Activity\Scratch;

use App\Http\Controllers\Controller;
use App\Models\Activity;
use App\Models\ActivityPrize;
use App\Models\ActivityRecord;
use App\Models\ActivityShareLog;
use App\Models\ActivityUserView;
use App\Models\Scrm\Members;
use App\Models\ScrmEnterprise\WxUserCorpExternal;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Validation\ValidationException;
use SimpleSoftwareIO\QrCode\BaconQrCodeGenerator;

class IndexController extends Controller
{

    /**
     * 获取活动信息
     * @return JsonResponse
     * @throws ValidationException
     */
    public function Info(): JsonResponse
    {
        $rules   = [
            'enterprise_id' => 'required',
            'activity_id'   => 'required',
            'openid'        => 'required'
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'activity_id.required'   => '活动ID不能为空',
            'openid.required'        => '用户身份错误'
        ];
        $params  = $this->validate($this->request, $rules, $message);

        $enterprise_id = $params['enterprise_id'];
        $id            = $params['activity_id'];
        try {
            $info = Activity::query()->where(['id' => $id, 'enterprise_id' => $enterprise_id])->first();
        } catch (Exception $e) {
            error($e->getMessage());
        }
        $activityPrizeList = ActivityPrize::query()->where([
            'activity_id'   => $id,
            'enterprise_id' => $params['enterprise_id'],
        ])->get([
            'id',
            'enterprise_id',
            'activity_id',
            'prize_short_name',
            'prize_level',
        ]);
        if (empty($info)) {
            return successJson();
        }
        $result      = $info->toArray();
        $prize_names = '';
        if ($activityPrizeList->count() > 0) {
            $tempList = $activityPrizeList->toArray();
            foreach ($tempList as $value) {
                $prize_names .= $value['prize_short_name'] . '/';
            }
            $prize_names = substr($prize_names, 0, -1);
        }

        $result['prize_names'] = $prize_names;
        $start_time_str        = date("Y年m月d日", $result['start_time']);
        $end_time_str          = date("m月d日", $result['end_time']);
        $result['time_str']    = $start_time_str . '-' . $end_time_str;

        $nums = ActivityRecord::query()->where([
            'activity_id'   => $id,
            'enterprise_id' => $params['enterprise_id'],
            'openid'        => $params['openid']
        ])->count();
        if ($result['start_time'] > time()) {
            $result['activity_status'] = Activity::StatusNotStarted;
        }
        if ($result['end_time'] < time()) {
            $result['activity_status'] = Activity::StatusEnd;
        }
        $canScratch = true;
        if ($result['activity_num'] <= $nums || $result['activity_status'] != Activity::StatusRun) {
            $canScratch = false;
        }
        $result['can_scratch'] = $canScratch;
        return successJson($result);
    }

    /**
     * 用户参与记录
     * @return JsonResponse
     * @throws ValidationException
     */
    public function userView(): JsonResponse
    {
        $params  = $this->request->all();
        $rules   = [
            'enterprise_id' => 'required',
            'activity_id'   => 'required',
            'openid'        => 'required',
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'activity_id.required'   => '活动ID不能为空',
            'openid.required'        => '用户openid不能为空',
        ];
        $this->validate($this->request, $rules, $message);

        $userName = '';
        $unionId  = $params['union_id'];
        if (!empty($unionId)) {
            $userInfos = WxUserCorpExternal::getExternalUserInfoByUnionIdNoCode($params['enterprise_id'], $unionId);
            $userInfo  = current($userInfos);
            if (!empty($userInfo)) {
                $userName = $userInfo['nick'];
            }
        }
        $data = [
            'activity_id'   => $params['activity_id'],
            'enterprise_id' => $params['enterprise_id'],
            'openid'        => $params['openid'],
            'unionid'       => $unionId,
            'user_name'     => $userName,
            'source'        => $params['source'],
            'share_num'     => 0
        ];
        try {
            $info = ActivityUserView::query()->where(['openid' => $params['openid']])->first();
            if (!empty($info) && empty($info['source'])) {
                $info->update(['source' => $params['source']]);
            }
            if (!empty($info) && $userName != $info['user_name']) {
                $info->update(['user_name' => $userName]);
            }
            if (empty($info)) {
                ActivityUserView::query()->create($data);
            }
            if ($params['source'] == "user_share") { //记录用户的打开分享链接的日志
                $data = [
                    'enterprise_id' => $params['enterprise_id'],
                    'activity_id'   => $params['activity_id'],
                    'openid'        => $params['openid'],
                    'source'        => 1
                ];
                ActivityShareLog::query()->create($data);
            }
        } catch (Exception $e) {
            error($e->getMessage());
        }
        return successJson();
    }

    /**
     * 导购核销券
     * @param Request $request
     * @return JsonResponse
     * @throws ValidationException
     */
    public function WriteOffCoupon(Request $request): JsonResponse
    {
        $params  = $request->all();
        $rules   = [
            'enterprise_id' => 'required',
            'mobile'        => 'required',
            'code'          => 'required',
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'mobile.required'        => '导购身份错误',
            'code.required'          => '核销码错误',
        ];
        $this->validate($request, $rules, $message);
        $memberInfo = Members::query()->where(['enterprise_id' => $params['enterprise_id'], 'mobile' => $params['mobile']])->first();
        if (empty($memberInfo)) {
            error("导购不存在");
        }
        $record = ActivityRecord::query()->where(['enterprise_id' => $params['enterprise_id'], 'activity_code' => $params['code']])->first();
        if (empty($record)) {
            error("无效的核销码");
        }
        if ($record['is_check']) {
            error("该券已被核销了");
        }
        $checkInfo = ActivityRecord::query()->where(['enterprise_id' => $params['enterprise_id'], 'activity_id' => $record['activity_id'], 'openid' => $record['openid'], 'is_check' => ActivityRecord::IsCheckYes])->first();
        if (!empty($checkInfo)) {
            error("该用户已核销过优惠券，无法再次核销");
        }
        $activityInfo = Activity::query()->where(['id' => $record['activity_id'], 'enterprise_id' => $params['enterprise_id']])->first();
        if (empty($activityInfo)) {
            error("活动不存在");
        }
        $nowTime = time();
        if ($activityInfo['check_time'] < $nowTime) {
            error("该券已过期");
        }
        $userName = $record['user_name'];
        if (empty($record['user_name'])) {
            $userInfos = WxUserCorpExternal::getExternalUserInfoByUnionIdNoCode($params['enterprise_id'], $record['unionid']);
            $userInfo  = current($userInfos);
            if (!empty($userInfo)) {
                $userName = $userInfo['nick'];
            }
        }
        $updateData = [
            'is_check'    => 1,
            'check_time'  => $nowTime,
            'version'     => $record['version'] + 1,
            'member_id'   => $memberInfo['id'],
            'member_name' => $memberInfo['name'],
            'user_name'   => $userName
        ];
        try {
            ActivityRecord::query()->where(['id' => $record['id'], 'version' => $record['version']])->update($updateData);
        } catch (Exception $e) {
            error("核销失败");
        }
        return successJson();
    }

    /**
     * 刮奖
     *
     * @throws ValidationException
     * @author  user <qiuapeng921@vchangyi.com>
     * @date    2022/8/15 10:00
     */
    public function ScratchPrize(): JsonResponse
    {
        $rules   = [
            'enterprise_id' => 'required',
            'activity_id'   => 'required',
            'openid'        => 'required',
            'unionid'       => 'required',
            'record_id'     => 'sometimes|required'
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'activity_id.required'   => '活动ID不能为空',
            'openid.required'        => '用户标识不能为空',
            'unionid.required'       => '用户标识不能为空',
            'record_id.required'     => '抽奖失败',
        ];
        $params  = $this->validate($this->request, $rules, $message);

        $activity = Activity::query()->where([
            'id'            => $params['activity_id'],
            'enterprise_id' => $params['enterprise_id'],
        ])->first();
        // 判断活动是否存在,或者是否进行中
        if (is_null($activity) || $activity['activity_status'] != Activity::StatusRun) {
            error("活动不存在");
        }


        // 获取用户刮卡次数
        $activityRecordCount = ActivityRecord::query()->where([
            'activity_id'   => $params['activity_id'],
            'enterprise_id' => $params['enterprise_id'],
            'openid'        => $params['openid'],
        ])->count();

        // 查看用户是否有刮奖资格
        if ($activityRecordCount >= $activity['activity_num']) {
            $responseData = [
                'coupon' => [],
                'times'  => -1,
            ];
            return successJson($responseData);
        }
        $result = [];
        if (!empty($params['record_id'])) {
            $activityRecord             = ActivityRecord::query()->where(['id' => $params['record_id']])->withTrashed()->first();
            $activityRecord->deleted_at = 0;
            $activityRecord->save();
            $result              = ActivityPrize::query()->where(['id' => $activityRecord['prize_id']])
                ->get(['id', 'enterprise_id',
                       'activity_id',
                       'prize_name',
                       'prize_short_name',
                       'prize_level'])->first();
            $result['record_id'] = $activityRecord['id'];
        } else {
            $activityRecordLast = ActivityRecord::query()->where([
                'activity_id'   => $params['activity_id'],
                'enterprise_id' => $params['enterprise_id'],
                'openid'        => $params['openid'],
            ])->where('deleted_at', '>', 0)->withTrashed()->first();
            if (!empty($activityRecordLast)) { //查询用户上次未完成刮卡的奖品，存在，继续返回上次的中奖信息
                $result              = ActivityPrize::query()->where(['id' => $activityRecordLast['prize_id']])
                    ->get(['id', 'enterprise_id',
                           'activity_id',
                           'prize_name',
                           'prize_short_name',
                           'prize_level'])->first();
                $result['record_id'] = $activityRecordLast['id'];
            } else {
                // 获取活动关联的奖品
                $activityPrizeList = ActivityPrize::query()->where([
                    'activity_id'   => $params['activity_id'],
                    'enterprise_id' => $params['enterprise_id'],
                ])->get([
                    'id',
                    'enterprise_id',
                    'activity_id',
                    'prize_name',
                    'prize_short_name',
                    'prize_level',
                    'prize_ratio'
                ]);


                if (is_null($activityPrizeList) || $activityPrizeList->count() == 0) {
                    error("活动奖品异常");
                }

                // 中奖概率
                $probability = rand(1, 100);

                $end = 0;

                $tmpList = $activityPrizeList->toArray();

                // 计算中奖
                foreach ($tmpList as $item) {
                    $start = $end + 1;
                    $end   = ($start - 1) + $item['prize_ratio'];
                    if ($probability >= $start && $probability <= $end) {
                        $result = $item;
                        unset($result['prize_ratio']);
                    }
                }

                // 生成活动码券码
                $activityCode = idGenerate();
                // 生成二维码
                $qrCodeGenerator = new BaconQrCodeGenerator();
                $qrcode          = base64_encode($qrCodeGenerator->format("png")->size(220)->generate($activityCode));

                // 奖品Id存在则中奖
                $status = $result ? ActivityRecord::StatusYes : ActivityRecord::StatusNo;

                // 获取用户昵称
                $userView = ActivityUserView::query()->where([
                    'activity_id'   => $params['activity_id'],
                    'enterprise_id' => $params['enterprise_id'],
                    'openid'        => $params['openid'],
                ])->first();
                $userName = "";
                if ($userView && $userView['user_name']) {
                    $userName = $userView['user_name'];
                }

                // 插入抽奖记录
                $saveData = [
                    'enterprise_id' => $params['enterprise_id'],
                    'activity_id'   => $params['activity_id'],
                    'openid'        => $params['openid'],
                    'unionid'       => $params['unionid'],
                    'user_name'     => $userName,
                    'qr_code'       => $qrcode,
                    'activity_code' => $activityCode,
                    'prize_id'      => $result['id'],
                    'status'        => $status,
                    'created_at'    => time(),
                    'updated_at'    => time(),
                    'deleted_at'    => time()
                ];
                $recordId            = ActivityRecord::query()->insertGetId($saveData);
                $result['record_id'] = $recordId;
            }
        }

        $responseData = [
            'coupon' => $result,
            'times'  => $activity['activity_num'] - ($activityRecordCount + 1),
        ];

        return successJson($responseData);
    }


    /**
     * 券列表
     *
     * @throws ValidationException
     * @author  user <qiuapeng921@vchangyi.com>
     * @date    2022/8/15 10:02
     */
    public function CouponList(): JsonResponse
    {
        $page = $this->request->post('page', 1);
        $size = $this->request->post('size', 10);
        if ($size > 999) {
            $size = 10;
        }

        $rules   = [
            'enterprise_id' => 'required',
            'openid'        => 'sometimes|required',
            'mobile'        => 'sometimes|required',
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'openid.required'        => '用户标识不能为空',
            'mobile.required'        => '导购标识不能为空'
        ];
        $params  = $this->validate($this->request, $rules, $message);

        // 获取用户刮卡记录
        $query = ActivityRecord::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->with(['activity:id,start_time,check_time', 'prize:id,prize_name,prize_level']);

        if (isset($params['openid'])) {
            $query->where('openid', $params['openid']);
        }

        if (isset($params['mobile'])) {
            $memberInfo = Members::query()->where([
                'enterprise_id' => $params['enterprise_id'],
                'mobile'        => $params['mobile']
            ])->first();
            if (!is_null($memberInfo)) {
                $query->where('member_id', $memberInfo['id'])->orderByDesc('check_time');
            }
        }

        $result = $query->paginate($size, [
            'id',
            'prize_id',
            'activity_id',
            'qr_code',
            'activity_code',
            'is_check',
            'user_name',
            'check_time',
            'created_at'
        ], '', $page);


        $status = 0;
        // 判断用户是否已核销过
        $isCheck = false;
        if (isset($params['openid'])) {
            // 获取用户已核销券
            $activityRecord = ActivityRecord::query()->where([
                'enterprise_id' => $params['enterprise_id'],
                'openid'        => $params['openid'],
                'is_check'      => ActivityRecord::IsCheckYes
            ])->first();
            if (!is_null($activityRecord)) {
                $isCheck = true;
                $status  = 1;
            }
        }

        $checkData = [];
        // 未核销取出最大奖品id
        if ($isCheck === false) {
            $checkData = collect($result->items())->sortBy('prize.prize_level')->first();
        }
        $nowTime = time();
        foreach ($result->items() as &$item) {
            if ($item['check_time'] > 0) {
                $item['check_time_string'] = date("Y年m月d日 H:i", $item['check_time']);
            }
            if ($status == 0 && $item['activity']['check_time'] < $nowTime) {
                $status = 2;
            }
            if ($item['activity'] && $item['activity']['start_time'] > 0 && $item['activity']['check_time'] > 0) {
                $start_time_str                       = date("Y年m月d日", $item['activity']['start_time']);
                $check_time_str                       = date("Y年m月d日", $item['activity']['check_time']);
                $item['activity']['expired_time']     = $start_time_str . '-' . $check_time_str;
                $start_time_string                    = date("Y/m/d", $item['activity']['start_time']);
                $check_time_string                    = date("Y/m/d", $item['activity']['check_time']);
                $item['activity']['expired_time_str'] = $start_time_string . '-' . $check_time_string;

            }
        }

        $data = [
            'meta'     => [
                'total' => $result->total(),
                'page'  => intval($page),
                'size'  => intval($size)
            ],
            'data'     => $result->items(),
            'is_check' => $isCheck,
            'check_id' => $checkData ? $checkData->id : 0,
            'status'   => $status
        ];

        return successJson($data);
    }

    /**
     * 检查用户的核销状态
     *
     * @return JsonResponse
     * @throws ValidationException
     */
    public function CheckCoupon(): JsonResponse
    {
        $rules   = [
            'enterprise_id' => 'required',
            'openid'        => 'required',
            'id'            => 'required',
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'openid.required'        => '用户标识不能为空',
            'id.required'            => '无效核销码'
        ];
        $params  = $this->validate($this->request, $rules, $message);

        $info = ActivityRecord::query()->where(['id' => $params['id'], 'openid' => $params['openid']])->get(["is_check"])->first();
        if (empty($info)) {
            error("核销码错误");
        }
        $result = $info->toArray();
        return successJson($result);
    }

    /**
     * 用户的分享上报记录
     * @return JsonResponse
     * @throws ValidationException
     */
    public function ShareInfo(): JsonResponse
    {
        $rules    = [
            'enterprise_id' => 'required',
            'openid'        => 'required',
            'activity_id'   => 'required',
        ];
        $message  = [
            'enterprise_id.required' => '企业ID不能为空',
            'openid.required'        => '用户标识不能为空',
            'activity_id.required'   => '无效活动信息'
        ];
        $params   = $this->validate($this->request, $rules, $message);
        $userInfo = ActivityUserView::query()->where(['openid' => $params['openid'], 'enterprise_id' => $params['enterprise_id']])->first();
        $userInfo->increment("share_num");
        $data = [
            'enterprise_id' => $params['enterprise_id'],
            'activity_id'   => $params['activity_id'],
            'openid'        => $params['openid']
        ];
        ActivityShareLog::query()->create($data);
        return successJson();
    }
}