<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

/**
 * 添加app_type字段到withdrawal_application表
 * Created by PhpStorm.
 * User: Hui Yang
 * Email: 2361313833@qq.com
 * DateTime: 2020年11月24日11:21:38
 * Class AddMemTypeInStoreWorkerCertificationTable
 */
class AddAppTypeInWithdrawalApplicationTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        //
        if (Schema::hasTable('withdrawal_application')) {
            Schema::table('withdrawal_application', function (Blueprint $table) {
                if (!Schema::hasColumn('withdrawal_application', 'app_type')) {
                    $table->tinyInteger('app_type')
                        ->default(1)
                        ->after('enterprise_id')
                        ->comment('身份类型：1:导购(默认) 2:合伙人');
                }
            });
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        //
    }
}
