<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddSceneTypeToSubscribeMessageRecordTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasColumn('subscribe_message_record', 'scene_type')) {
            Schema::table('subscribe_message_record', function (Blueprint $table) {
                $table->unsignedTinyInteger('scene_type')->default(0)->index()->comment('场景类别（1：直播订阅）');
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
        Schema::table('subscribe_message_record', function (Blueprint $table) {
        });
    }
}
