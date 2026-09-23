import * as GLM from "gl-matrix";
import type { RenderElement } from "$lib/renderer/element";
import DynamicModel from "$lib/renderer/model/dynamic";
import type Renderer from "$lib/renderer";
import HumanMale from "$lib/assets/human.male.glb?url";
import { getGLBChunks } from "$lib/renderer/model/glb";
import assert from "$lib/assert";
import type ModelInstance from "$lib/renderer/model/instance";
import type { Camera } from "$lib/renderer/camera";
import { getGLTFModelParams } from "$lib/renderer/model/gltf";
import Shader from "$lib/renderer/shader";
import vertex from "$lib/assets/shaders/character.vert?raw";
import fragment from "$lib/assets/shaders/character.frag?raw";

/**
 * A custom element for rendering editable characters.
 * Should NOT be used outside the character creator.
 */
export default class Character implements RenderElement {
  private model: DynamicModel;
  private instance: ModelInstance;

  private constructor(model: DynamicModel, instance: ModelInstance) {
    this.model = model;
    this.instance = instance;
  }

  static async create(renderer: Renderer): Promise<Character> {
    const response = await fetch(HumanMale, {
      credentials: import.meta.env.DEV ? "include" : "same-origin",
    });
    assert(response.ok, "failed to get glb");

    const blob = await response.blob();
    const { source, data } = getGLBChunks(await blob.arrayBuffer());
    const shader = new Shader(renderer.gl, vertex, fragment);
    shader.loadUniformLocation("u_model");
    shader.loadUniformLocation("u_view");
    shader.loadUniformLocation("u_projection");
    shader.loadUniformLocation("u_base_color");
    shader.loadUniformLocation("u_alpha_cutoff");
    shader.loadUniformLocation("u_joint_matrix[0]");

    const params = await getGLTFModelParams(shader, source, {
      instanced: false,
      preloadedBuffers: [data],
    });
    const model = new DynamicModel(params);
    const instance = model.createInstance();
    instance.updateTransforms();
    instance.computeSkinningMatrix();

    return new Character(model, instance);
  }

  destroy() {
    this.model.destroy();
  }

  draw() {
    this.model.draw();
  }

  use() {
    this.model.use();
  }

  setCamera(camera: Camera) {
    this.model.setCamera(camera);
  }

  setSkinTone(color: Float32Array) {
    this.model.defaultMaterial = {
      name: "default",
      baseColorFactor: color,
      alphaMode: "OPAQUE",
      alphaCutoff: 0.5,
      doubleSided: false,
    };
  }

  translate(v: GLM.vec3) {
    GLM.mat4.translate(this.instance.transform, this.instance.transform, v);
  }
}
