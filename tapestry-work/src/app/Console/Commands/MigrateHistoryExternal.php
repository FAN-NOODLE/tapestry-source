<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/3/30 0030
 * Time: 10:14
 */

namespace App\Console\Commands;


use App\Jobs\ExternalContactUpdateCustomerStatisticJob;
use App\Library\DataLake\UserRequestService;
use App\Library\Uc\External;
use App\Models\MemberUserRelation;
use App\Models\StoreWorker;
use App\Models\User;
use CrsSdk\Boss\Enterprise;
use Illuminate\Console\Command;

class MigrateHistoryExternal extends Command
{

    /**
     * 单次查询数量
     */
    const PER_PAGE = 2;


    /**
     * 命令
     * @var string
     */
    protected $signature = 'migrate_history_external {--enterprise_id=}';

    /**
     * 描述信息
     * @var string
     */
    protected $description = '同步历史外部联系人数据';

    /**
     * @var UserRequestService
     */
    protected $ucExternal;


    /**
     * SyncDataLakeInfo constructor.
     */
    public function __construct()
    {
        parent::__construct();
        $this->ucExternal = new External();

    }

    /**
     * @param int $enterpriseId
     * @param array $workerIds
     * @return bool
     * @throws \Exception User: zhaoyang@vchangyi.com
     * Date: 2021/4/23 14:16
     */
    public function handle($enterpriseId = 0, $workerIds = [])
    {
        try{
            $enterpriseId = intval($this->option('enterprise_id'));
        } catch (\Throwable $exception ){
        }
        if (empty($enterpriseId)) {
            echo '企业id必传';
            return true;
        }
        $enterpriseInfo = (new Enterprise())->detail(['enterprise_id' => $enterpriseId]);
        if (empty($enterpriseInfo)) {
            echo '找不到企业';
            return true;
        }
        $qyDomain  = $enterpriseInfo['data']['domain'] ?? '';
        $num       = 0;
        $insertIds = [];
        foreach ($this->workerList($enterpriseId, $workerIds) as $memUids) {
            foreach ($memUids as $memUid) {
                $time      = time();
                $oldUserId = MemberUserRelation::query()
                    ->where('mem_uid', $memUid)
                    ->where('enterprise_id', $enterpriseId)
                    ->pluck('ec_userid')->toArray();
                $data      = $this->ucExternal->listForMember($qyDomain, [
                    'memUid'    => $memUid,
                    'pageSize' => 99999
                ]);
                if (!empty($data['list'])) {
                    $trueUserId = array_column($data['list'], 'ecUserid');
                    $newUserId  = array_diff($trueUserId, $oldUserId);
                    if (!empty($newUserId)) {
                        foreach ($data['list'] as $val) {
                            if (in_array($val['ecUserid'], $newUserId)) {
                                if (!empty($val['ecUnionid'])) {
                                    $user = User::query()->where('union_id', $val['ecUnionid'])->first();
                                }
                                $insertData = [
                                    'enterprise_id'       => $enterpriseId,
                                    'mem_uid'             => $memUid,
                                    'ec_id'               => $val['ecId'],
                                    'ec_userid'           => $val['ecUserid'],
                                    'ec_name'             => $val['ecName'],
                                    'ec_avatar'           => $val['ecAvatar'],
                                    'ec_gender'           => $val['ecGender'],
                                    'user_name_index'     => $val['ecIndex'],
                                    'union_id'            => $val['ecUnionid'] ?? '',
                                    'relation_at'         => $time,
                                    'is_binding_customer' => MemberUserRelation::BINDING_CUSTOMER_YES,
                                    'source_type'         => MemberUserRelation::SOURCE_TAKE_SHOPPERS_ADD,
                                    'user_id'             => !empty($user) ? $user->id : 0,
                                    'created_at'          => $time,
                                    'updated_at'          => $time
                                ];

                                $result      = MemberUserRelation::query()->create($insertData);
                                $insertIds[] = $result->id;
                                ExternalContactUpdateCustomerStatisticJob::sendCreateMessage($insertData);
                                $num++;
                            }
                        }
                    }
                }
            }
        }

        echo "本次新增" . $num . "条数据" . join(',', $insertIds);
        return true;
    }


    /**
     * 获取企业导购
     *
     * @param $enterpriseId
     * @param array $workerIds
     * @return \Generator
     * User: zhaoyang@vchangyi.com
     * Date: 2021/4/19 16:53
     */
    private function workerList($enterpriseId, $workerIds = []): \Generator
    {
        $page = 1;
        do {
            $res = StoreWorker::query()
                ->where('enterprise_id', $enterpriseId)
                ->when(!empty($workerIds), function ($query) use ($workerIds) {
                    $query->whereIn('id', $workerIds);
                })
                ->where('worker_status', StoreWorker::IN_SERVICE)
                ->forPage($page, static::PER_PAGE)
                ->pluck('mem_uid')
                ->toArray();
            yield $res;
            $page++;
        } while ($res !== []);
    }
}