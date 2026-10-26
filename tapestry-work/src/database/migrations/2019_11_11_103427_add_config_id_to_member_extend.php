<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddConfigIdToMemberExtend extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('member_extend')) {
            Schema::table('member_extend', function (Blueprint $table) {
                $table->string('config_id')->default('')->after('enterprise_id')->comment('联系我二维码id');
                $table->string('contact_qr_code')->default('')->after('enterprise_id')->comment('联系我二维码');
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
        if (Schema::hasTable('member_extend')) {
            Schema::table('member_extend', function (Blueprint $table) {
                $table->dropColumn('config_id');
                $table->dropColumn('contact_qr_code');
            });
        }
    }
}
