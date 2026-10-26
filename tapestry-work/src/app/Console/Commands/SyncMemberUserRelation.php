<?php
/**
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/2/28
 * Time: 16:32
 */

namespace App\Console\Commands;

use App\Jobs\ExternalContactUpdateCustomerStatisticJob;
use App\Library\Uc\External;
use App\Library\Uc\Member;
use App\Models\MemberUserRelation;
use App\Models\Scrm\Members;
use App\Models\Scrm\WorkExternalUserRelation;
use App\Models\ScrmEnterprise\WxUserCorpExternal;
use App\Models\StoreWorker;
use CrsSdk\Boss\Enterprise;
use Illuminate\Console\Command;

class SyncMemberUserRelation extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto:sync-member-user-relation {--enterprise_id=}';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '同步scrm中的外部联系人到智慧导购';

    public function handle()
    {
        $enterpriseId = intval($this->option('enterprise_id'));
        if (!$enterpriseId) return;
        //如果数据不一致则同步数据
        try {
            $enterpriseInfo = (new Enterprise())->detail(['enterprise_id' => $enterpriseId]);
        } catch (\Throwable $e) {
            echo '企业信息获取失败,enterprise_id:' . $enterpriseId . PHP_EOL;
            return;
        }
        // 查询所有权限
        $workerList = StoreWorker::query()
            ->where('worker_status', StoreWorker::IN_SERVICE)
            ->where('enterprise_id', $enterpriseId)
            ->get(['mem_uid', 'mem_userid', 'enterprise_id'])
            ->toArray();
        $time       = time();
        foreach ($workerList as $val) {
            //获取该导购外部联系人openid（memberUserRelation表）
            $memUid             = $val['mem_uid'];
            $memberUserRelOpenIdArr = MemberUserRelation::query()
                ->where('enterprise_id', $enterpriseId)
                ->where('mem_uid', $memUid)
                ->pluck('ec_userid')->toArray();
            //获取该企业所有导购外部联系人openid（wx_user_corp_external表）
            $userId = $val['mem_userid'];
            //查询导购在scrm库中相关数据
            $memberId = Members::query()
                ->where('enterprise_id', $enterpriseId)
                ->where('userid', $userId)
                ->value('id');
            //查询导购在scrm库中外部联系人关系数据
            $workExternalUserRelUserIdArr = WorkExternalUserRelation::query()
                ->where('member_id', $memberId)
                ->pluck('work_external_user_id')->toArray();
            //查询导购在scrm苦衷外部联系人openid数据
            $wxUserCorpExternalOpenIdArr = WxUserCorpExternal::query()
                ->setConn($enterpriseId)
                ->whereIn('id', $workExternalUserRelUserIdArr)
                ->where('unionid', '<>', '')
                ->pluck('openid')->toArray();
            //将智慧导购外部联系人数据和scrm外部联系人数据对比
            $newOpenIdArr = array_values(array_diff($wxUserCorpExternalOpenIdArr, $memberUserRelOpenIdArr));
            if (empty($newOpenIdArr)) {
                echo '数据是一样的，没毛病,enterprise_id:' . $enterpriseId . PHP_EOL;
                continue;
            }

            $addNum = 0;
            foreach ($this->getWxUserCorpExternalDetail($enterpriseInfo['data']['domain'], $newOpenIdArr) as $info) {
                $insertData = [];
                foreach ($info['followMembers'] ?? [] as $v) {
                    if ($v['memUid'] == $memUid) {
                        $insertData[] = [
                            'enterprise_id'       => $enterpriseId,
                            'mem_uid'             => $memUid,
                            'ec_id'               => $info['ecId'],
                            'ec_userid'           => $info['ecUserid'],
                            'ec_name'             => $info['ecName'],
                            'ec_avatar'           => $info['ecAvatar'],
                            'ec_gender'           => $info['ecGender'],
                            'user_name_index'     => (preg_match('/^[a-zA-Z\x{4e00}-\x{9fa5}]+/u', $info['ecIndex'])) ? $info['ecIndex'] : '井',
                            'union_id'            => $info['ecUnionid'],
                            'relation_at'         => intval($info['ecCreated'] / 1000),
                            'is_binding_customer' => MemberUserRelation::BINDING_CUSTOMER_YES,
                            'source_type'         => 1,
                            'user_id'             => 0,
                            'created_at'          => $time,
                            'updated_at'          => $time
                        ];
                        $addNum++;
                    }
                }

                if (!empty($insertData)) {
                    MemberUserRelation::query()->insert($insertData);
                    foreach ($insertData as $insArr) {
                        ExternalContactUpdateCustomerStatisticJob::sendCreateMessage($insArr);
                    }
                }
            }
            echo $enterpriseId . '插入数据:' . $addNum . '条';

        }
        die('处理完成');
    }

    private function getWxUserCorpExternalDetail($qyDomain, $newOpenIdArr)
    {
        foreach ($newOpenIdArr as $key => $openId) {
            $condition = [
                'ecUserid' => $openId
            ];
            $data      = (new External())->listByUserId($qyDomain, $condition)['list'];

            yield $data[0] ?? [];
        }

    }

}
