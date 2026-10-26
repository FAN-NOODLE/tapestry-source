<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddAgentIdEnterpriseShopSettingTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::table('enterprise_shop_setting', function (Blueprint $table) {
            $table->string('agent_id')->default('')->after('corp_id')->comment('agent_id');
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('enterprise_shop_setting', function (Blueprint $table) {
            $table->dropColumn('agent_id');
        });
    }
}
