<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class EnterprisePayment extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('enterprise_payment')) {
            Schema::create('enterprise_payment', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary()->comment('自增id');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedTinyInteger('payment_category')->default(0)->comment('购买类目，0：红包，1：短信，2：智慧导购，3：SCRM，4：商城，5：综合系统');
                $table->unsignedBigInteger('payment_amount')->default(0)->comment('支付金额（单位：分）');
                $table->unsignedBigInteger('available_amount')->default(0)->comment('可用金额（单位：分）针对红包');
                $table->unsignedBigInteger('collection_amount')->default(0)->comment('收款金额（单位：分）');
                $table->unsignedBigInteger('server_charge')->default(0)->comment('服务费（单位：分）');
                $table->unsignedBigInteger('buy_number')->default(0)->comment('购买账号数量');
                $table->unsignedBigInteger('buy_day')->default(0)->comment('购买时长');
                $table->unsignedInteger('payment_time')->default(0)->comment('支付时间')->index();
                $table->unsignedInteger('use_time')->default(0)->comment('生效时间')->index();
                $table->unsignedInteger('fail_time')->default(0)->comment('失效时间')->index();
                $table->unsignedInteger('arrival_time')->default(0)->comment('到账时间');
                $table->unsignedTinyInteger('payment_type')->default(0)->comment('支付方式，0：其他，1：现金，2：微信，3：支付宝，4：银联，5：对公转账');
                $table->unsignedTinyInteger('payment_status')->default(0)->comment('收款状态，0：待收款，1：已收款');
                $table->unsignedTinyInteger('account_sync')->default(0)->comment('账号同步，0：未同步，1：已同步');
                $table->string('payment_order_sn', 255)->default('')->comment('支付流水号');
                $table->string('voucher_url', 255)->default('')->comment('购买凭证');
                $table->string('remark', 100)->default('')->comment('备注');
                $table->unsignedInteger('updated_at')->default(0);
                $table->unsignedInteger('created_at')->default(0);
                $table->unsignedInteger('deleted_at')->default(0);
            });

            DB::statement("ALTER TABLE `enterprise_payment` COMMENT '企业付费表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('enterprise_payment');
    }
}
