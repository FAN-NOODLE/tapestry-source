<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddIndexToEnterpriseMiniProgramCodeLibraryTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('enterprise_mini_program_code_library')) {
            Schema::table('enterprise_mini_program_code_library', function (Blueprint $table) {
                $table->index('audit_status', 'idx_audit_status');
                $table->index('version', 'idx_version');
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
        if (Schema::hasTable('enterprise_mini_program_code_library')) {
            Schema::table('enterprise_mini_program_code_library', function (Blueprint $table) {
                $table->dropIndex('idx_audit_status');
                $table->dropIndex('idx_version');
            });
        }
    }
}
