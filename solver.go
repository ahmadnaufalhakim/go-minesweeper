package main

import (
	"slices"
)

type Constraint struct {
	UnknownNeighbors [][2]int
	RemainingValue   int
}

func (m Minesweeper) DeterministicSolve(maxComponentSize int) (bool, map[[2]int]struct{}, map[[2]int]struct{}) {
	// Typed position sets helper functions
	newPositionSet := func() map[[2]int]struct{} {
		return make(map[[2]int]struct{})
	}
	add := func(positionSet map[[2]int]struct{}, position [2]int) {
		positionSet[position] = struct{}{}
	}
	has := func(positionSet map[[2]int]struct{}, position [2]int) bool {
		_, ok := positionSet[position]
		return ok
	}

	// Start cell guard
	if m.StartCell == nil {
		return false, newPositionSet(), newPositionSet()
	}

	// Initial flood reveal from clicking start cell
	revealed := newPositionSet()
	queue := [][2]int{{m.StartCellPosition[0], m.StartCellPosition[1]}}
	for len(queue) > 0 {
		position := queue[0]
		queue = queue[1:]
		pos := [2]int{position[0], position[1]}
		if has(revealed, pos) {
			continue
		}

		add(revealed, pos)
		// If clear/zero, add neighbors to queue (flood fill)
		if m.Grid[pos[0]][pos[1]].Value == CLEAR {
			for _, neighbor := range m.getNeighborsOf(pos[0], pos[1]) {
				if !has(revealed, neighbor) && !slices.Contains(m.BombPositions, neighbor) {
					queue = append(queue, neighbor)
				}
			}
		}
	}

	flagged := newPositionSet()
	changed := true

	// Helper function to compute remaining value (== number - flagged cells around)
	remainingValue := func(row, col int) int {
		count := 0
		for _, neighbor := range m.getNeighborsOf(row, col) {
			if has(flagged, neighbor) {
				count++
			}
		}
		return m.PositionToValue[[2]int{row, col}] - count
	}

	// Helper function to get all neighbors that are not revealed and not flagged
	unknownNeighbors := func(row, col int) [][2]int {
		res := make([][2]int, 0)
		for _, neighbor := range m.getNeighborsOf(row, col) {
			if !has(revealed, neighbor) && !has(flagged, neighbor) {
				res = append(res, neighbor)
			}
		}
		return res
	}

mainLoop:
	for {
		// If nothing is changed in previous iteration, then done
		if !changed {
			break
		}
		changed = false
		madeProgress := false

		// --- Build frontier and constraints ---
		frontier := newPositionSet()
		constraints := make(map[[2]int]Constraint)
		for row := range m.Rows {
			for col := range m.Cols {
				cell := [2]int{row, col}
				if !has(revealed, cell) || m.Grid[row][col].Value == CLEAR {
					continue
				}

				unk := unknownNeighbors(row, col)
				if len(unk) == 0 {
					continue
				}

				rem := remainingValue(row, col)
				constraints[cell] = Constraint{
					UnknownNeighbors: unk,
					RemainingValue:   rem,
				}
				for _, u := range unk {
					add(frontier, u)
				}
			}
		}

		// If there's nothing left to deduce
		if len(frontier) == 0 {
			break
		}

		// --- Apply simple local rules (only on constraints/frontier) ---
		for _, constraint := range constraints {
			rem := constraint.RemainingValue
			unk := constraint.UnknownNeighbors

			if rem == 0 && len(unk) > 0 {
				for _, u := range unk {
					if !has(revealed, u) {
						add(revealed, u)
						madeProgress = true
					}
				}
			} else if rem == len(unk) && len(unk) > 0 {
				for _, u := range unk {
					if !has(flagged, u) {
						add(flagged, u)
						madeProgress = true
					}
				}
			}
		}

		// If progress has been made, flood newly revealed zeros and repeat
		if madeProgress {
			// Flood neighbors of zero values
			// Also collect revealed into slice to avoid mutation race
			revealedList := make([][2]int, 0, len(revealed))
			for pos := range revealed {
				revealedList = append(revealedList, pos)
			}
			for _, pos := range revealedList {
				if m.PositionToValue[pos] == 0 {
					for _, neighbor := range m.getNeighborsOf(pos[0], pos[1]) {
						if !has(revealed, neighbor) && !has(flagged, neighbor) && !slices.Contains(m.BombPositions, neighbor) {
							add(revealed, neighbor)
						}
					}
				}
			}
			changed = true
			continue mainLoop
		}

		// --- Build adjacency among frontier unknowns (if they appear together in a constraint) ---
		adj := make(map[[2]int]map[[2]int]struct{})
		ensureAdj := func(a [2]int) {
			if _, ok := adj[a]; !ok {
				adj[a] = newPositionSet()
			}
		}
		for _, constraint := range constraints {
			unk := constraint.UnknownNeighbors
			for i := range len(unk) {
				for j := range len(unk) {
					if i == j {
						continue
					}
					a := unk[i]
					b := unk[j]
					ensureAdj(a)
					add(adj[a], b)
				}
			}
		}

		// Ensure nodes in frontier exist in adj (maybe isolated)
		for pos := range frontier {
			ensureAdj(pos)
		}

		// --- Connected components on frontier graph ---
		components := make([][][2]int, 0)
		seen := newPositionSet()
		for node := range frontier {
			if _, ok := seen[node]; ok {
				continue
			}
			// BFS
			component := make([][2]int, 0)
			queue := [][2]int{node}
			add(seen, node)
			for len(queue) > 0 {
				curNode := queue[0]
				queue = queue[1:]
				component = append(component, curNode)
				for neighbor := range adj[curNode] {
					if _, ok := seen[neighbor]; !ok {
						add(seen, neighbor)
						queue = append(queue, neighbor)
					}
				}
			}
			components = append(components, component)
		}

		// --- Evaluate each component by brute-force ---
		forcedSafe := newPositionSet()
		forcedMine := newPositionSet()

		for _, component := range components {
			if len(component) == 0 {
				continue
			}
			if len(component) > maxComponentSize {
				// Skip performing brute-force evaluation on large components - treat as non-deducible
				continue
			}

			componentList := component
			N := len(componentList)

			// Build relevant constraints: intersection with this component and adjusted remaining (subtract flagged outside component)
			type relevantConstraint struct {
				Inter  [][2]int
				RemAdj int
			}
			relevantConstraints := make([]relevantConstraint, 0)
			for _, constraint := range constraints {
				inter := make([][2]int, 0)
				for _, u := range constraint.UnknownNeighbors {
					if slices.Contains(componentList, u) {
						inter = append(inter, u)
					}
				}
				if len(inter) == 0 {
					continue
				}
				// rem adjusted by flagged outside component
				flaggedOutside := 0
				for _, u := range constraint.UnknownNeighbors {
					if !slices.Contains(componentList, u) {
						if has(flagged, u) {
							flaggedOutside++
						}
					}
				}
				relevantConstraints = append(relevantConstraints, relevantConstraint{
					Inter:  inter,
					RemAdj: constraint.RemainingValue - flaggedOutside,
				})
			}

			// Precompute map from cell -> index in componentList
			indexOf := make(map[[2]int]int, N)
			for i, u := range componentList {
				indexOf[u] = i
			}
			// Enumerate assignments: 0..(1<<N)-1
			bombCounts := make([]int, N)
			totalAssignments := 0
			total := 1 << uint(N)
			for mask := range total {
				ok := true
				for _, relevantConstraint := range relevantConstraints {
					s := 0
					for _, u := range relevantConstraint.Inter {
						idx := indexOf[u]
						if (mask>>uint(idx))&1 == 1 {
							s++
						}
					}
					if s != relevantConstraint.RemAdj {
						ok = false
						break
					}
				}

				if ok {
					totalAssignments++
					for i := range N {
						if (mask>>uint(i))&1 == 1 {
							bombCounts[i]++
						}
					}
				}
			}

			// If inconsistency happens, board is invalid
			if totalAssignments == 0 {
				return false, revealed, flagged
			}

			// Intersection across valid assignments using counts
			for i := range N {
				if bombCounts[i] == totalAssignments {
					add(forcedMine, componentList[i])
				}
				if bombCounts[i] == 0 {
					add(forcedSafe, componentList[i])
				}
			}
		}

		// If skipped a big component, treat as stuck
		if len(forcedSafe) == 0 && len(forcedMine) == 0 {
			// Nothing assignment forced by enumeration -> guessing needed -> stuck
			break
		}

		// Apply forced moves
		for pos := range forcedMine {
			if !has(flagged, pos) {
				add(flagged, pos)
				changed = true
			}
		}
		for pos := range forcedSafe {
			if !has(revealed, pos) {
				add(revealed, pos)
				changed = true
			}
		}
		// loop will rebuild frontier when it iterates again
	}

	// Finished deduction loop, check if all safe cells are revealed
	allSafeCellsRevealed := true
	for row := range m.Rows {
		for col := range m.Cols {
			pos := [2]int{row, col}
			if slices.Contains(m.BombPositions, pos) {
				continue
			}
			if !has(revealed, pos) {
				allSafeCellsRevealed = false
				break
			}
		}

		if !allSafeCellsRevealed {
			break
		}
	}

	return allSafeCellsRevealed, revealed, flagged
}
