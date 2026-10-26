<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class EnterpriseSmsStatistics extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('enterprise_sms_statistics')) {
            Schema::create('enterprise_sms_statistics', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary()->comment('自增id');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedBigInteger('sms_number')->default(0)->comment('短信可用数量');
                $table->unsignedBigInteger('sms_total_number')->default(0)->comment('短信累计数量');
                $table->unsignedBigInteger('sms_payment_total')->default(0)->comment('短信累计充值金额');
                $table->unsignedInteger('updated_at')->default(0);
                $table->unsignedInteger('created_at')->default(0);
                $table->unsignedInteger('deleted_at')->default(0);
            });

            DB::statement("ALTER TABLE `enterprise_sms_statistics` COMMENT '企业短信统计表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('enterprise_sms_statistics');
    }
}
