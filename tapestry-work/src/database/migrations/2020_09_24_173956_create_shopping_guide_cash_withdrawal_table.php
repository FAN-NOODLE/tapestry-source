<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateShoppingGuideCashWithdrawalTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('shopping_guide_cash_withdrawal')) {
            Schema::create('shopping_guide_cash_withdrawal', function (Blueprint $table) {
                $table->unsignedBigInteger('id', true);
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedBigInteger('worker_id')->default(0)->comment('导购ID');
                $table->string('mem_uid',255)->default('')->comment('用户Id');
                $table->bigInteger('expect_cash_amount')->default(0)->comment('预计可提现金额');
                $table->bigInteger('available_cash_amount')->default(0)->comment('可提现金额');
                $table->unsignedInteger('cash_withdrawal')->default(0)->comment('已提现金额');
                $table->unsignedInteger('frozen_amount')->default(0)->comment('冻结金额');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            DB::statement("alter table `shopping_guide_cash_withdrawal` comment '导购可提现金额表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('shopping_guide_cash_withdrawal');
    }
}
