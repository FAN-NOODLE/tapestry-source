<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddFieldToEnterpriseMessageScene extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('enterprise_message_scene')) {
            Schema::table('enterprise_message_scene', function (Blueprint $table) {
                $table->unsignedInteger('subscribe_num')->default(0)->comment('已订阅数');
                $table->unsignedInteger('send_num')->default(0)->comment('已发送数');
                $table->unsignedInteger('send_real_num')->default(0)->comment('实际已发送数');
                $table->unsignedTinyInteger('type')->default(1)->comment('策略类型 1 自动 2 手动');
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
        //
    }
}
