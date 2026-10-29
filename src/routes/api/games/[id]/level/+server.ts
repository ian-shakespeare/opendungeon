import { error, redirect } from "@sveltejs/kit";
import type { RequestHandler } from "./$types";
import { games } from "#lib/server/games.js";
import { ServerMessageType, type LevelLoaded } from "#lib/messages.js";
import { getUserLevel, type LevelData } from "#lib/server/database/levels.js";
import { files } from "#lib/server/files/index.js";

export const PUT: RequestHandler = async ({ locals, params, request }) => {
  const { session } = locals;
  if (!session) {
    redirect(303, "/sign-in");
  }

  const { id: gameId } = params;

  const game = games[gameId];
  if (!game) {
    error(404, "Game not found.");
  }

  const state = await game.getState();

  if (!state.players[session.user_id].isGameMaster) {
    error(403, "Players may not load levels.");
  }

  const { levelId }: { levelId: string } = await request.json();
  const level = await getUserLevel(session.user_id, levelId);
  if (!level) {
    error(404, "Level not found.");
  }

  const file = files.file(level.uri);
  const exists = await file.exists();
  if (!exists) {
    error(404, "Level data not found.");
  }

  const data: LevelData = await file.json();
  await game.setLevelUri(level.uri);

  const message: LevelLoaded = {
    type: ServerMessageType.LevelLoaded,
    name: level.name,
    data,
  };
  await game.publish(message);

  return new Response(null, { status: 204 });
};
