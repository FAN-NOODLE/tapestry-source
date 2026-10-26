<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateGroupOrderTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('group_order')) {
            Schema::create('group_order', function (Blueprint $table) {
                $table->increments('id')->comment('自增主键');
                $table->unsignedInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->string('qw_group_id', 32)->default('')->comment('微信群ID');
                $table->string('group_user_id', 50)->default('')->comment('用户ID');
                $table->string('union_id', 50)->default('')->comment('用户唯一ID');
                $table->string('order_no', 32)->default('')->comment('订单编号');
                $table->unsignedInteger('order_pay_time')->default(0)->comment('订单支付时间');
                $table->unsignedInteger('refund_money')->default(0)->comment('退款金额');
                $table->unsignedInteger('pay_money')->default(0)->comment('支付金额');
                $table->unsignedInteger('order_money')->default(0)->comment('订单金额');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            // 表注释
            DB::statement("ALTER TABLE `group_order` comment '企业微信群订单表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('group_order');
    }
}
