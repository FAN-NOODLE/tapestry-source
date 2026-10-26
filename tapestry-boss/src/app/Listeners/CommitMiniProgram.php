<?php
/*
 * This file is part of the car/chedianai_bc.
 *
 * (c) chedianai <i@chedianai.com>
 *
 * This source file is subject to the MIT license that is bundled
 * with this source code in the file LICENSE.
 */

namespace App\Listeners;

use App\Events\MiniProgramAuthed;
use App\Services\OpenPlatform\MiniProgramService;
use Illuminate\Contracts\Queue\ShouldQueue;

class CommitMiniProgram implements ShouldQueue
{
    protected $miniProgramService;

    public function __construct(MiniProgramService $miniProgramService)
    {
        $this->miniProgramService = $miniProgramService;
    }

    /**
     * Handle the event.
     *
     * @param MiniProgramAuthed $event
     */
    public function handle(MiniProgramAuthed $event)
    {
        $thirdPartAccount = $event->thirdPartAccount;
        if (!$thirdPartAccount) {
            return;
        }

        try {
            $this->miniProgramService->submitVersion($thirdPartAccount, true);
        } catch (\Exception $e) {
            app('log')->error($e);
        }
    }
}
