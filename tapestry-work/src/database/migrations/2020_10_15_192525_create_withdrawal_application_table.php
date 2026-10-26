<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateWithdrawalApplicationTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('withdrawal_application')) {
            Schema::create('withdrawal_application', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->comment('提现编号');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业ID');
                $table->unsignedBigInteger('operator_id')->default(0)->comment('操作人ID');
                $table->unsignedBigInteger('worker_id')->default(0)->comment('导购id');
                $table->unsignedInteger('application_amount')->default(0)->comment('申请金额');
                $table->unsignedInteger('request_initiated_at')->default(0)->comment('请求发起时间');
                $table->unsignedInteger('request_handle_at')->default(0)->comment('请求处理时间');
                $table->unsignedTinyInteger('status')->default(0)->comment('状态:1待审批2提现成功3提现失败4已驳回5转账中');
                $table->string('reason')->default('')->comment('原因:失败原因/驳回原因');

                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');

                $table->index('enterprise_id', 'enterprise_id_index');
            });
            \Illuminate\Support\Facades\DB::statement("ALTER TABLE `withdrawal_application` comment '提现申请表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('withdrawal_application');
    }
}
