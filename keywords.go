package crawlrate

type ByAlpha []string

func (a ByAlpha) Len() int           { return len(a) }
func (a ByAlpha) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAlpha) Less(i, j int) bool { return a[i] < a[j] }

/*
Usage:

kw := []string{"casino", "bingo", "poker"} 
sort.Sort(ByAlpha(kw))

*/
