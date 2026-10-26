<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddFieldOpenPlatformIdToEnterprise extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('enterprise')) {
            Schema::table('enterprise', function (Blueprint $table) {
                $table->unsignedBigInteger('open_platform_id')->default(0)->after('service_merchant_id')->comment('企业关联的开放平台id');
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
        if (Schema::hasTable('enterprise')) {
            Schema::table('enterprise', function (Blueprint $table) {
                $table->dropColumn('open_platform_id');
            });
        }
    }
}
