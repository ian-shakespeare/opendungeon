import { error, redirect } from "@sveltejs/kit";
import type { RequestHandler } from "./$types";
import { gamerooms } from "#lib/server/gamerooms.js";
import { getProfile } from "#lib/server/database/profiles.js";
import { ServerMessageType, type PlayerLeft, type ServerMessage } from "#lib/messages.js";
import { getPlayer } from "#lib/server/database/players.js";

export const GET: RequestHandler = async ({ locals, params }) => {
  const { session } = locals;
  if (!session) {
    redirect(303, "/sign-in");
  }

  const { id: gameId } = params;

  const [profile, player, state] = await Promise.all([
    getProfile(session.user_id),
    getPlayer(gameId, session.user_id),
    gamerooms.getGame(gameId),
  ]);
  if (!profile) {
    error(404, "Profile not found");
  }
  if (!player) {
    error(404, "Player not found");
  }
  if (!state) {
    error(404, "Game not found");
  }

  const isGameMaster = state.gameMaster.userId === session.user_id;

  let unsubscribe: () => void;

  const stream = new ReadableStream({
    async start(controller) {
      unsubscribe = await gamerooms.subscribe(gameId, (message) => {
        controller.enqueue(`data: ${JSON.stringify(message)}\n\n`);
      });

      const message: ServerMessage = isGameMaster
        ? await (async () => {
            await gamerooms.updateGameMaster(gameId, true);
            return {
              type: ServerMessageType.GameMasterJoined,
              userId: session.user_id,
              username: profile.username,
              avatarUri: profile.avatar_uri,
            };
          })()
        : await (async () => {
            const playerState = state.players[session.user_id];
            if (!playerState) {
              throw new Error("Player tried to join without being part of game.");
            }

            await gamerooms.updateGamePlayer(gameId, session.user_id, true);
            return {
              type: ServerMessageType.PlayerJoined,
              userId: profile.user_id,
              username: profile.username,
              avatarUri: profile.avatar_uri,
              color: playerState.color,
            };
          })();

      await gamerooms.publish(gameId, message);
    },
    async cancel() {
      unsubscribe();

      if (isGameMaster) {
        await gamerooms.updateGameMaster(gameId, false);
      } else {
        await gamerooms.updateGamePlayer(gameId, session.user_id, false);
      }

      const message: PlayerLeft = {
        type: ServerMessageType.PlayerLeft,
        userId: session.user_id,
      };
      await gamerooms.publish(gameId, message);
    },
  });

  return new Response(stream, {
    headers: {
      "Content-Type": "text/event-stream",
      "Cache-Control": "no-cache",
      Connection: "keep-alive",
    },
  });
};
