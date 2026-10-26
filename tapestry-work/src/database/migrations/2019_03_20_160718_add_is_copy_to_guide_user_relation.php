<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddIsCopyToGuideUserRelation extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('guide_user_relation')) {
            Schema::table('guide_user_relation', function (Blueprint $table) {
                $table->unsignedTinyInteger('is_copy_mobile')->default(0)->after('mem_uid')->comment('是否复制过手机号(0、否,1、是)');
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
        if (Schema::hasTable('guide_user_relation')) {
            Schema::table('guide_user_relation', function (Blueprint $table) {
                $table->dropColumn('is_copy_mobile');
            });

        }
    }
}
