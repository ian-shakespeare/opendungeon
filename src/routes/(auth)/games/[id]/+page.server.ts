import { error, redirect } from "@sveltejs/kit";
import type { Actions, PageServerLoad } from "./$types";
import { listUserCharacters } from "#lib/server/database/characters.js";
import { getUserGameWithPlayerProfiles, updateGame } from "#lib/server/database/games.js";
import { getProfile } from "#lib/server/database/profiles.js";
import { listUserLevels } from "#lib/server/database/levels.js";
import { files } from "#lib/server/files/index.js";
import { type GameState, gamerooms, PLAYER_COLORS } from "#lib/server/gamerooms.js";
import { listFriendProfiles } from "#lib/server/database/friends.js";
import { createPlayer } from "#lib/server/database/players.js";
import { fail } from "@sveltejs/kit";

export const load: PageServerLoad = async ({ locals, params }) => {
  const { session } = locals;
  if (!session) {
    redirect(303, "/sign-in");
  }

  const { id: gameId } = params;

  // three cases
  // 1. game is running
  // 2. game is not running and is stored
  // 3. game is not running and is not stored

  let [profile, game, state, levels, friends] = await Promise.all([
    getProfile(session.user_id),
    getUserGameWithPlayerProfiles(session.user_id, gameId),
    gamerooms.getGame(gameId),
    listUserLevels(session.user_id),
    listFriendProfiles(session.user_id),
  ]);

  if (!profile) {
    redirect(303, "/profile/new");
  }

  if (!game) {
    error(404, "Game not found");
  }

  const gameMaster = game.players.find(
    ({ permission_level }) => permission_level === "game_master",
  );
  if (!gameMaster) {
    error(400, "Game must have a game master");
  }

  if (!state) {
    const file = files.file(game.uri);
    const exists = await file.exists();
    state = exists
      ? ((await file.json()) as GameState)
      : {
          characters: {},
          levelUri: null,
          gameMaster: {
            userId: gameMaster.user_id,
            username: gameMaster.username,
            avatarUri: gameMaster.avatar_uri,
            online: session.user_id === gameMaster.user_id,
          },
          players: {},
          lastDisconnect: Date.now(),
        };
    await Promise.all([
      gamerooms.createGame(game.game_id, state),
      updateGame(session.user_id, gameId, { is_active: true }),
    ]);
  }

  return {
    friends: friends.filter((f) => f.accepted),
    profile,
    game,
    state,
    levels,
  };
};

// TODO: this shouldn't be a form action, move this to an endpoint
export const actions = {
  addplayer: async ({ params, request }) => {
    const data = await request.formData();

    const { id: gameId } = params;

    const userId = data.get("user-id") as string;
    if (!userId) {
      return fail(400, { userId, missing: true });
    }

    const playerProfile = await getProfile(userId);
    if (!playerProfile) {
      return fail(400, { userId, missing: true });
    }

    const game = await gamerooms.getGame(gameId);
    if (!game) {
      return fail(404, { success: false, message: "Game not found." });
    }

    await createPlayer(gameId, userId, "player");

    const availableColors = new Set(PLAYER_COLORS).difference(
      new Set(Object.values(game.players).map(({ color }) => color)),
    );
    if (availableColors.size === 0) {
      return fail(500, { success: false, message: "All colors are in use." });
    }
    const [color] = Array.from(availableColors);

    await gamerooms.addGamePlayer(gameId, {
      userId,
      username: playerProfile.username,
      avatarUri: playerProfile.avatar_uri,
      online: false,
      color,
    });
    return { success: true };
  },
} satisfies Actions;
