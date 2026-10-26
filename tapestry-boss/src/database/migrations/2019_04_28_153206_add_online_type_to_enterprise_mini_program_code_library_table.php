<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddOnlineTypeToEnterpriseMiniProgramCodeLibraryTable extends Migration
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
                $table->unsignedTinyInteger('online_type')->default(1)->after('using')->comment('上线方式（1：自动上线，2：手动上线）');
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
                $table->dropColumn('online_type');
            });
        }
    }
}
