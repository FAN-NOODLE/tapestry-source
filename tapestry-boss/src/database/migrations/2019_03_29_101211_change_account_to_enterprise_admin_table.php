<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class ChangeAccountToEnterpriseAdminTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('enterprise_admin')) {
            Schema::table('enterprise_admin', function (Blueprint $table) {
                $table->dropUnique('account_UNIQUE');
                $table->index('account');
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
        Schema::table('enterprise_admin', function (Blueprint $table) {
            //
        });
    }
}
