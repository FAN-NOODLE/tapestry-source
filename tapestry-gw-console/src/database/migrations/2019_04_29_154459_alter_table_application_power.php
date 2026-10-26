<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use \Illuminate\Support\Facades\DB;

class AlterTableApplicationPower extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('application_power')) {
            Schema::table('application_power', function (Blueprint $table) {
                if (!Schema::hasColumn('application_power', 'power_alias')) {
                    $table->string('power_alias', 128)->default('')->comment('前端使用活动别名')->after('available');
                }
            });

            $names = [
                '礼券' => 'myapp_app_myapplication_coupon',
                '拼团' => 'myapp_app_myapplication_group',
                '新人优惠' => 'myapp_app_myapplication_new',
                '秒杀' => 'myapp_app_myapplication_spike',
                '超低购' => 'myapp_app_myapplication_low-prices',
            ];

            foreach ($names as $name => $value) {
                $where['name'] = $name;
                DB::table("application_power")
                    ->where($where)
                    ->update(['power_alias' => $value]);
            }
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('application_power', function (Blueprint $table) {
            //
        });
    }
}
