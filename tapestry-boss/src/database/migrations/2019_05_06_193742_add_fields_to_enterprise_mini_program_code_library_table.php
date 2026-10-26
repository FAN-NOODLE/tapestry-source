<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddFieldsToEnterpriseMiniProgramCodeLibraryTable extends Migration
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
                $table->unsignedTinyInteger('release_model')
                    ->default(1)
                    ->after('online_type')
                    ->comment('发布模式（1：全量发布，2：分阶段发布，15天内全量发布）');
                $table->unsignedInteger('grp_gray_percentage')
                    ->default(0)
                    ->after('release_model')
                    ->comment('分阶段发布灰度的百分比，1到100的整数，当release_model=2时有值');
                $table->unsignedTinyInteger('grp_status')
                    ->default(0)
                    ->after('grp_gray_percentage')
                    ->comment('分阶段发布状态（0:初始状态 1:执行中 2:暂停中 3:执行完毕 4:被删除）');

                $table->index('grp_status', 'idx_grp_status');
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
                $table->dropColumn('release_model');
                $table->dropColumn('grp_gray_percentage');
                $table->dropColumn('grp_status');
                $table->dropIndex('idx_grp_status');
            });
        }
    }
}
