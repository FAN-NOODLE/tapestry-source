<?php
/**
 * Created by PhpStorm.
 * User: XuHuitao
 * Date: 2020/7/13
 * Time: 10:21
 */

namespace App\Console\Commands;

use CrsSdk\Shop\Order;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;

class CountOrderFromShop extends Command {
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'count_order_from_shop {--enterprise=}';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '从gw_shop库统计订单数据同步到work库中';

    /**
     * Create a new command instance.
     */
    public function __construct()
    {
        parent::__construct();
    }

    public function handle() {

        $enterpriseId = $this->option('enterprise') ?? 0;

        set_time_limit(0);

        $page = 1;
        $size = 2000;
        $gwShop = config('client.api_url.shop') . '/gw-shop';
        config(['client.api_url.shop' => $gwShop]);
        while (1) {

            $existUser = DB::table('member_order_statistics')
                ->when($enterpriseId > 0, function ($query) use ($enterpriseId) {
                    $query->where('enterprise_id', $enterpriseId);
                })
                ->forPage($page, $size)
                ->select(['user_id'])->get()->toArray();
            // 已经存在的数据进行更新操作
            $userIdList = array_column($existUser, 'user_id');

            $orderData = (new Order())->statisticsUserData(['user_ids' => $userIdList]);
            foreach ($orderData as $data) {
                $updateData = [
                    'latest_buy_time' => $data['latest_buy_time']
                ];
                DB::table('member_order_statistics')
                    ->where('user_id', $data['user_id'])
                    ->update($updateData);
            }

            if (count($userIdList) == 0) {
                break;
            }
            $page++;
        }

    }
}