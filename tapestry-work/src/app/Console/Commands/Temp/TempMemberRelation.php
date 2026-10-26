<?php


namespace App\Console\Commands\Temp;


use App\Common\Constant;
use App\Exceptions\Error;
use App\Library\Uc\External;
use App\Models\MemberUserRelation;
use Carbon\Carbon;
use CrsSdk\Boss\Enterprise;
use Illuminate\Console\Command;

class TempMemberRelation extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'temp:member-customer-relation {enterpriseId} {startTime} {endTime}';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '临时处理企业外部联系人数据';

    /**
     * 同步某个企业下某个时间内的外部联系人数据
     *
     * @throws \Exception
     */
    public function handle()
    {
        $enterpriseId = $this->argument('enterpriseId');
        if (empty($enterpriseId)) {
            die("请输入企业参数");
        }
        // 查询企业信息
        try {
            $enterprise = (new Enterprise())->detail(['enterprise_id' => $enterpriseId])['data'];
        } catch (Error $error) {
            die("查询企业信息失败：{$error->getCode()}:{$error->getMessage()}");
        }

        $startTime = $this->argument('startTime');
        if (empty($startTime)) {
            $startTime = Carbon::today()->timestamp;
        }

        $endTime = $this->argument('endTime');
        if (empty($endTime)) {
            $endTime = time();
        }
        $customerRelations = $this->yieldMemberCustomerRelation($enterprise['domain'], $startTime, $endTime);
        $total = $customerRelations->current()['total'];
        if ($total <= 0) {
            die("企业'{$enterprise['name']}'在" . date("Y-m-d H:i:s", $startTime) . "~" . date("Y-m-d H:i:s", $endTime) . "无可更新人员");
        }
        echo "更新企业'{$enterprise['name']}'的外部联系人,共计：{$total}\n";
        $memberCustomerRelationModel = MemberUserRelation::query();
        foreach ($customerRelations as $customerRelationLoop) {
            foreach ($customerRelationLoop['list'] as $customer) {
                $exist = $memberCustomerRelationModel->where('enterprise_id', $enterpriseId)
                    ->where('ec_id', $customer['ecId'])
                    ->where('mem_uid', $customer['memUid'])
                    ->first();
                if (empty($exist)) {
                    $relation = [
                        'ec_id'            => $customer['ecId'] ?? '',
                        'mem_uid'          => $customer['memUid'] ?? '',
                        'ec_avatar'        => $customer['ecAvatar'] ?? '',
                        'ec_gender'        => $customer['ecGender'] ?? 3,
                        'union_id'         => $customer['ecUnionid'] ?? '',
                        'source_type'      => $customer['emState'] ? intval($customer['emState']) : 1,
                        'ec_userid'        => $customer['ecUserid'] ?? '',
                        'ec_name'          => $customer['ecName'] ?? '',
                        'relation_at'      => $customer['emCreatetime'] ?? 0,
                        'user_name_index'  => $customer['ecIndex'] ?? '',
                        'user_remark_name' => $customer['emRemark'] ?? '',
                        'delete_type'      => 0,
                        'enterprise_id'    => $enterpriseId
                    ];

                    if ($customer['emRelationStatus'] == 1) {
                        $relation['delete_type'] = 0;
                    } else if ($customer['emRelationStatus'] == 3) {
                        $relation['delete_type'] = 2;
                        $relation['cancel_relation_at'] = $customer['emRelationDeleted'];
                    }
                    $memberCustomerRelationModel->create($relation);
                } else {
                    if ($customer['emRelationStatus'] == 1) {
                        $relation['delete_type'] = 0;
                    } else if ($customer['emRelationStatus'] == 3) {
                        $relation['delete_type'] = 2;
                        $relation['cancel_relation_at'] = $customer['emRelationDeleted'];
                    }

                    $memberCustomerRelationModel->update($relation);
                }
            }
        }

        die("更新完成~~~");
    }

    /**
     * 查询外部联系人数据
     *
     * @param string $domain
     * @param int $startTime
     * @param int $endTime
     * @return \Generator
     * @throws \Exception
     */
    private function yieldMemberCustomerRelation(string $domain, int $startTime, int $endTime)
    {
        $params = [
            'pageNum'  => Constant::PAGE,
            'pageSize' => 100,
            'emStart'  => $startTime,
            'emEnd'    => $endTime
        ];
        $relationSdk = new External();
        while (true) {
            $memberUserRelations = $relationSdk->relationList($domain, $params);
            if (empty($memberUserRelations['list'])) {
                break;
            }
            $params['pageNum'] += 1;
            yield $memberUserRelations;
        }
    }
}