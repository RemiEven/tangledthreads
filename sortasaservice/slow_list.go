package main

import (
	"context"
	"math/rand"
	"slices"
	"time"
)

type SlowList []int

func (sl SlowList) Get(i int) int {
	time.Sleep(5 * time.Millisecond)
	return sl[i]
}

func (sl SlowList) Set(i int, value int) {
	time.Sleep(10 * time.Millisecond)
	sl[i] = value
}

// RandomSort checks if sl is already sorted, and if it isn't, randomly swaps two consecutive elements
func (sl SlowList) RandomSort(ctx context.Context) error {
	if slices.IsSorted(sl) {
		return nil
	}
	firstIndex := rand.Intn(len(sl))

	secondIndex := (firstIndex + 1) % len(sl)

	if err := sl.swap(ctx, firstIndex, secondIndex); err != nil {
		return err
	}
	return sl.RandomSort(ctx)
}

func (sl SlowList) swap(ctx context.Context, i, j int) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	iValue := sl.Get(i)
	if err := ctx.Err(); err != nil {
		return err
	}
	jValue := sl.Get(j)
	sl.Set(j, iValue)
	if err := ctx.Err(); err != nil {
		return err
	}
	sl.Set(i, jValue)
	if err := ctx.Err(); err != nil {
		return err
	}
	return nil
}

func (sl SlowList) BubbleSort(ctx context.Context) error {
	hasPermutted := true
	for i := range len(sl) {
		if !hasPermutted { // no permutations during the last iteration - the list is sorted
			return nil
		}
		hasPermutted = false
		for j := range len(sl) - i - 1 {
			value := sl.Get(j)
			if err := ctx.Err(); err != nil {
				return err
			}
			nextValue := sl.Get(j + 1)
			if err := ctx.Err(); err != nil {
				return err
			}
			if value > nextValue {
				hasPermutted = true
				sl.Set(j, nextValue)
				if err := ctx.Err(); err != nil {
					return err
				}
				sl.Set(j+1, value)
				if err := ctx.Err(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (sl SlowList) QuickSort(ctx context.Context) error {
	return sl.quicksort(ctx, 0, len(sl)-1)
}

func (sl SlowList) quicksort(ctx context.Context, low int, high int) error {
	if low >= 0 && high >= 0 && low < high {
		pivot, err := sl.partition(ctx, low, high)
		if err != nil {
			return err
		}
		if err := sl.quicksort(ctx, low, pivot); err != nil {
			return err
		}
		if err := sl.quicksort(ctx, pivot+1, high); err != nil {
			return err
		}
	}
	return nil
}

func (sl SlowList) partition(ctx context.Context, low int, high int) (int, error) {
	pivot := sl.Get(low)
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	i := low
	j := high
	for {
		valueI := sl.Get(i)
		if err := ctx.Err(); err != nil {
			return 0, err
		}
		valueJ := sl.Get(j)
		if err := ctx.Err(); err != nil {
			return 0, err
		}
		for valueI < pivot {
			i++
			valueI = sl.Get(i)
			if err := ctx.Err(); err != nil {
				return 0, err
			}
		}
		for valueJ > pivot {
			j--
			valueJ = sl.Get(j)
			if err := ctx.Err(); err != nil {
				return 0, err
			}
		}
		if i >= j {
			return j, nil
		}
		sl.Set(i, valueJ)
		if err := ctx.Err(); err != nil {
			return 0, err
		}
		sl.Set(j, valueI)
		if err := ctx.Err(); err != nil {
			return 0, err
		}
	}
}
