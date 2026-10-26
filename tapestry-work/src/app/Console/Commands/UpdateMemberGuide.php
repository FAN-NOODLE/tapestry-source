<?php
/**
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/2/28
 * Time: 16:32
 */

namespace App\Console\Commands;

use App\Library\Uc\Member;
use App\Models\MemberExtend;
use App\Models\MemberGuide;
use App\Models\MemberUserRelation;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class UpdateMemberGuide extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'update-member-guide {type}';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '转移导购数据(恰恰,蒙牛传递参数1,其他传递2)';

    public function handle()
    {
        $type = $this->argument('type') ?? 0;
        if (empty($type) || !in_array($type, [1, 2])) {
            echo "恰恰,蒙牛传递参数1,其他传递2";
            exit;
        }
        $enterpriseIds = [
            13,
            86
        ];

        $model = MemberExtend::query();
        if ($type == 1) {
            $model->whereIn('enterprise_id', $enterpriseIds);
        }
        if ($type == 2) {
            $model->whereNotIn('enterprise_id',  $enterpriseIds);
        }

        // 获取支付导购数据
        $MemberExtend = $model->where('role', '导购')
            ->get()->toArray();

        $data = [];
        // 循环数据
        foreach ($MemberExtend as $value) {

            $data[] = [
                'enterprise_id'      => $value['enterprise_id'],
                'mem_uid'            => $value['mem_uid'],
                'guide_id'           => $value['guide_id'],
                'allocation_status'  => $value['allocation_status'],
                'trade_guide_status' => $value['trade_guide_status'],
                'limit_user_count'   => $value['limit_user_count'],
                'created_at'         => $value['created_at']
            ];
        }

        // 分割数据
        $list = array_chunk($data, 500);

        // 循环数据
        foreach ($list as $item) {
            MemberGuide::query()->insert($item);
        }
    }

}
