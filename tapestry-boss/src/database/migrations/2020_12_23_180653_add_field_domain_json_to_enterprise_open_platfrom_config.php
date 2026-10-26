<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddFieldDomainJsonToEnterpriseOpenPlatfromConfig extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('enterprise_open_platform_config')) {
            Schema::table('enterprise_open_platform_config', function (Blueprint $table) {
                $table->text('domain_json')->default('')->after('ext_build_json')->comment('小程序域名配置');
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
        if (Schema::hasTable('enterprise_open_platform_config')) {
            Schema::table('enterprise_open_platform_config', function (Blueprint $table) {
                $table->dropColumn('domain_json');
            });
        }
    }
}
