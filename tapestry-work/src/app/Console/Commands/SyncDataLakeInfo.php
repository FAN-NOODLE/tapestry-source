<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/3/30 0030
 * Time: 10:14
 */

namespace App\Console\Commands;


use App\Library\DataLake\UserRequestService;
use App\Models\MemberUserRelation;
use App\Models\User;
use Illuminate\Console\Command;

class SyncDataLakeInfo extends Command
{
    /**
     * 单次查询数量
     */
    const PER_PAGE = 500;

    /**
     * 命令
     * @var string
     */
    protected $signature = 'sync_datalake_info';

    /**
     * 描述信息
     * @var string
     */
    protected $description = '同步DL会员信息';

    /**
     * @var UserRequestService
     */
    protected $dataLakeUserService;

    /**
     * SyncDataLakeInfo constructor.
     */
    public function __construct()
    {
        parent::__construct();
        $this->dataLakeUserService = new UserRequestService();
        MemberUserRelation::query();
    }

    /**
     * 入口方法
     */
    public function handle()
    {
        $num = 0;
        foreach ($this->getUserGenerator() as $users) {
            foreach ($users as $user) {

                if (!$this->integrateInfo($user)) {
                    continue;
                }
                $num++;
            }
        }

        echo "本次同步".$num."条数据";
    }

    /**
     * 更新用户信息
     *
     * @param array $user
     * @return bool
     */
    private function integrateInfo(array $user): bool
    {
        $this->dataLakeUserService->setEnterpriseId($user['enterprise_id']);
        try {
            $memberInfo = $this->dataLakeUserService->queryMemberInfo('unionId', $user['union_id']);
            $point      = $this->dataLakeUserService->getCurrentPoint('unionId', $user['union_id']);
            $grade      = $this->dataLakeUserService->queryGradeInfo('unionId', $user['union_id']);

        } catch (\Throwable $e) {
            // 非会员
            return false;
        }
        $update['union_code'] = $memberInfo['union_code'];
        $update['mobile']     = $memberInfo['mobile'];
        $update['birthday']   = $memberInfo['birthday'];
        $update['point']      = $point['points'][0]['evalidPoints'] ?? 0;
        $update['grade']      = $grade['grade'];

        $effectRows = MemberUserRelation::query()
            ->where(['union_id' => $user['union_id']])
            ->update($update);
        echo "修改union_id".$user['union_id']."的用户".$effectRows."行\n";
        return true;

    }

    /**
     * 测试数据
     *
     * @return \Generator
     */
    private function fakeUserGen(): \Generator
    {
        yield [
            [
                "enterprise_id" => 1,
                "union_id"      => "ouOt506tsM-ybmk7w0LZ5rUI10zY"

            ],
            [
                "enterprise_id" => "3306120225798016",
                "union_id"      => "oUmNp1nMJ59vXzIFRfz6tTO2F8gY"

            ]
        ];
    }

    /**
     * member_user_relation获取用户列表，按union_id分组
     * @return \Generator
     */
    private function getUserGenerator(): \Generator
    {
        $page = 1;
        do {
            $res = MemberUserRelation::query()
                ->where('union_id','<>','')
//                ->where('union_id','oUmNp1p1kl9BHb542g9VnSw6m-Mo')
                ->groupBy('union_id')
                ->forPage($page, static::PER_PAGE)
                ->get(['enterprise_id', 'union_id'])
                ->toArray();
            yield $res;
            $page++;
        } while ($res !== []);
    }
}