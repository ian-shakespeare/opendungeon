import type { ServerMessage } from "#lib/messages.js";
import { keystore, type KeyStore } from "#lib/server/keystore/index.js";

export const PLAYER_COLORS = ["orange", "yellow", "green", "aqua", "blue", "purple"] as const;

export type PlayerColor = (typeof PLAYER_COLORS)[number];

export type GamePlayer = {
  userId: string;
  username: string;
  avatarUri: string | null;
  online: boolean;
  color: PlayerColor;
};

export type GameCharacter = {
  characterId: string;
  userId: string;
  uri: string;
  x: number;
  y: number;
  /** in radians */
  rotation: number;
};

export type GameState = {
  characters: Record<string, Omit<GameCharacter, "characterId">>;
  levelUri: string | null;
  gameMaster: Omit<GamePlayer, "color">;
  players: Record<string, Omit<GamePlayer, "userId">>;
  lastDisconnect: number;
};

export default class GameRooms {
  private keystore: KeyStore;

  constructor(keystore: KeyStore) {
    this.keystore = keystore;
  }

  async publish(gameId: string, message: ServerMessage) {
    this.keystore.publish(gameId, JSON.stringify(message));
  }

  async subscribe(
    gameId: string,
    callbackfn: (message: ServerMessage) => void,
  ): Promise<() => void> {
    const id = await this.keystore.subscribe(gameId, (value: string) => {
      const message: ServerMessage = JSON.parse(value);
      callbackfn(message);
    });

    return () => this.keystore.unsubscribe(id);
  }

  async createGame(gameId: string, initialState: GameState) {
    await this.keystore.jsonSet(gameId, "$", initialState);
  }

  async getGame(gameId: string): Promise<GameState | null> {
    const state = await this.keystore.jsonGet(gameId, "$");
    if (!state) {
      return null;
    }

    return state as GameState;
  }

  async deleteGame(gameId: string) {
    await this.keystore.delete(gameId);
  }

  async setGameLevelUri(gameId: string, levelUri: string) {
    await this.keystore.jsonSet(gameId, "$.levelUri", levelUri);
  }

  async updateGameMaster(gameId: string, online: boolean) {
    await this.keystore.jsonSet(gameId, `$.gameMaster.online`, online);
  }

  async addGamePlayer(gameId: string, { userId, username, avatarUri, color, online }: GamePlayer) {
    await this.keystore.jsonSet(gameId, `$.players.${userId}`, {
      username,
      avatarUri,
      color,
      online,
    });
  }

  async getGamePlayer(gameId: string, userId: string): Promise<GamePlayer | null> {
    const player = await this.keystore.jsonGet(gameId, `$.players.${userId}`);
    if (!player) {
      return null;
    }

    return { ...player, userId } as GamePlayer;
  }

  async updateGamePlayer(gameId: string, userId: string, online: boolean) {
    await this.keystore.jsonSet(gameId, `$.players.${userId}.online`, online);
  }

  async addGameCharacter(
    gameId: string,
    { userId, characterId, uri, x, y, rotation }: GameCharacter,
  ) {
    await this.keystore.jsonSet(gameId, `$.characters.${characterId}`, {
      userId,
      uri,
      x,
      y,
      rotation,
    });
  }

  async getGameCharacter(gameId: string, characterId: string): Promise<GameCharacter | null> {
    const character = await this.keystore.jsonGet(gameId, `$.characters.${characterId}`);
    if (!character) {
      return null;
    }

    return { ...character, characterId } as GameCharacter;
  }

  async updateGameCharacter(
    gameId: string,
    characterId: string,
    { userId, uri, x, y, rotation }: Omit<GameCharacter, "characterId">,
  ) {
    await this.keystore.jsonSet(gameId, `$.characters.${characterId}`, {
      userId,
      uri,
      x,
      y,
      rotation,
    });
  }
}

export const gamerooms = new GameRooms(keystore);
