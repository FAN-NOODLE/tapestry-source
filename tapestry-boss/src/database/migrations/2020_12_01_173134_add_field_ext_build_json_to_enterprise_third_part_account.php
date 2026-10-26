<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddFieldExtBuildJsonToEnterpriseThirdPartAccount extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('enterprise_third_part_account')) {
            Schema::table('enterprise_third_part_account', function (Blueprint $table) {
                $table->string('ext_build_json', 255)->default('')->after('enterprise_type')->comment('小程序打包发版信息');
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
        if (Schema::hasTable('enterprise_third_part_account')) {
            Schema::table('enterprise_third_part_account', function (Blueprint $table) {
                $table->dropColumn('ext_build_json');
            });
        }
    }
}
