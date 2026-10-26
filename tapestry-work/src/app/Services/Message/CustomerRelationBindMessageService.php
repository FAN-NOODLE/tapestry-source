<?php


namespace App\Services\Message;

use App\Models\MemberUserRelation;
use App\Services\CommonService;
use App\Services\StoreWorkerService;
use CrsSdk\Boss\Enterprise;
use CrsSdk\MicroUser\User;
use CrsSdk\Shop\Customer;
use Illuminate\Support\Facades\Log;
use \Exception;

/**
 * 顾客绑定导购 消息
 *
 * Class CustomerBindGuideMessageService
 * @package App\Services\Message
 */
class CustomerRelationBindMessageService
{

    const AUTO_EXEC_DEFAULT_INTERVAL = 300;

    const AUTO_EXEC_PAGE_SIZE = 100;

    /**
     * 定时执行
     * @param int $interval
     */
    public function autoExec(int $interval)
    {
        Log::info("CustomerRelationBindMessageService.autoExec start", ['time' => time()]);

        $startTime = time() - ($interval ?? self::AUTO_EXEC_DEFAULT_INTERVAL);

        $startRowId = MemberUserRelation::query()
            ->where('created_at', '>=', $startTime)
            ->orderBy('id')
            ->value('id');

        if (empty($startRowId)) {
            Log::info("customer relation bind message auto exec, no data");
            return;
        }

        $endRowId = MemberUserRelation::query()->orderBy('created_at', 'desc')->value('id');
        Log::info("customer relation bind message auto exec", [
            'startRowId' => $startRowId,
            'endRowId' => $endRowId
        ]);


        $relationTotal = 0;
        $sendTotal = 0;
        $page = 1;
        while (true) {
            $offset = ($page - 1) * self::AUTO_EXEC_PAGE_SIZE;
            $relationList = MemberUserRelation::query()
                ->select(['id', 'enterprise_id', 'union_id', 'mem_uid', 'created_at'])
                ->where('id', '>=', $startRowId)
                ->where('id', '<=', $endRowId)
                ->where('union_id', '!=', '')
                ->where('mem_uid', '!=', '')
                ->where('delete_type', MemberUserRelation::DEL_DEFAULT)
                ->limit(self::AUTO_EXEC_PAGE_SIZE)->offset($offset)->get()->toArray();
            if (empty($relationList)) {
                Log::info("customer relation bind message page break", ['page' => $page]);
                break;
            }
            $count = count($relationList);
            $relationTotal += $count;
            Log::info("customer relation bind message page exec", ['page' => $page, 'count' => $count]);
            foreach ($relationList as $relation) {
                if ($this->send($relation)) {
                    $sendTotal++;
                }
            }
            $page++;
        }

        Log::info("CustomerRelationBindMessageService.autoExec end", [
            'relationTotal' => $relationTotal,
            'sendTotal' => $sendTotal
        ]);
    }

    /**
     * 添加外部联系人、分享链路绑定、用户商城授权  三个条件同时满足
     *
     * @param $relation
     * @return bool
     */
    protected function send($relation): bool
    {
        Log::info("发送顾客绑定成功消息请求参数", ['relation' => $relation]);

        $flagKey = redisKey('REDIS_SHOP_AUTH_PREFIX', [
            'enterpriseId' => $relation['enterprise_id'],
            'unionId' => $relation['union_id'],
            'memUid' => $relation['mem_uid']
        ]);
        if (false != getOrSetData($flagKey)) {
            Log::info("发送顾客绑定成功消息已经发送");
            return false;
        }

        static $enterpriseList = [];
        // 查询企业信息
        $enterprise = $enterpriseList[$relation['enterprise_id']] ?? null;
        if ($enterprise == null) {
            $params = ['enterprise_id' => $relation['enterprise_id']];
            try {
                $enterprise = $enterpriseData = (new Enterprise())->detail($params)['data'];
                $enterpriseList[$relation['enterprise_id']] = $enterprise;
            } catch (Exception $e) {
                Log::error('企业信息查询失败', [
                    'params' => $params,
                    'code' => $e->getCode(),
                    'message' => $e->getMessage(),
                    'exception' => $e->getTraceAsString()
                ]);
            }
        }

        // 查询用户信息
        $user = null;
        $params = [
            'enterprise_id' => $relation['enterprise_id'],
            'union_id' => $relation['union_id']
        ];
        try {
            $_SERVER['enterprise_hash'] = $enterprise['hash'];
            $user = (new User())->userDetail($params);
        } catch (Exception $e) {
            Log::error('商城用户信息查询失败', [
                'params' => $params,
                'code' => $e->getCode(),
                'message' => $e->getMessage(),
                'exception' => $e->getTraceAsString()
            ]);
        }

        // 查询用户分享链路
        $shareLink = null;
        if ($user != null) {
            $params = [
                'enterprise_id' => $relation['enterprise_id'],
                'user_id' => $user['id']
            ];
            try {
                $shareLink = (new Customer())->customerBindingWorkerDetail($params);
            } catch (Exception $e) {
                Log::error('绑定关系查询失败', [
                    'params' => $params,
                    'code' => $e->getCode(),
                    'message' => $e->getMessage(),
                    'exception' => $e->getTraceAsString()
                ]);
            }
        }

        if (empty($enterprise) || empty($user) || empty($shareLink) || $relation['mem_uid'] != $shareLink['mem_uid']) {
            Log::info("顾客绑定成功消息条件不满足", [
                'enterprise' => $enterprise,
                'user' => $user,
                'shareLink' => $shareLink
            ]);
            return false;
        }

        try {
            $this->sendMiniMessage($relation['enterprise_id'], $enterprise['domain'], $user['id'], $user['nick_name'],
                $user['mobile'], $user['union_id'], $relation['mem_uid'], $relation['id']);

            getOrSetData($flagKey, 1);
            Log::info("发送顾客绑定成功消息发送成功");
            return true;
        } catch (Exception $e) {
            Log::error("发送顾客绑定成功消息发送失败", [
                'relation' => $relation,
                'enterprise' => $enterprise,
                'user' => $user,
                'code' => $e->getCode(),
                'message' => $e->getMessage(),
                'exception' => $e->getTraceAsString()
            ]);
            return false;
        }
    }

    /**
     * @param $enterpriseId
     * @param $qyDomain
     * @param $userId
     * @param $userNickName
     * @param $userMobile
     * @param $userUnionId
     * @param $workerMemUid
     * @param $relationId
     * @throws \Exception
     */
    protected function sendMiniMessage($enterpriseId,
                                       $qyDomain,
                                       $userId,
                                       $userNickName,
                                       $userMobile,
                                       $userUnionId,
                                       $workerMemUid,
                                       $relationId): void
    {
        Log::info("CustomerRelationBindMessageService.sendMiniMessage", func_get_args());

        $title = '顾客关系绑定成功';
        $content = [
            [
                'key' => '顾客昵称',
                'value' => !empty($userNickName) ? $userNickName : '',
            ],
            [
                'key' => '顾客手机号',
                'value' => !empty($userMobile) ? $userMobile : '顾客暂未授权',
            ],
            [
                'key' => '温馨提醒',
                'value' => '可以点击下方“进入小程序”查看顾客详情哦~',
            ]
        ];
        $httpBuildParams = [
            'store_id' => 0,
            'item' => json_encode([
                'relation_id' => $relationId,
                'mem_uid' => $workerMemUid,
                'user_id' => $userId,
                'union_id' => $userUnionId
            ], true)
        ];
        $page = 'pagesA/customer-info/customer-info?' . http_build_query($httpBuildParams);

        $memUids = [$workerMemUid];
        (new MessageService())->sendMiniMessage($enterpriseId, $content, $page, $title, $memUids, $qyDomain);
    }

}