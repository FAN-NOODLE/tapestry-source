<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateSettlementRefundRecordTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('settlement_refund_record')) {
            Schema::create('settlement_refund_record', function (Blueprint $table) {
                $table->increments('id')->comment('主键id');
                $table->string('order_random_code',32)->default('')->comment('结算单号');
                $table->string('settlement_code',32)->default('')->comment('第三方结算单号');
                $table->string('project_code',32)->default('')->comment('项目编码');
                $table->Integer('refund_merchant_amount')->default(0)->comment('余额退回金额');
                $table->Integer('refund_service_amount')->default(0)->comment('服务费抵扣券退回金额');
                $table->string('change_code',32)->default('')->comment('服务费变更单单流水号');
                $table->Integer('change_merchant_amount')->default(0)->comment('变更单余额退回金额');
                $table->Integer('change_service_amount')->default(0)->comment('变更单抵扣券退回金额');
                $table->unsignedTinyInteger('refund_status')->default(1)->comment('退款状态, 1=退款中 2=退款成功 3=退款失败');
                $table->unsignedInteger('refund_initiated_at')->default(0)->comment('退款发起时间，退款接口调用陈成功才有值');
                $table->unsignedInteger('refund_succeed_at')->default(0)->comment('退款成功时间');
                $table->unsignedSmallInteger('retry_num')->default(0)->comment('退款重试次数，达到重试次数则不再重试');
                $table->unsignedTinyInteger('refund_failed_type')->default(0)->comment('退款失败类型');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');

                $table->index('order_random_code');
                $table->index('settlement_code');
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
        Schema::dropIfExists('settlement_refund_record');
    }
}
