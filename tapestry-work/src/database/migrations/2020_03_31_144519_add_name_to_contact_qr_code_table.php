<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddNameToContactQrCodeTable extends Migration
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
                $table->string('name')->default('')->after('enterprise_id')->comment('二维码名称');
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
                $table->dropColumn('name');
            });
        }
    }
}
