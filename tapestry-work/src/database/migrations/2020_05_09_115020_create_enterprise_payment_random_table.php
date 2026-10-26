<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateEnterprisePaymentRandomTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('enterprise_payment_random')) {
            Schema::create('enterprise_payment_random', function (Blueprint $table) {
                $table->increments('id')->comment('自增主键');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedInteger('enterprise_payment_id')->default(0)->comment('企业微信红包表主键ID');
                $table->unsignedInteger('is_receive')->default(1)->comment('是否占用 1=否 2=是');
                $table->unsignedInteger('amount')->default(0)->comment('金额');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            // 表注释
            DB::statement("ALTER TABLE `enterprise_payment_random` comment '企业微信红包表随机金额表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('enterprise_payment_random');
    }
}
