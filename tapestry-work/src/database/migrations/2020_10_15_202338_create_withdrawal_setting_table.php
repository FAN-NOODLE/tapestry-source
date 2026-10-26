<?php

use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateWithdrawalSettingTable extends Migration
{
    /**
     * 提现设置
     */
    public function up()
    {
        if (!Schema::hasTable('withdrawal_setting')) {
            Schema::create('withdrawal_setting', function (Blueprint $table) {
                $table->increments('id')->primary();
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业ID');
                $table->string('transfer_no')->default('')->comment('转账项目编号');
                $table->string('type')->default('')->comment('导购类型：1、门店导购，2、社群导购，3、企业导购 ,分隔');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });
            DB::statement("ALTER TABLE `withdrawal_setting` comment '提现设置'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('withdrawal_setting');
    }
}
