<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateSettlementRecordTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('settlement_record')) {
            Schema::create('settlement_record', function (Blueprint $table) {
                $table->increments('id')->comment('主键id');
                $table->string('order_random_code',32)->default('')->comment('结算单号');
                $table->string('settlement_code',32)->default('')->comment('第三方结算单号');
                $table->string('account_num',32)->default('')->comment('账单号');
                $table->string('project_code',32)->default('')->comment('项目编码');
                $table->string('name',20)->default('')->comment('姓名');
                $table->unsignedTinyInteger('certificate_type')->default(0)->comment('证件类型');
                $table->string('certificate_num',32)->default('')->comment('证件号');
                $table->string('phone_num',20)->default('')->comment('手机号');
                $table->string('bank_name',20)->default('')->comment('银行名称');
                $table->string('bank_card_num',20)->default('')->comment('银行卡号');
                $table->unsignedInteger('settle_amount')->default(0)->comment('结算金额');
                $table->integer('status')->default(0)->comment('结算单状态');
                $table->tinyInteger('is_verification')->default(0)->comment('结算单校验是否成功,0=否1=是');
                $table->string('verification_error_code',20)->default('')->comment('错误码');
                $table->string('fail_reason',64)->default('')->comment('失败原因');
                $table->tinyInteger('hangup_flag')->default(0)->comment('挂单标识: true:挂单');
                $table->string('hangup_begintime',64)->default('')->comment('挂单时间');
                $table->string('hangup_msg',64)->default('')->comment('挂单原因');
                $table->string('pay_time',32)->default('')->comment('支付时间(扣除商户余额的时间)');
                $table->string('handle_time',32)->default('')->comment('处理时间(打款成功或失败的时间)');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');

                $table->index('order_random_code');
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
        Schema::dropIfExists('settlement_record');
    }
}
