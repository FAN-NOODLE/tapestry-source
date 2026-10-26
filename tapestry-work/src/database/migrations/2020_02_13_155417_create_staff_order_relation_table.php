<?php

use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateStaffOrderRelationTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('staff_order_relation')) {
            Schema::create('staff_order_relation', function (Blueprint $table) {
                $table->increments('id')->comment('自增主键');
                $table->unsignedInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->string('mem_uid', 32)->default('')->comment('员工ID');
                $table->string('order_no', 32)->default('')->comment('订单编号');
                $table->unsignedInteger('order_pay_time')->default(0)->comment('订单支付时间');
                $table->unsignedInteger('refund_money')->default(0)->comment('退款金额');
                $table->unsignedInteger('pay_money')->default(0)->comment('支付金额');
                $table->unsignedInteger('order_money')->default(0)->comment('订单金额');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
                $table->index('enterprise_id', 'enterprise_id');
                $table->index('mem_uid', 'mem_uid');
                $table->index('order_no', 'order_no');
            });

            // 表注释
            DB::statement("ALTER TABLE `staff_order_relation` comment '内购员工订单关系表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('staff_order_relation');
    }
}
