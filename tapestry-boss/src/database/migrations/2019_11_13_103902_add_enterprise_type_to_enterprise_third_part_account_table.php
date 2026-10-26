<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddEnterpriseTypeToEnterpriseThirdPartAccountTable extends Migration
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
                $table->tinyInteger('enterprise_type')->default(1)->after('release_tag')->comment('小程序发版类型(1：稳定版，2：单独的项目)');
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
        Schema::table('enterprise_third_part_account', function (Blueprint $table) {
            //
        });
    }
}
