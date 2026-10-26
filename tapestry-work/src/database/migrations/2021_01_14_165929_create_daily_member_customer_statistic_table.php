<?php

use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateDailyMemberCustomerStatisticTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('daily_member_customer_statistic')) {
            Schema::create('daily_member_customer_statistic', function (Blueprint $table) {
                $table->bigIncrements('id')->comment('自增主键');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id');
                $table->unsignedInteger('stat_time')->default(0)->comment('统计日期，当日0点的时间戳');
                $table->char('mem_uid', 32)->default('')->comment('导购mem_uid');
                $table->unsignedMediumInteger('cumulative_customer_num')->default(0)->comment('累计顾客数');
                $table->unsignedMediumInteger('today_customer_num')->default(0)->comment('今日顾客数');
                $table->unsignedMediumInteger('new_apply_cnt')->default(0)->comment('向客户发起的好友申请数');
                $table->unsignedMediumInteger('new_contact_cnt')->default(0)->comment('新添加的客户数量');
                $table->unsignedMediumInteger('chat_cnt')->default(0)->comment('主动发送过消息的单聊总数');
                $table->unsignedMediumInteger('message_cnt')->default(0)->comment('在单聊中发送的消息总数');
                $table->unsignedSmallInteger('reply_percentage')->default(0)->comment('已回复聊天占比，回复过消息的聊天数/客户主动发起的聊天数比例，实际值*100');
                $table->unsignedMediumInteger('avg_reply_time')->default(0)->comment('平均首次回复时长');
                $table->unsignedMediumInteger('negative_feedback_cnt')->default(0)->comment('删除/拉黑成员的客户数');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');

            });
            // 表注释
            DB::statement("ALTER TABLE `daily_member_customer_statistic` comment '每日员工顾客统计数据'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('daily_member_customer_statistic');
    }
}
