<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class PartnerEnterpriseMap extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        //
        if (Schema::hasTable('partner_enterprise_map')) {
            if (!Schema::hasColumn('partner_enterprise_map', 'worker_avatar')) {
                Schema::table('partner_enterprise_map',function (Blueprint $table) {
                    $table->string('worker_avatar')->default('')->comment('导购头像');
                });
            }
        }

    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        //
        Schema::table('partner_enterprise_map', function (Blueprint $table) {
            $table->dropColumn('worker_avatar');
        });
    }
}
