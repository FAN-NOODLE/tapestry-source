<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateCashCouponRecordTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('cash_coupon_record')) {
            Schema::create('cash_coupon_record', function (Blueprint $table) {
                $table->increments('id')->comment('自增主键');
                $table->unsignedInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedInteger('cash_coupon_id')->default(0)->comment('现金红包活动id');
                $table->string('bill_no')->default('')->comment('商户单号');
                $table->string('detail_id')->default('')->comment('红包单号（微信返回）');
                $table->string('mem_uid')->default('')->comment('发送人mem_uid');
                $table->string('ec_userid')->default('')->comment('领取人ec_userid');
                $table->string('open_id')->default('')->comment('领取人open_id');
                $table->unsignedInteger('send_time')->default(0)->comment('红包发送时间');
                $table->unsignedInteger('receive_refund_time')->default(0)->comment('红包领取/退款时间');
                $table->unsignedTinyInteger('status')->default(1)->comment('发送领取状态：1、发送中，2、已领取，3、已退回');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            // 表注释
            DB::statement("ALTER TABLE `cash_coupon_record` comment '现金红包活动对账记录表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('cash_coupon_record');
    }
}
