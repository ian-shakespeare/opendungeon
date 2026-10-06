import { error, ok, type Result } from "result";

export default class BitMap {
  private _size: number;
  readonly bytes: ArrayBuffer;

  constructor(size: number) {
    const byteSize = Math.ceil(size / 8);
    this.bytes = new ArrayBuffer(byteSize, { maxByteLength: byteSize });
    this._size = size;
  }

  get size(): number {
    return this._size;
  }

  private get view(): DataView {
    return new DataView(this.bytes);
  }

  resize(size: number) {
    const byteSize = Math.ceil(size / 8);
    this.bytes.resize(byteSize);
    this._size = size;
  }

  get(index: number): Result<boolean, string> {
    if (index >= this.size) {
      return error("Index out of range.");
    }

    const byteIndex = Math.floor(index / 8);
    const byte = this.view.getUint8(byteIndex);
    const bitIndex = index % 8;
    const mask = 0x00000001 << bitIndex;

    return ok((byte & mask) !== 0x0);
  }

  set(index: number, value: boolean): Result<void, string> {
    if (index >= this.size) {
      return error("Index out of range.");
    }

    const view = this.view;
    const byteIndex = Math.floor(index / 8);
    const byte = view.getUint8(byteIndex);
    const bitIndex = index % 8;
    const mask = 0x00000001 << bitIndex;

    if (value) {
      view.setUint8(byteIndex, byte | mask);
    } else {
      view.setUint8(byteIndex, byte & ~mask);
    }

    return ok(undefined);
  }
}
