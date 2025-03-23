package tangledthreads

import (
	"context"
	"slices"
	"testing"
	"testing/synctest"
	"time"
)

func TestSlowListGet(t *testing.T) {
	synctest.Run(func() {
		slowList := SlowList([]int{12})
		startTime := time.Now()
		if value := slowList.Get(0); value != 12 {
			t.Errorf("unexpected value: wanted 12, got %v", value)
			return
		}
		if d := time.Since(startTime); d.Milliseconds() != 5 {
			t.Errorf("unexpected operation duration: wanted 5, got %v", d)
			return
		}
	})
}

func TestSlowListSet(t *testing.T) {
	synctest.Run(func() {
		slowList := make(SlowList, 1)
		startTime := time.Now()
		slowList.Set(0, 49)
		if value := slowList[0]; value != 49 {
			t.Errorf("unexpected value: wanted 49, got %v", value)
			return
		}
		if d := time.Since(startTime); d.Milliseconds() != 10 {
			t.Errorf("unexpected operation duration: wanted 10, got %v", d)
			return
		}
	})
}

func TestSwap(t *testing.T) {
	synctest.Run(func() {
		slowList := SlowList([]int{0, 1, 2, 3})
		slowList.swap(context.Background(), 1, 2)
		expected := []int{0, 2, 1, 3}
		if !slices.Equal(slowList, expected) {
			t.Errorf("numbers weren't swapped as expected: wanted %+v, got %+v", expected, slowList)
		}
	})
}

func TestRandomSort(t *testing.T) {
	synctest.Run(func() {
		slowList := SlowList([]int{0, 2, 1})
		err := slowList.RandomSort(context.Background())
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		} else if !slices.IsSorted(slowList) {
			t.Errorf("list wasn't sorted as expected: got %+v", slowList)
		}
	})
}

func TestRandomSortWithContextTimeout(t *testing.T) {
	synctest.Run(func() {
		slowList := SlowList([]int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0})
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()
		err := slowList.RandomSort(ctx)
		if err == nil {
			t.Errorf("expected an error")
		} else if slices.IsSorted(slowList) {
			t.Errorf("list was unexpectedly sorted")
		}
	})
}

func TestBubbleSort(t *testing.T) {
	synctest.Run(func() {
		slowList := SlowList([]int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0})
		err := slowList.BubbleSort(context.Background())
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		} else if !slices.IsSorted(slowList) {
			t.Errorf("list wasn't sorted as expected: got %+v", slowList)
		}
	})
}

func TestBubbleSortWithContextTimeout(t *testing.T) {
	synctest.Run(func() {
		slowList := SlowList([]int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0})
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()
		err := slowList.BubbleSort(ctx)
		if err == nil {
			t.Errorf("expected an error")
		} else if slices.IsSorted(slowList) {
			t.Errorf("list was unexpectedly sorted")
		}
	})
}

func TestQuickSort(t *testing.T) {
	synctest.Run(func() {
		slowList := SlowList([]int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0})
		err := slowList.QuickSort(context.Background())
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		} else if !slices.IsSorted(slowList) {
			t.Errorf("list wasn't sorted as expected: got %+v", slowList)
		}
	})
}

func TestQuickSortWithContextTimeout(t *testing.T) {
	synctest.Run(func() {
		slowList := SlowList([]int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0})
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()
		err := slowList.QuickSort(ctx)
		if err == nil {
			t.Errorf("expected an error")
		} else if slices.IsSorted(slowList) {
			t.Errorf("list was unexpectedly sorted")
		}
	})
}

func TestRandomSortIsTheFastestForAnAlreadySortedList(t *testing.T) {
	createAlreadySortedList := func() SlowList {
		return []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	}
	synctest.Run(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		var (
			randomSortTime = timeRun(ctx, createAlreadySortedList().RandomSort)
			bubbleSortTime = timeRun(ctx, createAlreadySortedList().BubbleSort)
			quickSortTime  = timeRun(ctx, createAlreadySortedList().QuickSort)
		)
		if randomSortTime > bubbleSortTime {
			t.Errorf("bubble sort was unexpectedly faster than random sort")
		}
		if randomSortTime > quickSortTime {
			t.Errorf("quicksort was unexpectedly faster than random sort")
		}
	})
}

func TestBubbleSortIsTheFastestForANearlySortedList(t *testing.T) {
	createNearlySortedList := func() SlowList {
		return []int{0, 1, 2, 9, 3, 4, 5, 6, 7, 8}
	}
	synctest.Run(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		var (
			randomSortTime = timeRun(ctx, createNearlySortedList().RandomSort)
			bubbleSortTime = timeRun(ctx, createNearlySortedList().BubbleSort)
			quickSortTime  = timeRun(ctx, createNearlySortedList().QuickSort)
		)
		if bubbleSortTime > randomSortTime {
			t.Errorf("random sort was unexpectedly faster than bubble sort")
		}
		if bubbleSortTime > quickSortTime {
			t.Errorf("quicksort was unexpectedly faster than bubble sort")
		}
	})
}

func TestQuickSortIsTheFastestForAMessyList(t *testing.T) {
	createMessyList := func() SlowList {
		return []int{8, 4, 6, 2, 5, 3, 9, 7, 0}
	}
	synctest.Run(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		var (
			randomSortTime = timeRun(ctx, createMessyList().RandomSort)
			bubbleSortTime = timeRun(ctx, createMessyList().BubbleSort)
			quickSortTime  = timeRun(ctx, createMessyList().QuickSort)
		)
		if quickSortTime > bubbleSortTime {
			t.Errorf("bubble sort was unexpectedly faster than quicksort")
		}
		if quickSortTime > randomSortTime {
			t.Errorf("random sort was unexpectedly faster than quicksort")
		}
	})
}

func timeRun(ctx context.Context, f func(context.Context) error) time.Duration {
	startTime := time.Now()
	f(ctx)
	return time.Since(startTime)
}
