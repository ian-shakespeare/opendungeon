import { PriorityQueue } from "#lib/priorityqueue.js";

interface Node {
  getNeighbors: () => Iterable<this>;
  distance: (to: this) => number;
  isEqual: (b: this) => boolean;
  toString: () => string;
}

type IsAccessibleFunc<N> = (node: N) => boolean;

export function aStar<N extends Node>(start: N, end: N, isAccessible: IsAccessibleFunc<N>): N[] {
  if (start.isEqual(end)) {
    return [start];
  }

  if (!isAccessible(end)) {
    return [];
  }

  const frontier = new PriorityQueue<{ node: N; cost: number }>();
  frontier.push({ node: start, cost: 0 }, 0);
  const costSoFar = new Map<string, number>();
  costSoFar.set(start.toString(), 0);
  const cameFrom = new Map<string, N>();
  cameFrom.set(start.toString(), start);

  while (frontier.length > 0) {
    const {
      value: { node: current, cost },
    } = frontier.pop()!;

    const isMoreExpensive = cost > costSoFar.get(current.toString())!;
    if (isMoreExpensive) {
      continue;
    }

    if (current.isEqual(end)) {
      const path: N[] = [];
      let node = current;
      while (!node.isEqual(start)) {
        path.push(node);
        node = cameFrom.get(node.toString())!;
      }
      return path.reverse();
    }

    for (const next of current.getNeighbors()) {
      if (!isAccessible(next)) {
        continue;
      }

      const nextStr = next.toString();
      const newCost = costSoFar.get(current.toString())! + current.distance(next);
      const heuristic = next.distance(end);
      if (!costSoFar.has(nextStr) || newCost < costSoFar.get(nextStr)!) {
        costSoFar.set(nextStr, newCost);
        frontier.push({ node: next, cost: newCost }, newCost + heuristic);
        cameFrom.set(nextStr, current);
      }
    }
  }

  return [];
}
