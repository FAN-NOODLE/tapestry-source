<?php

use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateWithdrawalAutomaticReviewTable extends Migration
{
    /**
     * 自动批准设置
     */
    public function up()
    {
        if (!Schema::hasTable('withdrawal_automatic_review')) {
            Schema::create('withdrawal_automatic_review', function (Blueprint $table) {
                $table->increments('id')->primary();
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业ID');
                $table->unsignedTinyInteger('auto_status')->default(0)->comment('是否开启：0 未开启，1已开启 ');
                $table->string('auto_at')->default('')->comment('自动批准时间');
                $table->unsignedInteger('last_at')->default(0)->comment('最后一次执行时间');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });
            DB::statement("ALTER TABLE `withdrawal_automatic_review` comment '自动批准设置'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('withdrawal_automatic_review');
    }
}
