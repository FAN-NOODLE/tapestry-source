<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddSourceTypeToContactQrCodeTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('contact_qr_code')) {
            Schema::table('contact_qr_code', function (Blueprint $table) {
                $table->unsignedBigInteger('channel_id')->default(0)->after('enterprise_id')->comment('对应渠道id,0=默认正常创建');
                $table->unsignedTinyInteger('source_type')->default(0)->after('type')->comment('建立方式，0=默认正常创建，1=内购自建,2=门店商城');
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
        if (Schema::hasTable('contact_qr_code')) {
            Schema::table('contact_qr_code', function (Blueprint $table) {
                $table->dropColumn('source_type');
            });
        }
    }
}
