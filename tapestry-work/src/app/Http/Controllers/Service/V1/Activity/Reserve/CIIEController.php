<?php

namespace App\Http\Controllers\Service\V1\Activity\Reserve;

use App\Http\Controllers\Controller;
use App\Library\EasyWechatCacheRedis\MiniProgram;
use App\Models\Activity;
use App\Models\ActivityReserveRecord;
use App\Models\ActivityReserveTimeConfig;
use App\Models\ActivityUserView;
use App\Models\Scrm\Members;
use App\Models\ScrmEnterprise\WxUserCorpExternal;
use CrsSdk\Boss\EnterpriseWorkPlugin;
use EasyWeChat\Kernel\Exceptions\HttpException;
use EasyWeChat\Kernel\Exceptions\InvalidArgumentException;
use EasyWeChat\Kernel\Exceptions\InvalidConfigException;
use EasyWeChat\Kernel\Exceptions\RuntimeException;
use Exception;
use GuzzleHttp\Client;
use Illuminate\Http\JsonResponse;
use Illuminate\Support\Facades\Redis;
use Illuminate\Validation\ValidationException;
use SimpleSoftwareIO\QrCode\BaconQrCodeGenerator;

class CIIEController extends Controller
{
    protected $prefixKey = "CIIE:";

    protected $activityId = 4;

    protected $activityAddress = "上海国家会展中心消费品展区6.2H 泰佩思琦 （Tapestry）展台";

    /**
     * 活动预约详情
     * @return JsonResponse
     * @throws ValidationException
     */
    public function Info(): JsonResponse
    {
        $rules = [
            'enterprise_id' => 'required',
            'activity_id' => 'required',
            'openid' => 'required'
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'activity_id.required' => '活动ID不能为空',
            'openid.required' => '用户身份错误'
        ];
        $params = $this->validate($this->request, $rules, $message);

        $enterprise_id = $params['enterprise_id'];
        $id = $params['activity_id'];
        $openid = $params['openid'];
        $result = [];

        $activityRedisKey=$this->prefixKey.$id;

        $userRecordKey=$this->prefixKey."record:".$id;
        try {
            $activityInfo=Redis::connection()->get($activityRedisKey);
            $info=json_decode($activityInfo,true);
//            $info = Activity::query()->where(['id' => $id, 'enterprise_id' => $enterprise_id])->first();
            $res=Redis::connection()->hget($userRecordKey,$openid);
            $result=json_decode($res,true);
//            $record = ActivityReserveRecord::query()->where([
//                'activity_id' => $id,
//                'enterprise_id' => $enterprise_id,
//                'openid' => $openid
//            ])->first();
            $isReserve = 0;
            if (!empty($result)) {
                $isReserve = 1;
            }
            $result['is_reserve'] = $isReserve;
            $result['activity_status'] = Activity::StatusRun;
            if ($info['start_time'] > time()) {
                $result['activity_status'] = Activity::StatusNotStarted;
            }
            if ($info['end_time'] < time()) {
                $result['activity_status'] = Activity::StatusEnd;
            }
            $reserve_time=$result['reserve_time']??0;
            $start_time = date("Y年m月d号 H:i", $result['reserve_time']??0);
            $end_time = date("H:i", $reserve_time + 3600);
            $result['reserve_time_str'] = "{$start_time}一{$end_time}";
            $result['reserve_address'] = $this->activityAddress;
            $result['mp_url'] = trim($info['activity_content']);
        } catch (Exception $e) {
            error($e->getMessage());
        }

        return successJson($result);
    }

    /**
     * 取消预约活动
     * @return JsonResponse
     * @throws ValidationException
     */
    public function Cancel(): JsonResponse
    {
        $params = $this->request->all();
        $rules = [
            'enterprise_id' => 'required',
            'activity_id' => 'required',
            'openid' => 'required',
            'reserve_id' => 'required',
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'activity_id.required' => '活动ID不能为空',
            'openid.required' => '用户openid不能为空',
            'reserve_id.required' => '无效的预约信息',
        ];
        $this->validate($this->request, $rules, $message);
        $record = ActivityReserveRecord::query()->where(['id' => $params['reserve_id']])->first();
        if (empty($record)) {
            error("预约信息不存在");
        }
        if ($record['openid'] != $params['openid']) {
            error("预约身份不正确");
        }
        if ($record['is_check']==1){
            error("已核销不能取消");
        }
        $userRecordKey=$this->prefixKey."record:".$params['activity_id'];
        ActivityReserveRecord::query()->where(['id' => $params['reserve_id']])->delete();
        $recordKey = $this->prefixKey . $record['reserve_time'];
        ActivityReserveTimeConfig::query()->where(['reserve_time' => $record['reserve_time']])->decrement("reserved_num");
        ActivityReserveTimeConfig::query()->where(['reserve_time' => $record['reserve_time']])->increment("reserve_num");
        Redis::connection()->incrby($recordKey, 1);
        Redis::connection()->hdel($userRecordKey,$params['openid']);
        return successJson();
    }

    /**
     * 修改预约活动
     * @return JsonResponse
     * @throws ValidationException
     */
    public function Edit(): JsonResponse
    {
        $params = $this->request->all();
        $rules = [
            'enterprise_id' => 'required',
            'activity_id' => 'required',
            'openid' => 'required',
            'reserve_time' => 'required',
            'mobile' => 'required',
            'user_name' => 'required|max:15',
            'reserve_id' => 'required'
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'activity_id.required' => '活动ID不能为空',
            'openid.required' => '用户openid不能为空',
            'reserve_time.required' => "预约时间点不能为空",
            'mobile.required' => "手机号不能为空",
            'user_name.required' => '姓名不能为空',
            'user_name.max' => "姓名长度不能超过15个字符",
            'reserve_id.required' => '无效的预约信息',
        ];
        $this->validate($this->request, $rules, $message);
        $record = ActivityReserveRecord::query()->where(['id' => $params['reserve_id']])->first();
        if (empty($record)) {
            error("预约信息不存在");
        }
        if ($record['openid'] != $params['openid']) {
            error("预约身份不正确");
        }
        if ($record['is_check']==1){
            error("已核销不能修改");
        }
        $save_data = [
            'mobile' => $params['mobile'],
            'user_name' => $params['user_name'],
            'gender' => $params['gender'],
        ];
        $userRecordKey=$this->prefixKey."record:".$params['activity_id'];
        ActivityReserveRecord::query()->where(['id' => $params['reserve_id']])->update($save_data);
        if ($record['reserve_time'] == $params['reserve_time']) {
            return successJson();
        }
        $redisKey = $this->prefixKey . $params['reserve_time'];
        $reserve_num=Redis::connection()->decrby($redisKey, 1);
        if ($reserve_num < 0) {// 如果预约扣减的可预约人数已经小于0，则设置可预约人数为0
            Redis::connection()->set($redisKey, 0);
            error("当前的预约已满");
        }
        ActivityReserveRecord::query()->where(['id' => $params['reserve_id']])->update(['reserve_time'=>$params['reserve_time']]);
        $recordKey = $this->prefixKey . $record['reserve_time'];
        $reserve_config = ActivityReserveTimeConfig::query()->where(['reserve_time' => $params['reserve_time']])->first();
        $reserve_record_config = ActivityReserveTimeConfig::query()->where(['reserve_time' => $record['reserve_time']])->first();
        if (!empty($reserve_config)) {
            ActivityReserveTimeConfig::query()->where(['id' => $reserve_config['id']])->increment("reserved_num");
            if ($reserve_config['reserve_num'] <= 1) {
                ActivityReserveTimeConfig::query()->where(['id' => $reserve_config['id']])->update(["reserve_num" => 0]);
            } else {
                ActivityReserveTimeConfig::query()->where(['id' => $reserve_config['id']])->decrement("reserve_num");// 修改的预约时间剩余-1
            }
            ActivityReserveTimeConfig::query()->where(['id' => $reserve_record_config['id']])->decrement("reserved_num"); //修改的历史已预约数-1
            Redis::connection()->incrby($recordKey, 1);
            $result = ActivityReserveRecord::query()->where(['id' => $params['reserve_id']])->first();
            Redis::connection()->hset($userRecordKey,$params['openid'],json_encode($result));
        }
        return successJson();
    }

    /**
     * 用户的活动进入上报
     * @return JsonResponse
     * @throws ValidationException
     */
    public function userView(): JsonResponse
    {
        $params = $this->request->all();
        $rules = [
            'enterprise_id' => 'required',
            'activity_id' => 'required',
            'openid' => 'required',
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'activity_id.required' => '活动ID不能为空',
            'openid.required' => '用户openid不能为空',
        ];
        $this->validate($this->request, $rules, $message);

        $userName = '';
        $unionId = $params['union_id'];
        if (!empty($unionId)) {
            $userInfos = WxUserCorpExternal::getExternalUserInfoByUnionIdNoCode($params['enterprise_id'], $unionId);
            $userInfo = current($userInfos);
            if (!empty($userInfo)) {
                $userName = $userInfo['nick'];
            }
        }
        $data = [
            'activity_id' => $params['activity_id'],
            'enterprise_id' => $params['enterprise_id'],
            'openid' => $params['openid'],
            'unionid' => $unionId,
            'user_name' => $userName,
            'source' => "ciie",
        ];
        try {
            $info = ActivityUserView::query()->where(["activity_id" => $params['activity_id'], 'openid' => $params['openid'], "source" => "ciie"])->first();
            if (empty($info)) {
                ActivityUserView::query()->create($data);
            }
        } catch (Exception $e) {
            error($e->getMessage());
        }
        return successJson();
    }

    /**
     * 创建预约活动
     * @return JsonResponse
     * @throws ValidationException
     */
    public function Create(): JsonResponse
    {
        $params = $this->request->all();
        $rules = [
            'enterprise_id' => 'required',
            'activity_id' => 'required',
            'openid' => 'required',
            'union_id' => 'required',
            'reserve_time' => 'required',
            'mobile' => 'required',
            'user_name' => 'required|max:15',
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'activity_id.required' => '活动ID不能为空',
            'openid.required' => '用户openid不能为空',
            'union_id.required' => 'union_id不能为空',
            'reserve_time.required' => "预约时间点不能为空",
            'mobile.required' => "手机号不能为空",
            'user_name.required' => '姓名不能为空',
            'user_name.max' => "姓名长度不能超过15个字符"
        ];
        $this->validate($this->request, $rules, $message);
        $userLock = "lock:" . $params['openid'];
        $redisKey = $this->prefixKey . $params['reserve_time'];
        try {
            $lock =  Redis::connection()->setnx($userLock, 1);
            if (!$lock) {
                throw new Exception("请稍后再试",100001);
            }
            Redis::connection()->expire($userLock,3);
            $nowTime = time();
            $reserve_time = $params['reserve_time'];
            $reserve_end_time = $reserve_time + 3600;
            if ($nowTime > $reserve_end_time) {
                throw new Exception("您预约的时间已结束",100001);
            }
            $record = ActivityReserveRecord::query()->where(['openid' => $params['openid']])->first();
            if (!empty($record)) {
                throw new Exception("您已经有预约了,请勿重复预约",100001);
            }
//        if ((empty($reserve_config) || $reserve_config['reserve_num'] == 0) && $nowTime < $reserve_time) {
//            error("您预约的时间还未开始");
//        }

            $reserve_num = Redis::connection()->decrby($redisKey, 1);
            if ($reserve_num < 0) {// 如果预约扣减的可预约人数已经小于0，则设置可预约人数为0
                Redis::connection()->set($redisKey, 0);
                throw new Exception("当前的预约已满",100001);
            }
            $unionId = $params['union_id'];

            $reserve_config = ActivityReserveTimeConfig::query()->where(['activity_id' => $params['activity_id'], 'reserve_time' => $params['reserve_time']])->first();
            // 生成预约活动码券码
            $activityCode = idGenerate();
            // 生成二维码
            $qrCodeGenerator = new BaconQrCodeGenerator();
            $qrcode = base64_encode($qrCodeGenerator->format("png")->margin(0)->size(220)->generate($activityCode));
            $avatar = '';
            if (!empty($unionId)) {
                $userInfos = WxUserCorpExternal::getExternalUserInfoByUnionIdNoCode($params['enterprise_id'], $unionId);
                $userInfo = current($userInfos);
                if (!empty($userInfo)) {
                    $avatar = $userInfo['avatar'];
                }
            }
            $data = [
                'enterprise_id' => $params['enterprise_id'],
                'activity_id' => $params['activity_id'],
                'unionid' => $params['union_id'],
                'openid' => $params['openid'],
                'mobile' => $params['mobile'],
                'avatar' => $avatar,
                'user_name' => $params['user_name'],
                'reserve_time' => $params['reserve_time'],
                'qr_code' => $qrcode,
                'activity_code' => $activityCode,
                'gender' => $params['gender'],
                'created_at' => time(),
                'updated_at' => time(),
            ];
            $id = ActivityReserveRecord::query()->insertGetId($data);
            if (!empty($reserve_config)) {
                ActivityReserveTimeConfig::query()->where(['id' => $reserve_config['id']])->increment("reserved_num");
                if ($reserve_config['reserve_num'] <= 1) {
                    ActivityReserveTimeConfig::query()->where(['id' => $reserve_config['id']])->update(["reserve_num" => 0]);
                } else {
                    ActivityReserveTimeConfig::query()->where(['id' => $reserve_config['id']])->decrement("reserve_num");
                }
            }
            $data['id']=$id;
            $data['is_check']=0;
            $userRecordKey=$this->prefixKey."record:".$params['activity_id'];
            Redis::connection()->hset($userRecordKey,$params['openid'],json_encode($data));
        } catch (Exception $e) {
            Redis::connection()->del($userLock);
            if($e->getCode()==100001){
                error($e->getMessage());
            }else{
                Redis::connection()->incrby($redisKey,1);
                error("生成预约失败" . $e->getMessage());
            }
        }
        Redis::connection()->del($userLock);
        return successJson(['reserve_id' => $id]);

    }

    /**
     * 导购的核销列表
     * @return JsonResponse
     * @throws ValidationException
     */
    public function CheckLists(): JsonResponse
    {
        $params = $this->request->all();
        $search = $params['search'] ?? "";
        $rules = [
            'enterprise_id' => 'required',
            'activity_id' => 'required',
            'mobile' => 'required',
            'page' => 'required|integer|min:1|max:999',
            'size' => 'required|integer|min:1|max:999'
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'activity_id.required' => '活动id不能为空',
            'mobile.required' => '导购身份错误',
            'page.required' => '页码不能为空',
            'page.integer' => '页码格式错误',
            'page.min' => '页码不能小于1',
            'page.max' => '页码不能大于999',
            'size.required' => 'size不能为空',
            'size.integer' => 'size格式错误',
            'size.min' => 'size不能小于1',
            'size.max' => 'size不能大于999',
        ];
        $this->validate($this->request, $rules, $message);
        $page = $params['page'];
        $size = $params['size'];
        $enterprise_id = $params['enterprise_id'];
        $memberInfo = Members::query()->where(['enterprise_id' => $enterprise_id, 'mobile' => $params['mobile']])->first();
        if (empty($memberInfo)) {
            error("导购不存在");
        }
        $memberid = $memberInfo['id'];
        try {
            $query = ActivityReserveRecord::query()->where(['enterprise_id' => $enterprise_id, 'activity_id' => $params['activity_id'], 'member_id' => $memberid]);
            if (!empty($search)) {
                $query->where(function ($query) use ($search) {
                    $query->where('mobile', "like", "%" . $search . "%")->orWhere("user_name", "like", "%" . $search . "%");
                });
            }
            $nums = $query->count();
            $offset = $page - 1;
            $lists = $query->offset($offset)->limit($size)->get();
            foreach ($lists as &$item) {
                $start_time = date("Y-m-d H:i", $item['reserve_time']);
                $end_time = date("H:i", $item['reserve_time'] + 3600);
                $item['reseve_time_str'] = "{$start_time}~{$end_time}";
            }
            unset($item);
            $result = ['count' => $nums, 'lists' => $lists];
        } catch (Exception $e) {
            error($e->getMessage());
        }
        return successJson($result);

    }

    /**
     * 核销预约码
     * @return JsonResponse
     * @throws ValidationException
     */
    public function WriteOffCoupon(): JsonResponse
    {
        $params = $this->request->all();
        $rules = [
            'enterprise_id' => 'required',
            'activity_id' => 'required',
            'mobile' => 'required',
            'code' => 'required',
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'activity_id.required' => '活动id不能为空',
            'code.required' => '预约核销码无效',
            'mobile.required' => '导购身份错误',
        ];
        $this->validate($this->request, $rules, $message);
        if (!in_array($params['mobile'], array_unique(array_merge(ActivityReserveRecord::WHITELISTS, ActivityReserveRecord::MANAGEWHITELISTS)))) {
            error("没有核销权限");
        }
        $enterprise_id = $params['enterprise_id'];
        $memberInfo = Members::query()->where(['enterprise_id' => $enterprise_id, 'mobile' => $params['mobile']])->first();
        if (empty($memberInfo)) {
            error("导购不存在");
        }
        $record = ActivityReserveRecord::query()->where(['enterprise_id' => $enterprise_id, 'activity_id' => $params['activity_id'], 'activity_code' => $params['code']])->first();
        if (empty($record)) {
            error("预约信息不存在");
        }
        $todayTime = strtotime(date("Y-m-d", time()));
        if ($record['reserve_time'] < $todayTime || $record['is_check'] == 2) {
            error("该预约码已失效");
        }
        if ($record['is_check'] == 1) {
            error("该预约码已被核销");
        }
        $checkData = [
            'member_id' => $memberInfo['id'],
            'member_name' => $memberInfo['name'],
            'check_time' => time(),
            'is_check' => 1,
            'version' => $record['version'] + 1
        ];
        try {
            ActivityReserveRecord::query()->where(['id' => $record['id'], 'version' => $record['version']])->update($checkData);
            $info=ActivityReserveRecord::query()->where(['id' => $record['id']])->first();
            $userRecordKey=$this->prefixKey."record:".$params['activity_id'];
            Redis::connection()->hset($userRecordKey,$info['openid'],json_encode($info));
        } catch (Exception $e) {
            error($e->getMessage());
        }
        return successJson();
    }

    /**
     * 导购白名单
     * @return JsonResponse
     * @throws ValidationException
     */
    public function WhiteList(): JsonResponse
    {
        $params = $this->request->all();
        $rules = [
            'enterprise_id' => 'required',
            'mobile' => 'required',
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'code.required' => '预约核销码无效',
        ];
        $this->validate($this->request, $rules, $message);
        $result = [
            'sa' => 0,
            'ma' => 0
        ];
        if (in_array($params['mobile'], ActivityReserveRecord::WHITELISTS)) {
            $result['sa'] = 1;
        }
        if (in_array($params['mobile'], ActivityReserveRecord::MANAGEWHITELISTS)) {
            $result['ma'] = 1;
        }
        return successJson($result);
    }

    /**
     * 时间段列表
     * @return JsonResponse
     */
    public function TimeLists(): JsonResponse
    {
        $params = $this->request->all();
        $reserve_time = $params['reserve_time'] ?? time();
        $start_time = strtotime(date("Y-m-d", $reserve_time));
        $end_time = $start_time + 86400;
        $type = $params['type'] ?? 2;
        try {
            $result = ActivityReserveTimeConfig::query()->whereBetween('reserve_time', [$start_time, $end_time])->orderBy('reserve_time')->groupBy("reserve_time")->get()->toArray();
            foreach ($result as &$item) {
                $redisKey = $this->prefixKey . $item['reserve_time'];
                $reserve_num = Redis::connection()->get($redisKey) ?? 0;
                $item['reserve_num'] = max($reserve_num, 0);
            }
            unset($item);
            if ($type == 1) {
                $res['lists'] = $result;
            } else {
                $res = ['lists' => $result];
                $date = [];
                $activity = Activity::query()->where(['id' => $this->activityId])->first();
                $days = ($activity['end_time'] - $activity['start_time']) / 86400;
                for ($i = 0; $i < $days; $i++) {
                    $date[] = date("Y-m-d", $activity['start_time'] + ($i * 86400));
                }
                $res['days'] = $date;
            }
            $res['today'] = date("Y-m-d", time());
        } catch (Exception $e) {
            error($e->getMessage());
        }

        return successJson($res);


    }

    /**
     * 预约时间段配置
     * @return JsonResponse
     */
    public function TimeConfig(): JsonResponse
    {
        $params = $this->request->all();
        $configs = $params['configs'];
        $enterprise_id = $params['enterprise_id'];
        $activity_id = $params['activity_id'];
        $nowTime = time();
        try {
            foreach ($configs as $config) {
                $config['activity_id'] = $activity_id;
                $config['enterprise_id'] = $enterprise_id;
                $record = ActivityReserveTimeConfig::query()->where(['reserve_time' => $config['reserve_time']])->first();
                if ($record) {
                    if ($config['reserve_time'] > $nowTime || ($config['reserve_time'] < $nowTime && $nowTime < $config['reserve_time'] + 3600)) {
                        ActivityReserveTimeConfig::query()->where(['reserve_time' => $config['reserve_time']])->update($config);
                        $redisKey = $this->prefixKey . $config['reserve_time'];
                        Redis::connection()->set($redisKey, $config['reserve_num']);
                    }
                } else {
                    ActivityReserveTimeConfig::query()->where(['reserve_time' => $config['reserve_time']])->create($config);
                    $redisKey = $this->prefixKey . $config['reserve_time'];
                    Redis::connection()->set($redisKey, $config['reserve_num']);
                }
            }
        } catch (Exception $e) {
            error($e->getMessage());
        }
        return successJson();
    }

    /**
     * 获取手机号
     * @return JsonResponse
     * @throws HttpException
     * @throws InvalidArgumentException
     * @throws InvalidConfigException
     * @throws RuntimeException
     * @throws \Psr\SimpleCache\InvalidArgumentException
     */
    public function GetPhoneNumber(): JsonResponse
    {
        $params = $this->request->all();
        try {
            $pluginLists = (new EnterpriseWorkPlugin())->lists(['enterprise_id' => $params['enterprise_id']]);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        if (empty($pluginLists)) {
            error(652012);
        }
        $appid = "";
        $secret = "";
        foreach ($pluginLists as $item) {
            if ($item['identifier'] == "crs_store_assistant_self") {
                $appid = $item['mini_app_id'];
                $secret = $item['mini_app_secret'];
            }
        }

        $token = MiniProgram::init(['app_id' => $appid, 'secret' => $secret])->access_token->getToken();
        $accessToken = $token['access_token'];
        $client = new Client();
        $url = "https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=" . $accessToken;
        $header = [
            "Content-Type" => "application/json",
            "Accept" => "application/json",
        ];
        $body = json_encode(['code' => $params['code']]);
        $options['headers'] = $header;
        $options['body'] = $body;
        $response = $client->post($url, $options)->getBody()->getContents();
        $rep = json_decode($response, true);
        if ($rep['errcode'] != 0) {
            error($rep['errcode'], $rep['errmsg']);
        }
        return successJson($rep['phone_info']);
    }

    /**
     * 查询全部所有的预约列表
     * @return JsonResponse
     * @throws ValidationException
     */
    public function Lists(): JsonResponse
    {
        $params = $this->request->all();
        $search = $params['search'] ?? "";
        $rules = [
            'enterprise_id' => 'required',
            'page' => 'required|integer|min:1|max:999',
            'size' => 'required|integer|min:1|max:999'
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'page.required' => '页码不能为空',
            'page.integer' => '页码格式错误',
            'page.min' => '页码不能小于1',
            'page.max' => '页码不能大于999',
            'size.required' => 'size不能为空',
            'size.integer' => 'size格式错误',
            'size.min' => 'size不能小于1',
            'size.max' => 'size不能大于999',
        ];
        $this->validate($this->request, $rules, $message);
        $page = $params['page'];
        $size = $params['size'];
        $enterprise_id = $params['enterprise_id'];
        try {
            $query = ActivityReserveRecord::query()->where(['enterprise_id' => $enterprise_id, "activity_id" => $this->activityId]);
            if (!empty($search)) {
                $query->where(function ($query) use ($search) {
                    $query->where('mobile', "like", "%" . $search . "%")->orWhere("user_name", "like", "%" . $search . "%");
                });
            }
            $nums = $query->count();
            $offset = ($page - 1) * $size;
            $lists = $query->orderBy('reserve_time', "desc")->offset($offset)->limit($size)->get();
            foreach ($lists as &$item) {
                $start_time = date("Y-m-d H:i", $item['reserve_time']);
                $end_time = date("H:i", $item['reserve_time'] + 3600);
                $item['reseve_time_str'] = "{$start_time}~{$end_time}";
            }
            unset($item);
            $result = ['count' => $nums, 'lists' => $lists];
        } catch (Exception $e) {
            error($e->getMessage());
        }

        return successJson($result);
    }
}