<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddTagToEnterpriseMiniProgramCodeLibraryTable extends Migration
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
                $table->string('release_tag', 50)->default('')->after('grp_status')->comment('小程序发布版本分支tag');
                $table->string('release_branch', 50)->default('')->after('grp_status')->comment('小程序发布版本分支');
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
                $table->dropColumn('release_tag');
                $table->dropColumn('release_branch');
            });
        }
    }
}
