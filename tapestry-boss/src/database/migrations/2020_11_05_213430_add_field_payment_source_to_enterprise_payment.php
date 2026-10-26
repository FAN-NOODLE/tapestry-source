<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddFieldPaymentSourceToEnterprisePayment extends Migration
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
                $table->string('payment_source', 100)->default('')->after('payment_time')->comment('支付来源');
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
