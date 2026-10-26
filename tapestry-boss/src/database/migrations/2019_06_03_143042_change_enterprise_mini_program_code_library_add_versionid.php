<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class ChangeEnterpriseMiniProgramCodeLibraryAddVersionid extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::table('enterprise_mini_program_code_library', function (Blueprint $table) {
            // 增加更新字段
            $table->unsignedBigInteger('version_id')->default(0)->comment('关联shop版本id');
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('enterprise_mini_program_code_library', function (Blueprint $table) {
            //
        });
    }
}
