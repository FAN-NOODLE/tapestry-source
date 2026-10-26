<?php

use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateDailyWorkerCustomerStatisticTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('daily_worker_customer_statistic')) {
            Schema::create('daily_worker_customer_statistic', function (Blueprint $table) {
                $table->bigIncrements('id')->comment('自增主键');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id');
                $table->unsignedBigInteger('worker_id')->default(0)->comment('导购id');
                $table->char('mem_uid', 32)->default('')->comment('导购mem_uid');
                $table->unsignedInteger('stat_time')->default(0)->comment('统计日期，当日0点的时间戳');
                $table->unsignedMediumInteger('cumulative_contact_card_way_customer_num')->default(0)->comment('累计通过名片方式添加的顾客总数');
                $table->unsignedMediumInteger('cumulative_other_way_customer_num')->default(0)->comment('累计通过其他方式添加的顾客总数');
                $table->unsignedMediumInteger('cumulative_binding_customer_num')->default(0)->comment('累计绑定顾客数');
                $table->unsignedMediumInteger('today_contact_card_way_customer_num')->default(0)->comment('今日通过名片方式添加的顾客总数');
                $table->unsignedMediumInteger('today_other_way_customer_num')->default(0)->comment('今日通过其他方式添加的顾客总数');
                $table->unsignedMediumInteger('today_binding_customer_num')->default(0)->comment('今日绑定顾客数');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');

            });
            // 表注释
            DB::statement("ALTER TABLE `daily_worker_customer_statistic` comment '每日导购顾客统计数据'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('daily_worker_customer_statistic');
    }
}
