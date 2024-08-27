package former

// expressIterator - generates all possible unique expresses for the system depending on systemSize
// how algorithm works:
// iterator iterates expressSelIndexes and extracts express based on this index, for example:
// selections =  [ 0 => {sel1}, 1 => {sel2}, 2 => {sel3}, 3 => {sel4}, 4 => {sel5}], systemSize 3
// init expressSelIndexes will be [0, 1, 2];  express - [ 0 => {sel1}, 1 => {sel2}, 2 => {sel3}]
// increment next expressSelIndexes will be [0, 1, 3]; express - [ 0 => {sel1}, 1 => {sel2}, 3 => {sel4}]
// increment next expressSelIndexes will be [0, 1, 4]; express - [ 0 => {sel1}, 1 => {sel2}, 4 => {sel5}]
// iterator have reached the maximum value for selection position 2
// next step - previous selectionPosition is incremented and next indexes after previous are reset
// increment next expressSelIndexes will be [0, 2(incremented previous), 3(reset next after previous, (old value 4)];
// express - [ 0 => {sel1}, 2 => {sel3}, 3 => {sel4}]
// increment next expressSelIndexes will be [0, 2, 4]; express - [ 0 => {sel1}, 2 => {sel3}, 4 => {sel5}]
// increment next expressSelIndexes will be [0, 3, 4]; express - [ 0 => {sel1}, 3 => {sel4}, 4 => {sel5}]
// increment next expressSelIndexes will be [1, 2, 3]; express - [ 1 => {sel2}, 2 => {sel3}, 3 => {sel4}]
// increment next expressSelIndexes will be [1, 2, 4]; express - [ 1 => {sel2}, 2 => {sel3}, 4 => {sel5}]
type expressIterator struct {
	selectionsCount int
	systemSize      int
	// expressSelIndexes - slice of selections indexes for current express, for example:
	// init expressSelIndexes will be [0, 1, 2]
	// express based on this index will be [[ 0 => {sel1}, 1 => {sel2}, 2 => {sel3}, 3 => {sel4}]
	expressSelIndexes        []int
	lastSelPositionInExpress int
}

func newExpressIterator(systemSize int, selectionsCount int) expressIterator {
	return expressIterator{
		selectionsCount:          selectionsCount,
		systemSize:               systemSize,
		expressSelIndexes:        nil,
		lastSelPositionInExpress: systemSize - 1,
	}
}

func (e *expressIterator) next() bool {
	if e.expressSelIndexes == nil {
		e.initExpressSelIndexes()

		return true
	}

	if e.firstSelIndexIsMax() && e.lastSelIndexIsMax() {
		return false
	}

	selPositionInExpress := e.lastSelPositionInExpress

	for {
		// if can`t increment current selectionPosition, go to previous selectionPosition
		if !e.canIncrement(selPositionInExpress) {
			selPositionInExpress--
			continue
		}

		e.expressSelIndexes[selPositionInExpress]++

		if selPositionInExpress == e.lastSelPositionInExpress {
			return true
		}

		// if the incriminated index is not the last one - reset all subsequent indexes
		for i := selPositionInExpress + 1; i <= e.lastSelPositionInExpress; i++ {
			e.expressSelIndexes[i] = e.expressSelIndexes[i-1] + 1
		}

		return true
	}
}

// initExpressSelIndexes - generate initial index for first express in system
// for example:
// systemSize 3 - initial currentIndex index [0, 1, 2],
// systemSize 4 - initial currentIndex index [0, 1, 2, 3]
func (e *expressIterator) initExpressSelIndexes() {
	e.expressSelIndexes = make([]int, e.systemSize)
	for i := 0; i < e.systemSize; i++ {
		e.expressSelIndexes[i] = i
	}
}

func (e *expressIterator) firstSelIndexIsMax() bool {
	return e.expressSelIndexes[0] == e.maxSelIndexValue(0)
}

func (e *expressIterator) lastSelIndexIsMax() bool {
	return e.expressSelIndexes[e.lastSelPositionInExpress] == e.maxSelIndexValue(e.lastSelPositionInExpress)
}

func (e *expressIterator) canIncrement(selIndexPositionInExpress int) bool {
	return e.expressSelIndexes[selIndexPositionInExpress] < e.maxSelIndexValue(selIndexPositionInExpress)
}

// maxSelIndexValue - calculate max index value by selection position in express, for example
// full selections slice [ 0 => {sel1}, 1 => {sel2}, 2 => {sel3}, 3 => {sel4}, 4 => {sel5}], system size 3
// all unique combinations of selection indexes for express will be
// [0, 1, 2]
// [0, 1, 3]
// [0, 1, 4]
// [0, 2, 3]
// [0, 2, 4]
// [0, 3, 4]
// [1, 2, 3]
// [1, 2, 4]
// [1, 3, 4]
// [2, 3, 4]
// we can see for selection position 0 max index is 2
// for selection position 1 max index is 3
// for selection position 3 max index is 4
// this algorithm is valid for any system size
func (e *expressIterator) maxSelIndexValue(selPositionInExpress int) int {
	return selPositionInExpress + e.selectionsCount - e.systemSize
}
