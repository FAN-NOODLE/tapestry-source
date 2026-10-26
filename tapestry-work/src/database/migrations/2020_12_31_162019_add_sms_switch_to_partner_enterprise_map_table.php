<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddSmsSwitchToPartnerEnterpriseMapTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('partner_enterprise_map')) {
                if (!Schema::hasColumn('partner_enterprise_map', 'sms_switch')) {
                    Schema::table('partner_enterprise_map', function (Blueprint $table) {
                        $table->tinyInteger('sms_switch')->default('1')->comment('合伙人是否接受短信 默认 1开启 2关闭');
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
        Schema::table('partner_enterprise_map', function (Blueprint $table) {
            $table->dropColumn('sms_switch');
        });
    }
}
