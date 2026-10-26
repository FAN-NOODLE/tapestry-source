<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddFieldFailStatusToEnterprisePayment extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('enterprise_payment')) {
            Schema::table('enterprise_payment', function (Blueprint $table) {
                $table->unsignedTinyInteger('fail_status')->default(0)->after('fail_time')->comment('失效状态，0：未失效，1：已失效');
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
        if (Schema::hasTable('enterprise_payment')) {
            Schema::table('enterprise_payment', function (Blueprint $table) {
                $table->dropColumn('fail_status');
            });
        }
    }
}
