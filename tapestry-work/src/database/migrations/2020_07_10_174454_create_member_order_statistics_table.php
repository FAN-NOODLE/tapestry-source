<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateMemberOrderStatisticsTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('member_order_statistics', function (Blueprint $table) {
            $table->unsignedBigInteger('id', true)->comment('主键id');
            $table->string('union_id', 80)->default('')->comment('用户唯一ID')->index();
            $table->unsignedBigInteger('user_id')->default(0)->comment('用户id')->index();
            $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
            $table->unsignedInteger('total_order_num')->default(0)->comment('用户累计订单数')->index();
            $table->unsignedInteger('total_order_amount')->default(0)->comment('用户累计金额，单位：分')->index();
            $table->text('last_update_data')->default('')->comment('最近更新数据');
            $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
            $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
            $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
        });
        // 表注释
        DB::statement("ALTER TABLE `member_order_statistics` comment '顾客订单统计表'");
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('member_order_statistics');
    }
}
