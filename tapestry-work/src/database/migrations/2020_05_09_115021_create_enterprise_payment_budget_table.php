<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateEnterprisePaymentBudgetTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('enterprise_payment_budget')) {
            Schema::create('enterprise_payment_budget', function (Blueprint $table) {
                $table->increments('id')->comment('自增主键');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->string('number', '50')->default('')->comment('活动编号');
                $table->unsignedInteger('amount')->default(0)->comment('金额');
                $table->unsignedInteger('balance')->default(0)->comment('余额');
                $table->string('voucher_url')->default('')->comment('充值凭证URL');
                $table->unsignedTinyInteger('income_type')->default(1)->comment('收入类型：1、收入，2、支出');
                $table->unsignedTinyInteger('type')->default(1)->comment('收支类型：1、充值，2、活动扣款，3、活动返款');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            // 表注释
            DB::statement("ALTER TABLE `enterprise_payment_budget` comment '充值收支记录表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('enterprise_payment_budget');
    }
}
