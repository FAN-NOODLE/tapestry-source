<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddFieldIsRenewToEnterprisePayment extends Migration
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
                $table->unsignedBigInteger('parent_id')->default(0)->after('enterprise_id')->comment('续费自关联id');
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
                $table->dropColumn('parent_id');
            });
        }
    }
}
