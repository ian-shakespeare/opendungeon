import type { Result } from "result";

export const GRID_HEIGHT = 256;
export const GRID_WIDTH = 256;

export const MINIMUM_ZOOM = 5;
export const MAXIMUM_ZOOM = 150;

export const MAXIMUM_ROTATION = 360;
export const MINIMUM_SCALE = 0.5;
export const MAXIMUM_SCALE = 3;

export type GameMessage = {
  content: string;
  username?: string;
  avatarUri?: string | null;
  subcontent?: string;
};

export enum GameMenuTab {
  Chat = "game-icons:chat-bubble",
  Players = "game-icons:tabletop-players",
  Levels = "game-icons:treasure-map",
  Characters = "game-icons:character",
  Settings = "game-icons:cog",
}

export enum GameMenuToolIcon {
  Select = "bxs:pointer",
  Measure = "mdi:ruler",
  Shape = "material-symbols:shapes",
  Draw = "material-symbols:draw",
  Dice = "fa-solid:dice-d20",
}

export enum MeasureShape {
  Path = "boxicons:path",
  Line = "ant-design:line-outlined",
  Square = "akar-icons:square",
  Circle = "akar-icons:circle",
  Cone = "fluent:cone-16-regular",
}

export type GameTool = { type: string; icon: GameMenuToolIcon };

/*
type Player = { userId: string; username: string; avatarUri: string | null; online: boolean };
type Character = { userId: string, uri: string, x: number, y: number, rotation: number };
type Chat = {};

type PlayerJoinHandler = (player: Omit<Player, "online">) => void;
type PlayerLeaveHandler = (userId: string) => void;
type CharacterLoadHandler = (character: any) => void;
type CharacterMoveHandler = (a: any) => void;
type ChatReceiveHandler = (chat: )

export class Game {
  private playerJoinHandler: PlayerJoinHandler | undefined;
  private playerLeaveHandler: PlayerLeaveHandler | undefined;

  private constructor() {}

  static async join(gameId: string): Promise<Result<Game, string>> {
  }

  onPlayerJoin(cb: PlayerJoinHandler) {
    this.playerJoinHandler = cb;
  }

  onPlayerLeave(cb: PlayerLeaveHandler) {
    this.playerLeaveHandler = cb;
  }

  onCharacterLoad() {}

  onCharacterMove() {}

  onChatReceive() {}

  onDiceRoll() {}

  onPing() {}

  characterMove() {}

  sendChat() {}

  ping(x: number, y: number) {}
}
*/
