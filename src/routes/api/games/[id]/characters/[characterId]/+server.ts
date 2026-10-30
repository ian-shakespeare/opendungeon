import { error, redirect } from "@sveltejs/kit";
import type { RequestHandler } from "./$types";
import { games } from "#lib/server/games.js";
import { ServerMessageType, type ServerMessage } from "#lib/messages.js";

export const PUT: RequestHandler = async ({ locals, params, request }) => {
  const { session } = locals;
  if (!session) {
    redirect(303, "/sign-in");
  }

  const { id: gameId, characterId } = params;
  const game = await games.get(gameId);
  if (!game) {
    error(404, "Game not found.");
  }

  const { x, y, rot }: { x: number; y: number; rot: number } = await request.json();

  // TODO: validate character CAN move and that the player has perms
  await game.updateCharacter(characterId, { userId: session.user_id, x, y, rot });

  // TODO: include rot
  const message: ServerMessage = {
    type: ServerMessageType.CharacterMoved,
    characterId,
    x,
    y,
  };
  await game.publish(message);

  return new Response(null, { status: 204 });
};
