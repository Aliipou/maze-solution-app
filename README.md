<div align="center">

[![Python](https://img.shields.io/badge/Python-3.10+-3776AB?style=flat&amp;logo=python)](https://python.org)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat)](LICENSE)

# Maze Solver

**Interactive maze solver with BFS, DFS, and A* visualization.**

</div>

## Overview

Generate random mazes, then watch three classic pathfinding algorithms find their way through. Each algorithm is visualized step-by-step so you can see exactly how they differ in behavior and efficiency.

## Algorithms

**BFS (Breadth-First Search)**
Explores all neighbors at the current depth before going deeper. Always finds the shortest path. Can be slow on large mazes because it explores in all directions equally.

**DFS (Depth-First Search)**
Dives as deep as possible before backtracking. Uses much less memory than BFS. Does not guarantee the shortest path.

**A* (A-Star)**
Guided by a heuristic (Manhattan distance to goal). Finds the shortest path much faster than BFS on average. The right choice for most practical pathfinding problems.

## Quick Start

```bash
git clone https://github.com/Aliipou/maze-solution-app.git
cd maze-solution-app
pip install -r requirements.txt
python main.py
```

Use arrow keys to select algorithm, Enter to run, G to generate a new maze.

## Complexity Comparison

| Algorithm | Time | Space | Shortest Path? |
|-----------|------|-------|----------------|
| BFS | O(V+E) | O(V) | Yes |
| DFS | O(V+E) | O(V) | No |
| A* | O(E log V) | O(V) | Yes (with admissible heuristic) |

## License

MIT
