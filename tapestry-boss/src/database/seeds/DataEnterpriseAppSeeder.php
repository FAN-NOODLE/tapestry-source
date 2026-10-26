<?php

use Illuminate\Database\Seeder;
use App\Models\EnterpriseApp;

class DataEnterpriseAppSeeder extends Seeder
{
    /**
     * Run the database seeds.
     *
     * @return void
     */
    public function run()
    {
        $data = [
            [
                'identity' => 'coupon', // 礼劵
                'created_at' => time()
            ],
            [
                'identity' => 'novice', // 新人优惠
                'created_at' => time()
            ],
            [
                'identity' => 'sec_kill', // 秒杀限时购
                'created_at' => time()
            ],
            [
                'identity' => 'group', // 拼团
                'created_at' => time()
            ],
        ];

        // 新增默认企业APP
        try {
            EnterpriseApp::query()->insert($data);
        } catch (\Exception $exception) {

        }
    }
}
