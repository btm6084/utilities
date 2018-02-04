## InArray
Simple utilty functions that check for the presence of a given element in a given slice, and returns of the found element. -1 is used to denote no match.

```
arr := []int{1, 2, 3}
inarray.Ints(2, arr) // 1
inarray.Ints(4, arr) // -1
```