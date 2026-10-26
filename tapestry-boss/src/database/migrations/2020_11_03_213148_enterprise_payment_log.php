<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class EnterprisePaymentLog extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('enterprise_payment_log')) {
            Schema::create('enterprise_payment_log', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary()->comment('自增id');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedBigInteger('enterprise_payment_id')->default(0)->comment('企业付费id')->index();
                $table->string('old_data_json')->default('')->comment('旧数据json');
                $table->string('new_data_json')->default('')->comment('新数据json');
                $table->unsignedInteger('updated_at')->default(0);
                $table->unsignedInteger('created_at')->default(0);
                $table->unsignedInteger('deleted_at')->default(0);
            });

            DB::statement("ALTER TABLE `enterprise_payment_log` COMMENT '企业付费日志表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('enterprise_payment_log');

    }
}
