<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddFieldsStaffSetting extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('staff_setting')) {
            Schema::table('staff_setting', function (Blueprint $table) {
                $table->unsignedTinyInteger('fission_type')->default(1)->after('achievement_desc')->comment('客户分享是否统计导购业绩 1:不计算业绩 2:计算业绩');
                $table->unsignedTinyInteger('type')->default(1)->after('fission_type')->comment('导购类型 1:内购，2:导购');
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
        Schema::table('staff_setting', function (Blueprint $table) {
            //
        });
    }
}
