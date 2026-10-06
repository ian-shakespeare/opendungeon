import { describe, expect, test } from "bun:test";
import BitMap from "./bitmap";

describe.concurrent("BitMap", () => {
  test("constructor sets correct size", () => {
    const map = new BitMap(8);
    expect(map.size).toBe(8);
    expect(map.get(8).ok).toBeFalse();
  });

  test("constructor allocates multiple bytes", () => {
    const map = new BitMap(9);
    expect(map.size).toBe(9);

    const { ok } = map.set(8, true);
    expect(ok).toBeTrue();

    const res = map.get(8);
    expect(res.ok).toBeTrue();

    if (res.ok) {
      expect(res.value).toBeTrue();
    }
  });

  test("set to false", () => {
    const map = new BitMap(3);

    let res = map.set(2, true);
    expect(res.ok).toBeTrue();

    res = map.set(2, false);
    expect(res.ok).toBeTrue();

    const result = map.get(2);
    expect(result.ok).toBeTrue();

    if (result.ok) {
      expect(result.value).toBeFalse();
    }
  });

  test("resize adds capacity", () => {
    const map = new BitMap(3);
    expect(map.size).toBe(3);
    expect(map.get(3).ok).toBeFalse();

    map.resize(4);
    expect(map.size).toBe(4);
    expect(map.get(3).ok).toBeTrue();
  });

  test("resize removes capacity", () => {
    const map = new BitMap(4);
    expect(map.size).toBe(4);
    expect(map.get(3).ok).toBeTrue();

    map.resize(3);
    expect(map.size).toBe(3);
    expect(map.get(3).ok).toBeFalse();
  });
});
