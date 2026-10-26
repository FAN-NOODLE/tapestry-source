<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class ChangeEnterpriseShopSettingAddcolumnTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::table('enterprise_shop_setting', function (Blueprint $table) {
            // 增加更新字段
            $table->string('app_secret')->default('')->after('corp_id')->comment('小程序secret');
            $table->string('app_id')->default('')->after('corp_id')->comment('小程序id');
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
            $table->dropColumn('app_secret');
            $table->dropColumn('app_id');
        });
    }
}
