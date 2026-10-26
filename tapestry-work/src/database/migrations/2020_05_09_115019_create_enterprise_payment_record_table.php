<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateEnterprisePaymentRecordTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('enterprise_payment_record')) {
            Schema::create('enterprise_payment_record', function (Blueprint $table) {
                $table->increments('id')->comment('自增主键');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedInteger('enterprise_payment_id')->default(0)->comment('企业微信红包表主键ID');
                $table->string('bill_no')->default('')->comment('商户单号');
                $table->string('detail_id')->default('')->comment('红包单号（微信返回）');
                $table->unsignedInteger('enterprise_payment_random_id')->default(0)->comment('随机金额主表ID');
                $table->string('union_id', 50)->default('')->comment('uc外部联系人union_id')->index();
                $table->unsignedInteger('amount')->default(0)->comment('金额');
                $table->string('mem_uid')->default('')->comment('发送人mem_uid');
                $table->unsignedBigInteger('user_id')->default(0)->comment('领取人user_id');
                $table->unsignedInteger('send_time')->default(0)->comment('红包发送时间');
                $table->unsignedInteger('receive_refund_time')->default(0)->comment('红包领取/退款时间');
                $table->unsignedTinyInteger('status')->default(1)->comment('发送领取状态：1、发送中，2、已领取，3、已退回');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            // 表注释
            DB::statement("ALTER TABLE `enterprise_payment_record` comment '企业微信红包表用户领取记录表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('enterprise_payment_record');
    }
}
