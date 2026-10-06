import { error, unwrap, ok, type Result } from "result";
import BitMap from "./bitmap";

const MAX_ENTITY_COUNT = 1_000_000;

type Entity = number;

class ComponentArray<C> {
  readonly activeEntities: BitMap;
  readonly components: C[];

  constructor(size: number) {
    this.activeEntities = new BitMap(size);
    this.components = new Array(size);
  }

  add(entity: Entity, component: C): Result<void, string> {
    if (entity > MAX_ENTITY_COUNT) {
      return error("Exceeded max entity count.");
    }

    if (entity >= this.components.length) {
      this.components.length = Math.min(2 * entity, MAX_ENTITY_COUNT);
    }

    if (entity >= this.activeEntities.size) {
      this.activeEntities.resize(Math.min(2 * entity, MAX_ENTITY_COUNT));
    }

    const result = this.activeEntities.set(entity, true);
    if (!result.ok) {
      return result;
    }

    this.components[entity] = component;

    return ok(undefined);
  }

  remove(entity: Entity): Result<void, string> {
    if (entity >= this.activeEntities.size) {
      return ok(undefined);
    }

    return this.activeEntities.set(entity, false);
  }

  map(callbackfn: (component: C) => C | void) {
    const view = new DataView(this.activeEntities.bytes);

    for (let i = 0; i < this.activeEntities.bytes.byteLength; i++) {
      const byte = view.getUint8(i);
      if (byte === 0x0) {
        continue;
      }

      const entityOffset = 8 * i;
      for (let j = 0; j < 8; j++) {
        const mask = 0x00000001 << j;
        if ((byte & mask) !== 0) {
          const entity = entityOffset + j;
          const value = callbackfn(this.components[entity]!);
          if (value) {
            this.components[entity]! = value;
          }
        }
      }
    }
  }

  static *intersect(
    ...components: ComponentArray<unknown>[]
  ): Iterator<{ entity: Entity; values: unknown[] }> {
    const views = [];
    let minLength = Infinity;
    for (let i = 0; i < components.length; i++) {
      const { activeEntities } = components[i]!;
      if (activeEntities.bytes.byteLength < minLength) {
        minLength = activeEntities.bytes.byteLength;
      }
      views.push(new DataView(activeEntities.bytes));
    }

    for (let i = 0; i < minLength; i++) {
      let byte = 0x0;
      for (let j = 0; j < views.length; j++) {
        byte &= views[j]!.getUint8(i);
      }

      const entityOffset = 8 * i;
      for (let j = 0; j < 8; j++) {
        const mask = 0x00000001 << j;
        if ((byte & mask) !== 0) {
          const entity = entityOffset + j;
          const values: unknown[] = [];

          for (let k = 0; k < components.length; k++) {
            values.push(components[k]!.components[entity]);
          }

          yield { entity, values };
        }
      }
    }
  }

  static *union(
    ...components: ComponentArray<unknown>[]
  ): Iterator<{ entity: Entity; values: (unknown | null)[] }> {
    const views = [];
    let maxLength = 0;
    for (let i = 0; i < components.length; i++) {
      const { activeEntities } = components[i]!;
      if (activeEntities.bytes.byteLength > maxLength) {
        maxLength = activeEntities.bytes.byteLength;
      }
      views.push(new DataView(activeEntities.bytes));
    }

    for (let i = 0; i < maxLength; i++) {
      let byte = 0x0;
      for (let j = 0; j < views.length; j++) {
        if (views[j]!.byteLength < i) {
          continue;
        }
        byte &= views[j]!.getUint8(i);
      }

      const entityOffset = 8 * i;
      for (let j = 0; j < 8; j++) {
        const mask = 0x00000001 << j;
        if ((byte & mask) !== 0) {
          const entity = entityOffset + j;
          const values: unknown[] = new Array(components.length).fill(null);

          for (let k = 0; k < components.length; k++) {
            const component = components[k]!;
            if (unwrap(component.activeEntities.get(entity))) {
              values[k] = component.components[entity];
            }
          }

          yield { entity, values };
        }
      }
    }
  }
}

type ComponentID = number;

const DEFAULT_COMPONENT_ARRAY_SIZE = 5000;

class System {
  private entityHandle: number;
  private deletedEntities: number[];
  private components: ComponentArray<unknown>[];

  constructor() {
    this.entityHandle = 0;
    this.deletedEntities = [];
    this.components = [];
  }

  registerComponent(): ComponentID {
    const id = this.components.length;
    this.components.push(new ComponentArray(5000));
    return id;
  }

  spawnEntity(): Entity {
    const entity = this.entityHandle;
    this.entityHandle++;
    return entity;
  }
}
