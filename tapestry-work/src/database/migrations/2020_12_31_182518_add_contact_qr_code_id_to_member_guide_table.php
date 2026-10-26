<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddContactQrCodeIdToMemberGuideTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('member_guide')) {
            Schema::table('member_guide', function (Blueprint $table) {
                if (!Schema::hasColumn('member_guide', 'contact_qr_code_id')) {
                    $table->unsignedInteger('contact_qr_code_id')->default(0)->after('mem_uid')->comment('联系我二维码id');
                }
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
        if (Schema::hasTable('member_guide')) {
            Schema::table('member_guide', function (Blueprint $table) {
                if (Schema::hasColumn('member_guide', 'contact_qr_code_id')) {
                    $table->dropColumn('contact_qr_code_id');
                }
            });
        }
    }
}
