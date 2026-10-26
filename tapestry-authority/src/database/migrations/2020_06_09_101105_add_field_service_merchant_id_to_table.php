<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddFieldServiceMerchantIdToTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasColumn('enterprise', 'service_merchant_id')) {
            Schema::table('enterprise', function (Blueprint $table) {
                $table->unsignedBigInteger('service_merchant_id')->default(1)->after('available_system')->comment('服务商ID');
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
        Schema::table('enterprise', function (Blueprint $table) {
            $table->dropColumn('service_merchant_id');
        });
    }
}
