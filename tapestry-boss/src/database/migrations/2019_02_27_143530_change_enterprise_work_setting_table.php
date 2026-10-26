<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class ChangeEnterpriseWorkSettingTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::table('enterprise_work_setting', function (Blueprint $table) {
            // 增加更新字段
            $table->string('ep_enumber')->default('')->after('corp_id')->comment('UC返回');
            $table->string('pl_plugin_id')->default('')->after('corp_id')->comment('UC返回');
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('enterprise_work_setting', function (Blueprint $table) {
            //
        });
    }
}
