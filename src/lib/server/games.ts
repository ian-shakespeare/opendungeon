import type { ServerMessage } from "#lib/messages.js";
import { keystore as defaultKeystore, type KeyStore } from "#lib/server/keystore/index.js";
import type { PubSub } from "./pubsub";

export const PLAYER_COLORS = ["orange", "yellow", "green", "aqua", "blue", "purple"] as const;

export type PlayerColor = (typeof PLAYER_COLORS)[number];

export type GamePlayer =
  | {
      isGameMaster: true;
      userId: string;
      username: string;
      avatarUri: string | null;
      online: boolean;
    }
  | {
      isGameMaster: false;
      color: PlayerColor;
      userId: string;
      username: string;
      avatarUri: string | null;
      online: boolean;
    };

export type GameCharacter = {
  characterId: string;
  userId: string;
  uri: string;
  x: number;
  y: number;
  /** in radians */
  rot: number;
};

export type GameState = {
  characters: Record<string, Omit<GameCharacter, "characterId">>;
  levelUri: string | null;
  players: Record<string, Omit<GamePlayer, "userId">>;
  lastDisconnect: number;
};

function getInitialGameState(): GameState {
  return {
    characters: {},
    levelUri: null,
    players: {},
    lastDisconnect: Date.now(),
  };
}

export class GameServer {
  private gameId: string;
  private keystore: KeyStore;
  private pubsub: PubSub;

  private constructor(gameId: string, keystore: KeyStore, pubsub: PubSub) {
    this.gameId = gameId;
    this.keystore = keystore;
    this.pubsub = pubsub;
  }

  static async start(
    gameId: string,
    initialState = getInitialGameState(),
    keystore = defaultKeystore,
    pubsub = keystore as unknown as PubSub,
  ): Promise<GameServer> {
    await keystore.jsonSet(gameId, "$", initialState);
    return new GameServer(gameId, keystore, pubsub);
  }

  async destroy() {
    await this.keystore.jsonDelete(this.gameId, "$");
  }

  async publish(message: ServerMessage) {
    this.pubsub.publish(this.gameId, JSON.stringify(message));
  }

  async subscribe(callbackfn: (message: ServerMessage) => void): Promise<() => void> {
    const id = await this.pubsub.subscribe(this.gameId, (value: string) => {
      const message: ServerMessage = JSON.parse(value);
      callbackfn(message);
    });

    return () => this.pubsub.unsubscribe(id);
  }

  async getState(): Promise<GameState> {
    const state = await this.keystore.jsonGet(this.gameId, "$");
    if (!state) {
      throw new Error("Failed to get state for active game.");
    }

    return state as GameState;
  }

  async setLevelUri(levelUri: string) {
    await this.keystore.jsonSet(this.gameId, "$.levelUri", levelUri);
  }

  // automatically determine color here
  async addPlayer(player: GamePlayer) {
    await this.keystore.jsonSet(this.gameId, `$.players.${player.userId}`, player);
  }

  async getPlayer(userId: string): Promise<GamePlayer | null> {
    const player = await this.keystore.jsonGet(this.gameId, `$.players.${userId}`);
    if (!player) {
      return null;
    }

    return { ...player, userId } as GamePlayer;
  }

  async updatePlayer(userId: string, online: boolean) {
    await this.keystore.jsonSet(this.gameId, `$.players.${userId}.online`, online);
  }

  async addCharacter({ userId, characterId, uri, x, y, rot }: GameCharacter) {
    await this.keystore.jsonSet(this.gameId, `$.characters.${characterId}`, {
      userId,
      uri,
      x,
      y,
      rot,
    });
  }

  async getCharacter(characterId: string): Promise<GameCharacter | null> {
    const character = await this.keystore.jsonGet(this.gameId, `$.characters.${characterId}`);
    if (!character) {
      return null;
    }

    return { ...character, characterId } as GameCharacter;
  }

  async updateCharacter(
    characterId: string,
    { userId, uri, x, y, rot }: Omit<GameCharacter, "characterId">,
  ) {
    await this.keystore.jsonSet(this.gameId, `$.characters.${characterId}`, {
      userId,
      uri,
      x,
      y,
      rot,
    });
  }
}

export const games: Record<string, GameServer> = {};
