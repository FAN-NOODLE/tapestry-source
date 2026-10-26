<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddFieldToServiceMerchantSetting extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('service_merchant_setting')) {
            Schema::table('service_merchant_setting', function (Blueprint $table) {
                $table->unsignedTinyInteger('is_default')->default(0)->after('status')->comment('是否默认 0 否 1 是');
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
        if (Schema::hasTable('service_merchant_setting')) {
            Schema::table('service_merchant_setting', function (Blueprint $table) {
                $table->dropColumn('is_default');
            });
        }
    }
}
