<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/12/25
 * Time: 6:17 PM
 */

namespace App\Http\Controllers\Callback\V1\Customer;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Http\Controllers\Callback\V1\BaseController;
use App\Jobs\ExternalContactUpdateCustomerStatisticJob;
use App\Jobs\Message\ExternalContactMessageAdapter;
use App\Library\Uc\External;
use App\Models\Channel;
use App\Models\ContactQrCode;
use App\Models\GuideUserRelation;
use App\Models\MemberGuide;
use App\Models\MemberUserRelation;
use App\Models\QwCallbackEnterprise;
use App\Models\StoreWorker;
use App\Models\User;
use App\Models\WelcomeSpeech;
use App\Services\CommonService;
use App\Services\MemberService;
use App\Services\StoreWorkerService;
use App\Services\Message\CustomerRelationBindMessageService;
use App\Services\Message\MessageService;
use CrsSdk\Boss\Enterprise;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;

class AddController extends BaseController
{
    /**
     * UC外部联系人添加成功回调
     *
     * @param Request $request
     * @return string
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params = json_decode(file_get_contents('php://input'), true);

        Log::info('UC 外部联系人回调进来了 enterprise_id: ' . $request->get('enterprise_id') . 'hash:' . $request->header('enterprise_hash'));
        Log::info('UC 外部联系人回调进来了 params: ' . var_export($params, true));

        $params['enterprise_id'] = $request->get('enterprise_id');

        $type = $params['addType'] ?? $params['delType'];
        switch ($type) {
            // 正常添加成员
            case 'add_external_contact' :
                $params['em_state'] = Constant::ACTIVE_ADDITION;
                if (!empty($params['emState'])) {
                    $params['em_state'] = $params['emState'];
                }
//                $externalDetail = $this->externalDetail($params);
//                if (
//                    isset($externalDetail['followMembers']) &&
//                    !empty($externalDetail['followMembers'])
//                ) {
//                    $members = array_column($externalDetail['followMembers'], null, 'memUid');
//                    if (
//                        isset($members[$params['memUid']]) &&
//                        !empty($members[$params['memUid']]) &&
//                        !empty($members[$params['memUid']]['emState'])
//                    ) {
//                        $params['em_state'] = $members[$params['memUid']]['emState'];
//                    }
//                }

                //这里需要加锁，防止重复提交
                $lockKey = 'add_contact:' . $params['enterprise_id'] . '-' . $params['memUid'] . '-' . $params['ecId'];
                // 查询锁
                $getCache = Redis::get($lockKey);
                // 检查是否已经超出时间 超出时间释放锁（防止死锁）
                if (!empty($getCache) && $getCache < time()) {
                    Redis::del([$lockKey]);
                }
                // 开始加锁
                $setCache = Redis::setnx($lockKey, time() + 60);
                // 如果已经发起直接跳出
                if (!$setCache) {
                    return 'success';
                }

                // 发送欢迎语
                try {
                    $res = $this->addExternalContact($params);
                    Redis::del([$lockKey]);
                    if ($res) {
                        $this->welcomeSend($params);
                    }
                } catch (\Exception $e) {
                    Redis::del([$lockKey]);
                    Log::info('发送欢迎语失败: ' . 'code---' . $e->getCode() . 'message---' . $e->getMessage());
                }
                break;
            // 免验证添加成员
            case 'add_half_external_contact' :
                $params['em_state'] = Constant::ACTIVE_ADDITION;
                if (isset($params['emState']) && !empty($params['emState'])) {
                    $params['em_state'] = $params['emState'];
                }
                // 发送欢迎语
                try {
                    $this->welcomeSend($params);
                } catch (\Exception $e) {
                    Log::info('发送欢迎语失败: ' . 'code---' . $e->getCode() . 'message---' . $e->getMessage());
                }
                break;
            // 删除外部联系人
            case 'del_external_contact' :
                $params['delete_type'] = MemberUserRelation::DEL_EXTERNAL_CONTACT;
                $this->delete($params);
                break;
            // 删除跟进成员
            case 'del_follow_user' :
                $params['delete_type'] = MemberUserRelation::DEL_FOLLOW_USER;
                $this->delete($params);
                break;
            default :
                break;
        }

        return 'success';
    }

    /**
     * 添加成员
     *
     * @param $params
     * @return bool|string
     * @throws \Exception
     */
    protected function addExternalContact($params)
    {
        $enterpriseId = $params['enterprise_id'];
        $memUid = $params['memUid'];
        $ecId = $params['ecId'] ?? '';
        $ecUserId = $params['ecUserid'] ?? '';
        $ecName = $params['ecName'] ?? '';
        $ecAvatar = $params['ecAvatar'] ?? '';
        $ecGender = $params['ecGender'] ?? 3;
        $userNameIndex = $params['userNameIndex'] ?? '';
        $unionId = empty($params['unionId']) ? '' : $params['unionId'];
        $contactQrCodeId = intval($params['em_state']);

        if (empty($enterpriseId)) {
            throw new Error(100000, '企业ID不能为空');
        }

        if (empty($memUid)) {
            throw new Error(100000, 'memUid不能为空');
        }

        if (empty($ecUserId)) {
            throw new Error(100000, 'ecUserid不能为空');
        }

        $memUser = MemberUserRelation::query()->where([
            'enterprise_id' => $enterpriseId,
            'mem_uid'       => $memUid,
            'ec_id'         => $ecId
        ])->first();

        // 存在了就不需要再加了
        if (!empty($memUser)) {
            if ($memUser->delete_type != 0) {
                $memUser->delete_type = 0;
                $memUser->source_type = $contactQrCodeId ?? 1;
                $memUser->created_at = time();
                $memUser->save();

                ExternalContactUpdateCustomerStatisticJob::sendCreateMessage($memUser->toArray());
                Log::info('单方面删除的客户恢复为正常客户,delete_type:' . $memUser->delete_type);
                return true;
            } else {
                Log::info('该客户已经存在了');
                return false;
            }
        }

        // 查询当前外部联系人绑定的客户
        $bindingMemUser = MemberUserRelation::query()->where([
            'enterprise_id'       => $enterpriseId,
            'is_binding_customer' => MemberUserRelation::BINDING_CUSTOMER_YES,
            'ec_id'               => $ecId
        ])->first();

        // 查询user
        $user = User::query()->where('union_id', $unionId)->first();

        $insertData = [
            'enterprise_id'       => $enterpriseId,
            'mem_uid'             => $memUid,
            'ec_id'               => $ecId,
            'ec_userid'           => $ecUserId,
            'ec_name'             => $ecName,
            'ec_avatar'           => $ecAvatar,
            'ec_gender'           => $ecGender,
            'user_name_index'     =>
                (preg_match('/^[a-zA-Z\x{4e00}-\x{9fa5}]+/u', $userNameIndex)) ? $userNameIndex : '井',
            'union_id'            => $unionId,
            'relation_at'         => time(),
            'is_binding_customer' => !empty($bindingMemUser) ? MemberUserRelation::BINDING_CUSTOMER_NO : MemberUserRelation::BINDING_CUSTOMER_YES,
            'source_type'         => $contactQrCodeId ?? 1,
            'user_id'             => $user ? $user->id : 0,
            'created_at'          => time(),
            'updated_at'          => time()
        ];

        $result = MemberUserRelation::query()->create($insertData);
        ExternalContactUpdateCustomerStatisticJob::sendCreateMessage($insertData);
        return $result;
    }

    /**
     * 删除成员
     * @param $params
     * @return bool|string
     */
    protected function delete($params)
    {
        $memUid = $params['memUid'];
        $ecId = $params['ecId'];
        $deleteType = $params['delete_type'];

        $memUser = MemberUserRelation::query()->where([
            'mem_uid' => $memUid,
            'ec_id'   => $ecId
        ])->first();

        if (empty($memUser)) {
            Log::info('该客户不存在');
            return 'success';
        }

        // 如果存在单方面删除,另一方删除就直接删除数据
        if ($memUser->delete_type != 0) {
            $memUser->delete();
        } else {
            $memUser->delete_type = $deleteType;
            $memUser->cancel_relation_at = time();
            $memUser->save();
        }

        ExternalContactUpdateCustomerStatisticJob::sendDeleteMessage($memUser->toArray());

        return true;
    }

    /**
     * 获取外部联系人详情
     *
     * @param $params
     * @return mixed
     * @throws \Exception
     */
    public function externalDetail($params)
    {
        // 获取企业domain
        $enterpriseData = (new Enterprise())->detail(
            ['enterprise_id' => $params['enterprise_id']]
        );
        $qyDomain = $enterpriseData['data']['domain'] ?? '';
        return (new External())->detail($qyDomain, ['ecId' => $params['ecId']]);
    }

    /**
     * 发送欢迎语
     *
     * @param $params
     * @throws \Exception
     */
    public function welcomeSend($params)
    {
        return true; //小程序顾客列表移除，发消息功能暂时移除 2021.12.2
        /**$emState = intval($params['em_state']);
        $welcomeCode = $params['welcomeCode'];

        // 获取企业domain
        $enterpriseData = (new Enterprise())->detail(
            ['enterprise_id' => $params['enterprise_id']]
        );

        $qyDomain       = $enterpriseData['data']['domain'] ?? '';
        $params['qy_domain'] = $qyDomain;

        // 2020.12.08，关闭欢迎语，由SCRM给用户发送
//        if ($welcomeCode) {
//            $message = [
//                'welcomeCode' => $welcomeCode,
//                'text'        => [
//                    'content' => WelcomeSpeech::DEFAULT_CONTENT['content']
//                ]
//            ];
//
//            // 默认欢迎语数据
//            $default = WelcomeSpeech::query()
//                ->where('enterprise_id', $params['enterprise_id'])
//                ->where('is_default', WelcomeSpeech::IS_DEFAULT_YES)
//                ->first();
//
//            if ($emState > 1) {
//                $contactQrCode = ContactQrCode::query()->where('id', $emState)->first();
//                if ($contactQrCode) {
//                    $welcomeMsg = WelcomeSpeech::query()
//                        ->where('enterprise_id', $params['enterprise_id'])
//                        ->where('channel_id', $contactQrCode->channel_id)
//                        ->first();
//
//                    if ($welcomeMsg) {
//                        $default = $welcomeMsg;
//                    }
//                }
//
//            }
//
//            if ($default) {
//                $message = $this->buildMessage($message, $default, $params);
//            }
//
//            Log::info('welcome-params: ' . var_export($message, true));
//
//            $sendResult = (new External())->welcomeSend($qyDomain, $message);
//            // 判断uc返回数据，针对"不合法的媒体文件id"情况，刷新媒体文件id重新发送
//            if (!is_array($sendResult)) {
//                $resultArr = json_decode($sendResult, true);
//                if (!empty($resultArr['code']) && $resultArr['code'] == '40007' && $default) {
//                    $default = $this->refreshWelcomeMedia($default);
//                    $message = $this->buildMessage($message, $default, $params);
//                    Log::info('welcome-params: ' . var_export($message, true));
//                    (new External())->welcomeSend($qyDomain, $message);
//                }
//            }
//        }

        try {
            if (!empty($params['em_state'])) {
                $memUser = MemberUserRelation::query()->where([
                    'enterprise_id' => $params['enterprise_id'],
                    'mem_uid'       => $params['memUid'],
                    'ec_userid'     => $params['ecUserid'],
                    'delete_type'   => 0
                ])->first();
                if ($memUser) {
                    $workerInfo = StoreWorkerService::detail(
                        [
                            'contact_qr_code_id' => $params['em_state'],
                            'enterprise_id'      => $params['enterprise_id']
                        ]);
                    if (!empty($workerInfo)) {
                        try {
                            $user = (new \CrsSdk\MicroUser\User())->userDetail([
                                'enterprise_id'   => $params['enterprise_id'],
                                'union_id'        => $params['unionId'],
                            ]);
                        } catch (\Exception $exception) {
                            Log::info('商城用户信息不存在');
                        }
                        $content = [
                            [
                                'key' => '顾客昵称',
                                'value' => $params['ecName'] ?? '',
                            ],
                            [
                                'key' => '温馨提醒',
                                'value' => '你可以在顾客管理列表中和这个顾客打招呼哦~',
                            ]
                        ];
                        $title   = '新增顾客提醒';
                        $memUids = [$params['memUid']];
                        $httpBuildParams = [
                            'store_id' => $workerInfo['store_id'] ?? 0,
                            'item' => json_encode([
                                'relation_id' => $memUser->id,
                                'mem_uid' => $params['memUid'],
                                'user_id' => $user['id'] ?? 0,
                                'union_id' => $params['unionId']
                            ], true)
                        ];
                        $page    = 'pages/customer-info/customer-info?' . http_build_query($httpBuildParams);

                        (new MessageService())->sendMiniMessage($params['enterprise_id'], $content, $page, $title, $memUids, $params['qy_domain']);
                    }
                }

            }
        } catch (\Exception $exception) {
            Log::info('新增顾客提醒异常' . $exception);
        }*/
    }

    /**
     * buildMessage
     * @desc 构建消息体
     * @param $message array
     * @param $welcomeSpeech WelcomeSpeech
     * @param $params array
     * @return array
     * @date 2020-08-11
     */
    public function buildMessage($message, $welcomeSpeech, $params)
    {
        $message['text'] = [
            'content' => $welcomeSpeech->content
        ];

        switch ($welcomeSpeech->type) {
            case 2 : // 图片
                $messageType = 'image';
                break;
            case 3 : // 网页
                $messageType = 'link';
                break;
            case 4 : // 小程序
                $messageType = 'miniprogram';
                break;
            default :
                $messageType = '';
        }

        if ($messageType) {
            $data = json_decode($welcomeSpeech->object, true);
            $httpBuildParams = [
                'mem_uid' => $params['memUid']
            ];
            if ($messageType == 'miniprogram') {
                // 引流红包
                if (strpos($data['page'], 'other/red-envelope/red-envelop') !== false) {
                    $data['page'] = $data['page'] . '&m=' . $params['memUid'];
                }

                // 携带导购workerId
                if (!empty($params['em_state'])) {
                    switchAppHostByDomain($params['qy_domain']);
                    $workerInfo = StoreWorkerService::detail(
                        [
                            'contact_qr_code_id' => $params['em_state'],
                            'enterprise_id'      => $params['enterprise_id'],
                            'mem_uid'            => $params['memUid'],
                            'worker_status'      => StoreWorker::IN_SERVICE
                        ]);
                    if (!empty($workerInfo)) {
                        $httpBuildParams['worker_id'] = $workerInfo['id'];
                        $httpBuildParams['store_id'] = $workerInfo['store_id'];
                    }
                }

                if (!empty($httpBuildParams)) {
                    $concatStr = (strpos($data['page'], '?') === false) ? '?' : '&';
                    $data['page'] = $data['page'] . $concatStr . http_build_query($httpBuildParams);
                }
            }

            $message[$messageType] = $data;
        }
        return $message;
    }

    /**
     * refreshWelcomeMedia
     * @desc 更新微信临时媒体id
     * @param $welcomeSpeech Model
     * @return Model
     * @date 2020-08-11
     *
     * @throws \Exception
     */
    public function refreshWelcomeMedia($welcomeSpeech)
    {
        $appIdentifier = 'contact';
        if (empty($welcomeSpeech->original_data)) {
            Log::info('refreshWelcomeMedia: original_data is empty', $welcomeSpeech->toArray());
            return $welcomeSpeech;
        }
        $originalData = json_decode($welcomeSpeech->original_data, true);
        if (empty($originalData)) {
            Log::info('refreshWelcomeMedia: original_data decode failed', $welcomeSpeech->toArray());
            return $welcomeSpeech;
        }
        $pic = $originalData['pic'];
        $object = '';
        $enterpriseId = $welcomeSpeech->enterprise_id;
        $enterprise = (new Enterprise)->detail(['enterprise_id' => $enterpriseId]);
        $qyDomain = $enterprise['data']['domain'] ?? '';
        if ($qyDomain) {
            if (WelcomeSpeech::IMAGE == $welcomeSpeech->type) {
                $object = json_decode($welcomeSpeech->object, true);

                $object['mediaId'] = MemberService::getWorkWeChartMediaId($qyDomain, $appIdentifier, $pic);

                $object = json_encode($object, JSON_UNESCAPED_UNICODE);
            }
            if (WelcomeSpeech::MINI_PROGRAM == $welcomeSpeech->type) {
                $object = json_decode($welcomeSpeech->object, true);

                $object['picMediaId'] = MemberService::getWorkWeChartMediaId($qyDomain, $appIdentifier, $pic);

                $object = json_encode($object, JSON_UNESCAPED_UNICODE);
            }
            if ($object) {
                $welcomeSpeech->object = $object;
                $welcomeSpeech->save();
            }
        }
        return $welcomeSpeech;
    }
}

