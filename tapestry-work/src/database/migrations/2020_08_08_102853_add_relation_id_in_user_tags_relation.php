<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddRelationIdInUserTagsRelation extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        //
        if (Schema::hasTable('user_customize_tags_relation')) {
            Schema::table('user_customize_tags_relation', function (Blueprint $table) {
                if (!Schema::hasColumn('user_customize_tags_relation', 'relation_id')) {
                    $table->unsignedInteger('relation_id')->default(0)->comment('员工标记外部联系人ID')->after('worker_id')->index();
                }
                if (Schema::hasColumn('user_customize_tags_relation', 'union_id')) {
                    $table->dropColumn('union_id');
                }
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
        //
        Schema::table('user_customize_tags_relation', function ($table) {

        });
    }
}
