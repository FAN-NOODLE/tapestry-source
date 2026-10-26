<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddPlPluginIdEnterpriseWorkPluginTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::table('enterprise_work_plugin', function (Blueprint $table) {
            $table->string('pl_plugin_id')->default('')->after('check_file')->comment('uc返回id');
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('enterprise_work_plugin', function (Blueprint $table) {
            $table->dropColumn('pl_plugin_id');
        });
    }
}
