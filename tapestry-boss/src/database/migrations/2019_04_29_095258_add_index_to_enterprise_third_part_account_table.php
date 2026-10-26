<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddIndexToEnterpriseThirdPartAccountTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('enterprise_third_part_account')) {
            Schema::table('enterprise_third_part_account', function (Blueprint $table) {
                $table->index('type', 'idx_type');
                $table->index('name', 'idx_name');
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
        if (Schema::hasTable('enterprise_third_part_account')) {
            Schema::table('enterprise_third_part_account', function (Blueprint $table) {
                $table->dropIndex('idx_type');
                $table->dropIndex('idx_name');
            });
        }
    }
}
