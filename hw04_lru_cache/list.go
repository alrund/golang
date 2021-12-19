package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	length int
	front  *ListItem
	back   *ListItem
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{Value: v, Next: l.front}
	l.pushFrontItem(newItem)

	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{Value: v, Prev: l.back}
	l.pushBackItem(newItem)

	return newItem
}

func (l *list) Remove(i *ListItem) {
	if l.Len() == 1 {
		l.empty()
		return
	}

	l.length--

	if l.isFront(i) {
		l.front = i.Next
		l.front.Prev = nil
		return
	}

	if l.isBack(i) {
		l.back = i.Prev
		l.back.Next = nil
		return
	}

	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev
}

func (l *list) MoveToFront(i *ListItem) {
	if l.Len() == 1 || l.front == i {
		return
	}

	l.Remove(i)
	l.pushFrontItem(i)
}

func (l *list) empty() {
	l.front = nil
	l.back = nil
	l.length = 0
}

func (l *list) addFirstItem(newItem *ListItem) {
	l.back = newItem
	l.front = newItem
	newItem.Prev = nil
	newItem.Next = nil
	l.length++
}

func (l *list) pushFrontItem(item *ListItem) {
	oldFrontItem := l.front

	if l.Len() == 0 {
		l.addFirstItem(item)
		return
	}

	oldFrontItem.Prev = item
	item.Next = oldFrontItem

	l.front = item
	l.length++
}

func (l *list) pushBackItem(item *ListItem) {
	oldBackItem := l.back

	if l.Len() == 0 {
		l.addFirstItem(item)
		return
	}

	oldBackItem.Next = item
	item.Prev = oldBackItem

	l.back = item
	l.length++
}

func (l *list) isFront(item *ListItem) bool {
	return l.front == item
}

func (l *list) isBack(item *ListItem) bool {
	return l.back == item
}

func NewList() List {
	return new(list)
}
