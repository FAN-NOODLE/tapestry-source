<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddShortNameToEnterpriseThirdPartAccountTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::table('enterprise_third_part_account', function (Blueprint $table) {
            $table->string('short_name', 191)->default('')->after('logo')->comment('商户简称');
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('enterprise_third_part_account', function (Blueprint $table) {
            $table->dropColumn('short_name');
        });
    }
}
